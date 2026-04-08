'use client';

import type { LucideIcon } from 'lucide-react';
import { FileJson2, FolderGit2, FolderOpen, Globe, History, Tags } from 'lucide-react';
import {
  buildProjectApiSpecsRoute,
  buildProjectCategoriesRoute,
  buildProjectCollectionsRoute,
  buildProjectEnvironmentsRoute,
  buildProjectFlowsRoute,
  buildProjectHistoriesRoute,
} from '@/constants/routes';

export type ProjectWorkspaceModule =
  | 'collections'
  | 'api-specs'
  | 'environments'
  | 'categories'
  | 'histories'
  | 'flows';

export interface ProjectWorkspaceModuleMeta {
  value: ProjectWorkspaceModule;
  label: string;
  shortLabel: string;
  description: string;
  icon: LucideIcon;
}

export const PROJECT_WORKSPACE_MODULES: ProjectWorkspaceModuleMeta[] = [
  {
    value: 'collections',
    label: 'Collections',
    shortLabel: 'Collections',
    description: 'Browse Postman-style request groups and nested drafts.',
    icon: FolderOpen,
  },
  {
    value: 'api-specs',
    label: 'API Specs',
    shortLabel: 'Specs',
    description: 'Browse saved interface definitions and documentation.',
    icon: FileJson2,
  },
  {
    value: 'environments',
    label: 'Environments',
    shortLabel: 'Envs',
    description: 'Inspect runtime targets, variables, and headers.',
    icon: Globe,
  },
  {
    value: 'categories',
    label: 'Categories',
    shortLabel: 'Categories',
    description: 'Organize resources into nested project groupings.',
    icon: Tags,
  },
  {
    value: 'histories',
    label: 'History',
    shortLabel: 'History',
    description: 'Review activity and execution records scoped to this project.',
    icon: History,
  },
  {
    value: 'flows',
    label: 'Flows',
    shortLabel: 'Flows',
    description: 'Open reusable workflow assets and orchestration definitions.',
    icon: FolderGit2,
  },
];

export const getProjectWorkspaceModuleMeta = (module: ProjectWorkspaceModule) =>
  PROJECT_WORKSPACE_MODULES.find((item) => item.value === module) ??
  PROJECT_WORKSPACE_MODULES[0];

export const buildProjectWorkspaceRoute = (
  projectId: string | number,
  module: ProjectWorkspaceModule
) => {
  switch (module) {
    case 'collections':
      return buildProjectCollectionsRoute(projectId);
    case 'api-specs':
      return buildProjectApiSpecsRoute(projectId);
    case 'environments':
      return buildProjectEnvironmentsRoute(projectId);
    case 'categories':
      return buildProjectCategoriesRoute(projectId);
    case 'histories':
      return buildProjectHistoriesRoute(projectId);
    case 'flows':
      return buildProjectFlowsRoute(projectId);
    default:
      return buildProjectApiSpecsRoute(projectId);
  }
};
