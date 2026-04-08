// Common translations - Simplified Chinese
const messages = {
  loading: '加载中...',
  error: '发生错误',
  retry: '重试',
  retryLater: '请稍后重试',
  save: '保存',
  cancel: '取消',
  confirm: '确认',
  delete: '删除',
  edit: '编辑',
  create: '创建',
  search: '搜索',
  noData: '暂无数据',
  success: '成功',
  failed: '失败',
  
  // User CRUD messages
  userCreateSuccess: '用户创建成功',
  userCreateFailed: '创建用户失败',
  userUpdateSuccess: '用户更新成功',
  userUpdateFailed: '更新用户失败',
  userDeleteSuccess: '用户删除成功',
  userDeleteFailed: '删除用户失败',
  userCreated: '用户 {name} 已创建',
  userUpdated: '用户 {name} 已更新',
  
  // Date Picker
  year: '年',
  month: '月',
  hour: '时',
  minute: '分',
  second: '秒',
  now: '现在',
  datePlaceholder: '选择日期',
  toggleLanguage: '切换语言',
  selectLanguage: '选择语言',
};

export default messages;

// Export type for other locales to implement
export type CommonMessages = typeof messages;
