export interface OAuthProvider {
  name: 'facebook' | 'instagram';
  getAuthUrl(): string;
  exchangeCodeForToken(code: string): Promise<string>;
  getUserProfile(accessToken: string): Promise<OAuthProfile>;
}

export interface OAuthProfile {
  id: string;
  name: string;
  email: string;
  avatar_url?: string;
}

export class FacebookOAuthProvider implements OAuthProvider {
  name = 'facebook' as const;
  private clientId: string;
  private clientSecret: string;
  private redirectUri: string;

  constructor(clientId: string, clientSecret: string, redirectUri: string) {
    this.clientId = clientId;
    this.clientSecret = clientSecret;
    this.redirectUri = redirectUri;
  }

  getAuthUrl(): string {
    const params = new URLSearchParams({
      client_id: this.clientId,
      redirect_uri: this.redirectUri,
      scope: 'email,public_profile',
      response_type: 'code',
    });
    return `https://www.facebook.com/v18.0/dialog/oauth?${params}`;
  }

  async exchangeCodeForToken(code: string): Promise<string> {
    const params = new URLSearchParams({
      client_id: this.clientId,
      client_secret: this.clientSecret,
      redirect_uri: this.redirectUri,
      code,
    });

    const response = await fetch(`https://graph.facebook.com/v18.0/oauth/access_token?${params}`);
    const data: any = await response.json();

    if (data.error) {
      throw new Error(data.error.message);
    }

    return data.access_token;
  }

  async getUserProfile(accessToken: string): Promise<OAuthProfile> {
    const params = new URLSearchParams({
      fields: 'id,name,email,picture',
      access_token: accessToken,
    });

    const response = await fetch(`https://graph.facebook.com/v18.0/me?${params}`);
    const data: any = await response.json();

    if (data.error) {
      throw new Error(data.error.message);
    }

    return {
      id: data.id,
      name: data.name,
      email: data.email,
      avatar_url: data.picture?.data?.url,
    };
  }
}

export class InstagramOAuthProvider implements OAuthProvider {
  name = 'instagram' as const;
  private clientId: string;
  private clientSecret: string;
  private redirectUri: string;

  constructor(clientId: string, clientSecret: string, redirectUri: string) {
    this.clientId = clientId;
    this.clientSecret = clientSecret;
    this.redirectUri = redirectUri;
  }

  getAuthUrl(): string {
    const params = new URLSearchParams({
      client_id: this.clientId,
      redirect_uri: this.redirectUri,
      scope: 'user_profile',
      response_type: 'code',
    });
    return `https://api.instagram.com/oauth/authorize?${params}`;
  }

  async exchangeCodeForToken(code: string): Promise<string> {
    const response = await fetch('https://api.instagram.com/oauth/access_token', {
      method: 'POST',
      body: JSON.stringify({
        client_id: this.clientId,
        client_secret: this.clientSecret,
        grant_type: 'authorization_code',
        redirect_uri: this.redirectUri,
        code,
      }),
    });

    const data: any = await response.json();

    if (data.error) {
      throw new Error(data.error.message);
    }

    return data.access_token;
  }

  async getUserProfile(accessToken: string): Promise<OAuthProfile> {
    const response = await fetch('https://graph.instagram.com/me', {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });

    const data: any = await response.json();

    if (data.error) {
      throw new Error(data.error.message);
    }

    return {
      id: data.id,
      name: data.username,
      email: '',
      avatar_url: '',
    };
  }
}

export class OAuthManager {
  private providers: Map<string, OAuthProvider> = new Map();

  registerProvider(provider: OAuthProvider) {
    this.providers.set(provider.name, provider);
  }

  getProvider(name: string): OAuthProvider | undefined {
    return this.providers.get(name);
  }
}
