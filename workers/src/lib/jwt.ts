import { SignJWT, jwtVerify, type JWTPayload as JoseJWTPayload } from 'jose';
import type { JWTPayload, User } from '../types';

const JWT_ALGORITHM = 'HS256';
const ACCESS_TOKEN_EXPIRY = '15m';
const REFRESH_TOKEN_EXPIRY = '7d';

export class JWTManager {
  constructor(private secret: string) {
    if (secret.length < 32) {
      throw new Error('JWT_SECRET must be at least 32 characters');
    }
  }

  private getSecretKey(): Uint8Array {
    return new TextEncoder().encode(this.secret);
  }

  async generateAccessToken(user: User): Promise<string> {
    return new SignJWT({
      user_id: user.id,
      user_name: user.name,
      is_admin: user.is_admin,
    } as JWTPayload & JoseJWTPayload)
      .setProtectedHeader({ alg: JWT_ALGORITHM })
      .setIssuedAt()
      .setExpirationTime(ACCESS_TOKEN_EXPIRY)
      .sign(this.getSecretKey());
  }

  async generateRefreshToken(user: User): Promise<string> {
    return new SignJWT({
      user_id: user.id,
    })
      .setProtectedHeader({ alg: JWT_ALGORITHM })
      .setIssuedAt()
      .setExpirationTime(REFRESH_TOKEN_EXPIRY)
      .sign(this.getSecretKey());
  }

  async generateTokens(user: User): Promise<{ access: string; refresh: string }> {
    const [access, refresh] = await Promise.all([
      this.generateAccessToken(user),
      this.generateRefreshToken(user),
    ]);

    return { access, refresh };
  }

  async verifyAccessToken(token: string): Promise<JWTPayload> {
    const { payload } = await jwtVerify(token, this.getSecretKey());
    return payload as unknown as JWTPayload;
  }

  async verifyRefreshToken(token: string): Promise<{ user_id: number }> {
    const { payload } = await jwtVerify(token, this.getSecretKey());
    return payload as { user_id: number };
  }
}
