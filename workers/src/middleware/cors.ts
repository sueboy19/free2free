import { cors } from 'hono/cors';
import type { Env } from '../types';

export const corsMiddleware = (env: Env) => {
  return cors({
    origin: env.CORS_ORIGINS.split(','),
    credentials: true,
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS'],
    allowHeaders: ['Content-Type', 'Authorization'],
    exposeHeaders: ['Content-Length'],
  });
};
