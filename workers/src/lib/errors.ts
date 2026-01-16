import { AppError } from '../middleware/error';

export const ErrorCodes = {
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  UNAUTHORIZED: 'UNAUTHORIZED',
  FORBIDDEN: 'FORBIDDEN',
  NOT_FOUND: 'NOT_FOUND',
  CONFLICT: 'CONFLICT',
  INTERNAL_ERROR: 'INTERNAL_ERROR',
  OAUTH_ERROR: 'OAUTH_ERROR',
  INVALID_TOKEN: 'INVALID_TOKEN',
  TOKEN_EXPIRED: 'TOKEN_EXPIRED',
} as const;

export const Errors = {
  validation: (message: string) => new AppError(400, ErrorCodes.VALIDATION_ERROR, message),

  unauthorized: (message: string = 'Authentication required') =>
    new AppError(401, ErrorCodes.UNAUTHORIZED, message),

  forbidden: (message: string = 'Access denied') =>
    new AppError(403, ErrorCodes.FORBIDDEN, message),

  notFound: (resource: string) => new AppError(404, ErrorCodes.NOT_FOUND, `${resource} not found`),

  conflict: (message: string) => new AppError(409, ErrorCodes.CONFLICT, message),

  internal: (message: string = 'Internal server error') =>
    new AppError(500, ErrorCodes.INTERNAL_ERROR, message),

  oauth: (message: string) => new AppError(500, ErrorCodes.OAUTH_ERROR, message),

  invalidToken: () => new AppError(401, ErrorCodes.INVALID_TOKEN, 'Invalid token'),

  tokenExpired: () => new AppError(401, ErrorCodes.TOKEN_EXPIRED, 'Token expired'),
};
