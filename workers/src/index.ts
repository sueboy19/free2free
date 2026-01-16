import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { logger } from 'hono/logger';
import { errorHandler } from './middleware/error';
import type { Env } from './types';
import authRoutes from './routes/auth';
import adminRoutes from './routes/admin';
import userRoutes from './routes/user';
import organizerRoutes from './routes/organizer';
import reviewRoutes from './routes/review';

const app = new Hono<{ Bindings: Env }>();

app.use('*', logger());
app.use('*', async (c, next) => {
  const corsMiddleware = cors({
    origin: c.env.CORS_ORIGINS.split(','),
    credentials: true,
    allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS'],
    allowHeaders: ['Content-Type', 'Authorization'],
  });
  return corsMiddleware(c, next);
});
app.use('*', errorHandler);

app.get('/', (c) => {
  return c.json({
    status: 'ok',
    message: 'Free2Free API is running',
    timestamp: new Date().toISOString(),
  });
});

app.route('/', authRoutes);
app.route('/', adminRoutes);
app.route('/', userRoutes);
app.route('/', organizerRoutes);
app.route('/', reviewRoutes);

export default app;
