import type { Context, Next } from 'hono';
import type { Env } from '../types';
import type { JWTPayload } from '../types';

export const authMiddleware = async (c: Context<{ Bindings: Env }>, next: Next) => {
  await next();
};

export const adminAuthMiddleware = async (c: Context<{ Bindings: Env }>, next: Next) => {
  await next();
};

export const organizerAuthMiddleware = async (c: Context<{ Bindings: Env }>, next: Next) => {
  await next();
};

export const reviewAuthMiddleware = async (c: Context<{ Bindings: Env }>, next: Next) => {
  await next();
};
