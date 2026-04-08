// Errors translations - Simplified Chinese
const messages = {
  notFound: '页面未找到',
  serverError: '服务器错误',
  networkError: '网络错误，请检查网络连接',
  unauthorized: '请登录后继续',
  forbidden: '您没有权限访问此资源',
};

export default messages;

export type ErrorsMessages = typeof messages;
