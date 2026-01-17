import { Hono } from 'hono';
import type { Env } from '../types';
import { FacebookOAuthProvider, InstagramOAuthProvider } from '../lib/oauth';
import { JWTManager } from '../lib/jwt';
import { SessionManager } from '../lib/session';
import { decodeJwt } from 'jose';

const router = new Hono<{ Bindings: Env }>();

router.get('/auth/:provider', async (c) => {
  const provider = c.req.param('provider');

  if (!['facebook', 'instagram'].includes(provider)) {
    throw new Error('Invalid OAuth provider');
  }

  const baseUrl = c.env.BASE_URL || `https://${c.req.url.split('/')[2]}`;
  const redirectUri = `${baseUrl}/auth/${provider}/callback`;

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

  return c.redirect(authUrl, 302);
});

router.get('/auth/:provider/callback', async (c) => {
  const provider = c.req.param('provider');
  const code = c.req.query('code');

  if (!code) {
    throw new Error('Authorization code is required');
  }

  const baseUrl = c.env.BASE_URL || `https://${c.req.url.split('/')[2]}`;
  const redirectUri = `${baseUrl}/auth/${provider}/callback`;

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

  const decoded = decodeJwt(tokens.refresh);
  const expiresAt = new Date((((decoded.payload as any).exp as number) || 0) * 1000).toISOString();

  await c.env.DB.prepare(
    `INSERT INTO refresh_tokens (user_id, token, expires_at, created_at)
     VALUES (?, ?, ?, datetime('now'))`
  )
    .bind(user.id, tokens.refresh, expiresAt)
    .run();

  const sessionManager = new SessionManager(c.env.DB);
  await sessionManager.createSession(user.id as number, { ...userData });

  // Return HTML page with postMessage like Go backend
  const userJSON = JSON.stringify({
    id: userData.id,
    social_id: userData.social_id,
    social_provider: userData.social_provider,
    name: userData.name,
    email: userData.email,
    avatar_url: userData.avatar_url,
    is_admin: userData.is_admin,
    created_at: userData.created_at,
    updated_at: userData.updated_at,
  });

  const html = `<!DOCTYPE html>
<html>
<head>
  <title>${provider === 'facebook' ? 'Facebook' : 'Instagram'} 登入成功</title>
  <script>
    (function() {
      var response = {
        type: 'auth_success',
        payload: {
          user: ${userJSON},
          token: "${tokens.access}"
        }
      };

      if (window.opener) {
        window.opener.postMessage(response, '*');
      }

      setTimeout(function() {
        window.close();
      }, 1000);
    })();
  </script>
</head>
<body>
  <p>登入成功，正在返回...</p>
</body>
</html>`;

  c.header('Content-Type', 'text/html; charset=utf-8');
  return c.body(html);
});

router.post('/auth/refresh', async (c) => {
  const body = await c.req.json();
  const refreshToken = body.refresh_token;

  if (!refreshToken) {
    throw new Error('Refresh token is required');
  }

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
