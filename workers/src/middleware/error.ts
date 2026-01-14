import type { Context, Next } from 'hono';

export class AppError extends Error {
  constructor(
    public statusCode: number,
    public code: string,
    message: string
  ) {
    super(message);
    this.name = 'AppError';
  }
}

export const errorHandler = async (c: Context, next: Next) => {
  try {
    await next();
  } catch (error) {
    console.error('Error:', error);

    if (error instanceof AppError) {
      return c.json(
        {
          code: error.statusCode,
          error: error.message,
          code_error: error.code,
        },
        error.statusCode as any
      );
    }

    return c.json(
      {
        code: 500,
        error: 'Internal server error',
        code_error: 'INTERNAL_ERROR',
      },
      500 as any
    );
  }
};
