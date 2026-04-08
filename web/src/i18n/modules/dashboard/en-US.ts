// Dashboard translations - English (US)
import type { DashboardMessages } from './zh-Hans';

const messages: DashboardMessages = {
  welcome: 'Welcome back',
  welcomeDescription: "Here's an overview of your AI-powered workflows and integrations",
  newWorkflow: 'New Workflow',
  stats: {
    activeWorkflows: 'Active Workflows',
    apiCalls: 'API Calls Today',
    successRate: 'Success Rate',
    responseTime: 'Avg. Response Time',
    vsLastMonth: 'vs last month',
    last24Hours: 'Last 24 hours',
    withinSla: 'Within SLA',
    pendingApproval: '{count} pending approval',
    avgPerHour: 'Avg. {count} per hour',
  },
  charts: {
    apiCallsTitle: 'API Calls This Week',
    responseTimeTitle: 'Response Time Trend',
    ms: 'ms',
  },
  activity: {
    title: 'System Activity',
    viewAll: 'View All',
  },
  workflows: {
    title: 'Active Workflows',
    id: 'ID',
    name: 'Name',
    status: 'Status',
    calls: 'API Calls',
    success: 'Success Rate',
  },
  actions: {
    title: 'Quick Actions',
    createWorkflow: 'Create Workflow',
    createWorkflowDesc: 'Build a new AI pipeline',
    manageUsers: 'Manage Users',
    manageUsersDesc: 'Team member permissions',
    settings: 'Settings',
    settingsDesc: 'Configure your console',
  }
};

export default messages;
