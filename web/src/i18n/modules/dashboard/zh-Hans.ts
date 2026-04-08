// Dashboard translations - Simplified Chinese
const messages = {
  welcome: '欢迎回来',
  welcomeDescription: '这里是您的 AI 工作流和集成概览',
  newWorkflow: '新工作流',
  stats: {
    activeWorkflows: '活跃工作流',
    apiCalls: '今日 API 调用',
    successRate: '成功率',
    responseTime: '平均响应时间',
    vsLastMonth: '较上月',
    last24Hours: '最近 24 小时',
    withinSla: '符合 SLA',
    pendingApproval: '{count} 个待审批',
    avgPerHour: '平均每小时 {count} 次',
  },
  charts: {
    apiCallsTitle: '本周 API 调用',
    responseTimeTitle: '典型响应时间趋势',
    ms: '毫秒',
  },
  activity: {
    title: '系统活动',
    viewAll: '查看全部',
  },
  workflows: {
    title: '活跃工作流',
    id: 'ID',
    name: '名称',
    status: '状态',
    calls: 'API 调用',
    success: '成功率',
  },
  actions: {
    title: '快捷操作',
    createWorkflow: '创建工作流',
    createWorkflowDesc: '构建新的 AI 管道',
    manageUsers: '管理用户',
    manageUsersDesc: '团队成员权限',
    settings: '设置',
    settingsDesc: '配置您的控制台',
  }
};

export default messages;

export type DashboardMessages = typeof messages;
