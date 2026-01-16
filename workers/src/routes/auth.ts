import { Hono } from 'hono';
import type { Env } from '../types';

const router = new Hono<{ Bindings: Env }>();

router.get('/auth/:provider', async (c) => {
  const provider = c.req.param('provider');

  if (!['facebook', 'instagram'].includes(provider)) {
    throw new Error('Invalid OAuth provider');
  }

  const redirectUri = `${c.env.BASE_URL}/auth/${provider}/callback`;

  const { FacebookOAuthProvider, InstagramOAuthProvider } = await import('../lib/oauth');

  let oauthProvider;
  if (provider === 'facebook') {
    oauthProvider = new FacebookOAuthProvider(
      c.env.FACEBOOK_KEY,
      c.env.FACEBOOK_SECRET,
      redirectUri
    );
  } else {
    oauthProvider = new InstagramOAuthProvider(
      c.env.INSTAGRAM_KEY,
      c.env.INSTAGRAM_SECRET,
      redirectUri
    );
  }

  const authUrl = oauthProvider.getAuthUrl();

  return c.json({ auth_url: authUrl });
});

router.get('/auth/:provider/callback', async (c) => {
  const provider = c.req.param('provider');
  const code = c.req.query('code');

  if (!code) {
    throw new Error('Authorization code is required');
  }

  const redirectUri = `${c.env.BASE_URL}/auth/${provider}/callback`;

  const { FacebookOAuthProvider, InstagramOAuthProvider } = await import('../lib/oauth');

  let oauthProvider;
  if (provider === 'facebook') {
    oauthProvider = new FacebookOAuthProvider(
      c.env.FACEBOOK_KEY,
      c.env.FACEBOOK_SECRET,
      redirectUri
    );
  } else {
    oauthProvider = new InstagramOAuthProvider(
      c.env.INSTAGRAM_KEY,
      c.env.INSTAGRAM_SECRET,
      redirectUri
    );
  }

  const accessToken = await oauthProvider.exchangeCodeForToken(code);
  const profile = await oauthProvider.getUserProfile(accessToken);

  let user = await c.env.DB.prepare(
    'SELECT * FROM users WHERE social_id = ? AND social_provider = ?'
  )
    .bind(profile.id, provider)
    .first();

  if (!user) {
    const result = await c.env.DB.prepare(
      `INSERT INTO users (social_id, social_provider, name, email, avatar_url, is_admin)
         VALUES (?, ?, ?, ?, ?, 0)`
    )
      .bind(profile.id, provider, profile.name, profile.email || '', profile.avatar_url || null)
      .run();

    user = await c.env.DB.prepare('SELECT * FROM users WHERE id = ?')
      .bind(result.meta.last_row_id)
      .first();
  }

  if (!user) {
    throw new Error('Failed to create user');
  }

  const { JWTManager } = await import('../lib/jwt');
  const jwtManager = new JWTManager(c.env.JWT_SECRET);

  const userData = {
    id: user.id as number,
    social_id: user.social_id as string,
    social_provider: user.social_provider as 'facebook' | 'instagram',
    name: user.name as string,
    email: user.email as string,
    avatar_url: user.avatar_url as string | undefined,
    is_admin: (user.is_admin as unknown as number) === 1,
    created_at: user.created_at as number,
    updated_at: user.updated_at as number,
  };

  const tokens = await jwtManager.generateTokens(userData);

  const { decodeJwt } = await import('jose');
  const decoded = decodeJwt(tokens.refresh);
  const expiresAt = new Date((((decoded.payload as any).exp as number) || 0) * 1000).toISOString();

  await c.env.DB.prepare(
    `INSERT INTO refresh_tokens (user_id, token, expires_at, created_at)
     VALUES (?, ?, ?, datetime('now'))`
  )
    .bind(user.id, tokens.refresh, expiresAt)
    .run();

  const { SessionManager } = await import('../lib/session');
  const sessionManager = new SessionManager(c.env.DB);
  const session = await sessionManager.createSession(user.id as number, { ...userData });

  return c.json({
    user: userData,
    tokens: {
      access: tokens.access,
      refresh: tokens.refresh,
    },
    session_id: session.id,
  });
});

router.post('/auth/refresh', async (c) => {
  const body = await c.req.json();
  const refreshToken = body.refresh_token;

  if (!refreshToken) {
    throw new Error('Refresh token is required');
  }

  const { JWTManager } = await import('../lib/jwt');
  const jwtManager = new JWTManager(c.env.JWT_SECRET);
  const payload = await jwtManager.verifyRefreshToken(refreshToken);

  const tokenRecord = await c.env.DB.prepare(
    'SELECT * FROM refresh_tokens WHERE token = ? AND expires_at > datetime("now")'
  )
    .bind(refreshToken)
    .first();

  if (!tokenRecord) {
    throw new Error('Invalid token');
  }

  const user = await c.env.DB.prepare('SELECT * FROM users WHERE id = ?')
    .bind(payload.user_id)
    .first();

  if (!user) {
    throw new Error('User not found');
  }

  const userData = {
    id: user.id as number,
    social_id: user.social_id as string,
    social_provider: user.social_provider as 'facebook' | 'instagram',
    name: user.name as string,
    email: user.email as string,
    avatar_url: user.avatar_url as string | undefined,
    is_admin: (user.is_admin as unknown as number) === 1,
    created_at: user.created_at as number,
    updated_at: user.updated_at as number,
  };

  const newTokens = await jwtManager.generateTokens(userData);

  await c.env.DB.prepare('DELETE FROM refresh_tokens WHERE token = ?').bind(refreshToken).run();

  const { decodeJwt } = await import('jose');
  const decoded = decodeJwt(newTokens.refresh);
  const expiresAt = new Date((((decoded.payload as any).exp as number) || 0) * 1000).toISOString();

  await c.env.DB.prepare(
    `INSERT INTO refresh_tokens (user_id, token, expires_at, created_at)
     VALUES (?, ?, ?, datetime('now'))`
  )
    .bind(user.id, newTokens.refresh, expiresAt)
    .run();

  return c.json({
    tokens: {
      access: newTokens.access,
      refresh: newTokens.refresh,
    },
  });
});

router.post('/auth/logout', async (c) => {
  const body = await c.req.json();
  const refreshToken = body.refresh_token;
  const sessionId = body.session_id;

  if (refreshToken) {
    await c.env.DB.prepare('DELETE FROM refresh_tokens WHERE token = ?').bind(refreshToken).run();
  }

  if (sessionId) {
    await c.env.DB.prepare('DELETE FROM sessions WHERE id = ?').bind(sessionId).run();
  }

  return c.json({ message: 'Logged out successfully' });
});

router.get('/auth/me', async (c) => {
  const user = c.get('user' as never);

  if (!user) {
    throw new Error('Authentication required');
  }

  return c.json({ user });
});

export default router;
