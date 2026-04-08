'use client';

import { useState } from 'react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogBody,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import type {
  ApiProject,
  CreateProjectRequest,
  ProjectPlatform,
  ProjectStatus,
  UpdateProjectRequest,
} from '@/types/project';

// 项目平台选项。
// 作用：统一维护项目表单可选的平台枚举与展示文案。
export const PLATFORM_OPTIONS: Array<{ value: ProjectPlatform; label: string }> = [
  { value: 'go', label: 'Go' },
  { value: 'javascript', label: 'JavaScript' },
  { value: 'python', label: 'Python' },
  { value: 'java', label: 'Java' },
  { value: 'ruby', label: 'Ruby' },
  { value: 'php', label: 'PHP' },
  { value: 'csharp', label: 'C#' },
];

export type ProjectFormMode = 'create' | 'edit';

interface ProjectFormDraft {
  name: string;
  slug: string;
  platform: string;
  status: `${ProjectStatus}`;
}

// 项目表单默认值生成器。
// 作用：根据传入项目生成创建/编辑弹窗的初始草稿。
const getDefaultDraft = (project?: ApiProject | null): ProjectFormDraft => ({
  name: project?.name ?? '',
  slug: project?.slug ?? '',
  platform: project?.platform ?? '',
  status: String(project?.status ?? 1) as `${ProjectStatus}`,
});

export const resolvePlatformLabel = (platform: string) =>
  PLATFORM_OPTIONS.find((option) => option.value === platform)?.label || 'Not set';

// 项目状态文案解析器。
// 作用：把数字状态转换成界面可读的标签文本。
const resolveStatusLabel = (status: number) => (status === 1 ? 'Active' : 'Inactive');

/**
 * 项目状态徽章。
 * 作用：统一展示项目启用/停用状态，避免多个页面各自维护文案与颜色。
 */
export function ProjectStatusBadge({
  status,
}: {
  status: ProjectStatus;
}) {
  return (
    <Badge variant={status === 1 ? 'default' : 'secondary'}>
      {resolveStatusLabel(status)}
    </Badge>
  );
}

/**
 * 项目表单弹窗。
 * 作用：
 * 1. 在创建模式下收集新项目字段，调用 `POST /v1/projects`
 * 2. 在编辑模式下修改已有项目字段，调用 `PATCH /v1/projects/:id`
 */
export function ProjectFormDialog({
  open,
  mode,
  project,
  isSubmitting,
  onOpenChange,
  onSubmit,
}: {
  open: boolean;
  mode: ProjectFormMode;
  project?: ApiProject | null;
  isSubmitting: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (payload: CreateProjectRequest | UpdateProjectRequest) => Promise<void>;
}) {
  const formKey = `${mode}-${project?.id ?? 'new'}-${open ? 'open' : 'closed'}`;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <ProjectFormDialogBody
        key={formKey}
        mode={mode}
        project={project}
        isSubmitting={isSubmitting}
        onOpenChange={onOpenChange}
        onSubmit={onSubmit}
      />
    </Dialog>
  );
}

