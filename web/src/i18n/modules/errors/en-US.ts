// Errors translations - English (US)
import type { ErrorsMessages } from './zh-Hans';

const messages: ErrorsMessages = {
  notFound: 'Page not found',
  serverError: 'Server error',
  networkError: 'Network error, please check your connection',
  unauthorized: 'Please login to continue',
  forbidden: "You don't have permission to access this resource",
  unexpected: 'An unexpected error occurred',
  error: 'Error',
  sessionExpiredLoginAgain: 'Session expired, please login again',
  permissionDenied: 'You do not have permission to perform this action',
  resourceNotFound: 'The requested resource was not found',
  serverTryLater: 'Something went wrong on our server. Please try again later',
  tokenExpired: 'Your session has expired. Please log in again.',
  errorCode: 'Error Code: {code}',
};

export default messages;
