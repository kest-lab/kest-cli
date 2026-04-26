// Errors translations - Simplified Chinese
const messages = {
  notFound: '页面未找到',
  serverError: '服务器错误',
  networkError: '网络错误，请检查网络连接',
  unauthorized: '请登录后继续',
  forbidden: '您没有权限访问此资源',
  unexpected: '发生了意外错误',
  error: '错误',
  sessionExpiredLoginAgain: '会话已过期，请重新登录',
  permissionDenied: '您没有权限执行此操作',
  resourceNotFound: '请求的资源不存在',
  serverTryLater: '服务器发生错误，请稍后重试',
  tokenExpired: '您的会话已过期，请重新登录。',
  errorCode: '错误代码：{code}',
};

export default messages;

export type ErrorsMessages = typeof messages;