function ProjectFormDialogBody({
  mode,
  project,
  isSubmitting,
  onOpenChange,
  onSubmit,
}: {
  mode: ProjectFormMode;
  project?: ApiProject | null;
  isSubmitting: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (payload: CreateProjectRequest | UpdateProjectRequest) => Promise<void>;
}) {
  const [draft, setDraft] = useState<ProjectFormDraft>(() => getDefaultDraft(project));
  const [errors, setErrors] = useState<{ name?: string; slug?: string }>({});

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    const nextErrors: { name?: string; slug?: string } = {};
    const trimmedName = draft.name.trim();
    const trimmedSlug = draft.slug.trim();

    if (!trimmedName) {
      nextErrors.name = 'Project name is required.';
    }

    if (trimmedSlug.length > 50) {
      nextErrors.slug = 'Slug must be 50 characters or fewer.';
    }

    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors);
      return;
    }

    if (mode === 'create') {
      await onSubmit({
        name: trimmedName,
        slug: trimmedSlug || undefined,
        platform: (draft.platform || undefined) as ProjectPlatform | undefined,
      });
      return;
    }

    await onSubmit({
      name: trimmedName,
      platform: (draft.platform || undefined) as ProjectPlatform | undefined,
      status: Number(draft.status) as ProjectStatus,
    });
  };

  return (
    <DialogContent size="default">
      <DialogHeader>
        <DialogTitle>{mode === 'create' ? 'Create Project' : 'Edit Project'}</DialogTitle>
        <DialogDescription>
          {mode === 'create'
            ? 'Create a new project through POST /v1/projects.'
            : 'Update project fields through PATCH /v1/projects/:id.'}
        </DialogDescription>
      </DialogHeader>

      <DialogBody>
        <form id="project-form" onSubmit={handleSubmit} className="space-y-4 py-1">
          <div className="space-y-2">
            <Label htmlFor="project-name">Project Name</Label>
            <Input
              id="project-name"
              value={draft.name}
              onChange={(event) => setDraft((current) => ({ ...current, name: event.target.value }))}
              placeholder="Payments Platform"
              errorText={errors.name}
              root
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="project-slug">Slug</Label>
            <Input
              id="project-slug"
              value={draft.slug}
              onChange={(event) => setDraft((current) => ({ ...current, slug: event.target.value }))}
              placeholder={mode === 'create' ? 'Optional, auto-generated when left blank' : ''}
              readOnly={mode === 'edit'}
              disabled={mode === 'edit'}
              errorText={errors.slug}
              root
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="project-platform">Platform</Label>
            <Select
              value={draft.platform || 'none'}
              onValueChange={(value) =>
                setDraft((current) => ({ ...current, platform: value === 'none' ? '' : value }))
              }
            >
              <SelectTrigger id="project-platform" className="w-full">
                <SelectValue placeholder="Select a platform" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">Not set</SelectItem>
                {PLATFORM_OPTIONS.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {mode === 'edit' ? (
            <div className="space-y-2">
              <Label htmlFor="project-status">Status</Label>
              <Select
                value={draft.status}
                onValueChange={(value) =>
                  setDraft((current) => ({ ...current, status: value as `${ProjectStatus}` }))
                }
              >
                <SelectTrigger id="project-status" className="w-full">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="1">Active</SelectItem>
                  <SelectItem value="0">Inactive</SelectItem>
                </SelectContent>
              </Select>
            </div>
          ) : null}
        </form>
      </DialogBody>

      <DialogFooter>
        <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
          Cancel
        </Button>
        <Button type="submit" form="project-form" loading={isSubmitting}>
          {mode === 'create' ? 'Create Project' : 'Save Changes'}
        </Button>
      </DialogFooter>
    </DialogContent>
  );
}

/**
 * 项目删除确认弹窗。
 * 作用：
 * 1. 明确提示用户删除是不可逆操作
 * 2. 在用户确认后触发 `DELETE /v1/projects/:id`
 */
export function DeleteProjectDialog({
  open,
  project,
  isDeleting,
  onOpenChange,
  onConfirm,
}: {
  open: boolean;
  project?: ApiProject | null;
  isDeleting: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => Promise<void>;
}) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent size="sm">
        <DialogHeader>
          <DialogTitle>Delete Project</DialogTitle>
          <DialogDescription>
            This will permanently delete {project ? `"${project.name}"` : 'the selected project'}.
          </DialogDescription>
        </DialogHeader>

        <DialogBody>
          <Alert variant="destructive">
            <AlertTitle>Irreversible action</AlertTitle>
            <AlertDescription>
              The backend deletes the project record immediately through DELETE /v1/projects/:id.
            </AlertDescription>
          </Alert>
        </DialogBody>

        <DialogFooter>
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button type="button" variant="destructive" loading={isDeleting} onClick={() => void onConfirm()}>
            Delete Project
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
