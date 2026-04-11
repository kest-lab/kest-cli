'use client';

import Link from 'next/link';
import { startTransition, useDeferredValue, useEffect, useMemo, useState } from 'react';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useQueryClient } from '@tanstack/react-query';
import {
  ArrowRight,
  FolderKanban,
  Plus,
  Search,
} from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Separator } from '@/components/ui/separator';
import {
  ProjectFormDialog,
  type ProjectFormMode,
  ProjectStatusBadge,
} from '@/components/features/project/project-shared';
import {
  buildProjectApiSpecsRoute,
  buildProjectCollectionsRoute,
  buildProjectDetailRoute,
  buildProjectEnvironmentsRoute,
  buildProjectTestCasesRoute,
} from '@/constants/routes';
import { apiSpecKeys, useApiSpecs } from '@/hooks/use-api-specs';
import {
  projectKeys,
  useCreateProject,
  useProjectStats,
  useProjects,
  useUpdateProject,
} from '@/hooks/use-projects';
import { apiSpecService } from '@/services/api-spec';
import { projectService } from '@/services/project';
import type { ApiProject, CreateProjectRequest, UpdateProjectRequest } from '@/types/project';
import { formatDate } from '@/utils';

const PROJECTS_PAGE_SIZE = 1000;
const MAX_PREVIEW_SPECS = 5;
const EMPTY_PROJECTS: ApiProject[] = [];

const getProjectCreatedAt = (project: ApiProject) => project.created_at || '';

const sortProjectsByCreatedAtDesc = (left: ApiProject, right: ApiProject) =>
  getProjectCreatedAt(right).localeCompare(getProjectCreatedAt(left));

const parseProjectId = (value: string | null) => {
  const numericValue = Number(value);
  return Number.isInteger(numericValue) && numericValue > 0 ? numericValue : null;
};

const formatProjectTimestamp = (value?: string | null) => {
  if (!value) {
    return null;
  }

  const parsed = new Date(value);
  return Number.isNaN(parsed.getTime()) ? null : formatDate(value, 'YYYY-MM-DD HH:mm');
};

const buildDashboardHref = (
  pathname: string,
  searchParams: URLSearchParams,
  previewProjectId?: number | null
) => {
  const nextParams = new URLSearchParams(searchParams.toString());

  if (previewProjectId) {
    nextParams.set('preview', String(previewProjectId));
  } else {
    nextParams.delete('preview');
  }

  const queryString = nextParams.toString();
  return queryString ? `${pathname}?${queryString}` : pathname;
};

export function ProjectDashboardPage() {
  const pathname = usePathname();
  const router = useRouter();
  const searchParams = useSearchParams();
  const queryClient = useQueryClient();

  const [searchQuery, setSearchQuery] = useState('');
  const [formMode, setFormMode] = useState<ProjectFormMode>('create');
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingProject, setEditingProject] = useState<ApiProject | null>(null);
  const [hasAutoSelectedProject, setHasAutoSelectedProject] = useState(false);

  const deferredSearch = useDeferredValue(searchQuery);

  const projectsQuery = useProjects({ page: 1, perPage: PROJECTS_PAGE_SIZE });
  const createProjectMutation = useCreateProject();
  const updateProjectMutation = useUpdateProject();

  const projects = projectsQuery.data?.items ?? EMPTY_PROJECTS;
  const previewProjectId = parseProjectId(searchParams.get('preview'));

  const filteredProjects = useMemo(() => {
    const normalizedQuery = deferredSearch.trim().toLowerCase();

    if (!normalizedQuery) {
      return [...projects].sort(sortProjectsByCreatedAtDesc);
    }

    return projects
      .filter((project) =>
        [project.name, project.slug, project.platform]
          .filter(Boolean)
          .some((value) => value.toLowerCase().includes(normalizedQuery))
      )
      .sort(sortProjectsByCreatedAtDesc);
  }, [deferredSearch, projects]);

  const selectedProject =
    previewProjectId !== null
      ? projects.find((project) => project.id === previewProjectId) ?? null
      : null;
  const fallbackProject = useMemo(
    () => [...projects].sort(sortProjectsByCreatedAtDesc)[0] ?? null,
    [projects]
  );

  const prefetchProjectPreview = (projectId: number) => {
    void queryClient.prefetchQuery({
      queryKey: projectKeys.projectStats(projectId),
      queryFn: () => projectService.getStats(projectId),
    });

    void queryClient.prefetchQuery({
      queryKey: apiSpecKeys.list({ projectId, page: 1, pageSize: MAX_PREVIEW_SPECS }),
      queryFn: () =>
        apiSpecService.list({
          projectId,
          page: 1,
          pageSize: MAX_PREVIEW_SPECS,
        }),
    });
  };

  const navigateToPreview = (projectId?: number | null) => {
    if (projectId) {
      prefetchProjectPreview(projectId);
    }

    startTransition(() => {
      router.replace(
        buildDashboardHref(pathname, new URLSearchParams(searchParams.toString()), projectId)
      );
    });
  };

  useEffect(() => {
    filteredProjects.slice(0, 3).forEach((project) => {
      prefetchProjectPreview(project.id);
    });
  }, [filteredProjects, queryClient]);

  useEffect(() => {
    if (
      hasAutoSelectedProject ||
      projectsQuery.isLoading ||
      previewProjectId !== null ||
      !fallbackProject
    ) {
      return;
    }

    setHasAutoSelectedProject(true);
    navigateToPreview(fallbackProject.id);
  }, [fallbackProject, hasAutoSelectedProject, previewProjectId, projectsQuery.isLoading]);

  const openCreateDialog = () => {
    setFormMode('create');
    setEditingProject(null);
    setIsFormOpen(true);
  };

  const openEditDialog = (project: ApiProject) => {
    setFormMode('edit');
    setEditingProject(project);
    setIsFormOpen(true);
  };

  const handleProjectSubmit = async (payload: CreateProjectRequest | UpdateProjectRequest) => {
    try {
      if (formMode === 'create') {
        const project = await createProjectMutation.mutateAsync(payload as CreateProjectRequest);
        navigateToPreview(project.id);
      } else if (editingProject) {
        await updateProjectMutation.mutateAsync({
          id: editingProject.id,
          data: payload as UpdateProjectRequest,
        });
        navigateToPreview(editingProject.id);
      }

      setIsFormOpen(false);
      setEditingProject(null);
    } catch {
      // Global HTTP error handling already surfaces failure feedback.
    }
  };

  return (
    <div className="flex h-full min-h-0 flex-col overflow-hidden lg:flex-row">
      <aside className="w-full shrink-0 border-b border-border/60 bg-bg-surface/70 lg:w-[296px] lg:border-b-0 lg:border-r">
        <div className="flex h-full max-h-[42vh] flex-col overflow-hidden lg:max-h-none">
          <div className="space-y-4 p-4">
            <div className="rounded-2xl border border-primary/10 bg-linear-to-br from-primary/10 via-primary/5 to-transparent p-4">
              <div className="flex items-center justify-between gap-3">
                <div>
                  <p className="text-xs font-medium uppercase tracking-[0.18em] text-text-muted">
                    Project Dashboard
                  </p>
                  <h1 className="mt-2 text-xl font-semibold tracking-tight">All projects</h1>
                  <p className="mt-1 text-sm text-text-muted">
                    Preview projects here, then enter the scoped workspace when you are ready to work.
                  </p>
                </div>
                <div className="rounded-2xl bg-background/70 p-3 text-primary shadow-sm">
                  <FolderKanban className="h-5 w-5" />
                </div>
              </div>
            </div>

            <div className="relative">
              <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted" />
              <Input
                value={searchQuery}
                onChange={(event) => setSearchQuery(event.target.value)}
                placeholder="Search projects"
                className="pl-9"
              />
            </div>

            <Button type="button" onClick={openCreateDialog} className="w-full">
              <Plus className="h-4 w-4" />
              Create Project
            </Button>
          </div>

          <Separator />

          <div className="min-h-0 flex-1 overflow-y-auto p-3">
            <div className="mb-3 flex items-center justify-between px-2 text-xs font-medium uppercase tracking-[0.18em] text-text-muted">
              <span>Projects</span>
              <span>{filteredProjects.length}</span>
            </div>

            {projectsQuery.isLoading ? (
              <div className="space-y-2">
                {Array.from({ length: 6 }).map((_, index) => (
                  <div
                    key={index}
                    className="rounded-2xl border border-border/60 bg-background/60 p-3"
                  >
                    <div className="h-4 w-24 animate-pulse rounded bg-muted" />
                    <div className="mt-2 h-3 w-40 animate-pulse rounded bg-muted" />
                    <div className="mt-3 h-3 w-20 animate-pulse rounded bg-muted" />
                  </div>
                ))}
              </div>
            ) : projectsQuery.error ? (
              <Alert>
                <AlertTitle>Unable to load projects</AlertTitle>
                <AlertDescription>
                  The dashboard could not load the project list from the current API.
                </AlertDescription>
              </Alert>
            ) : filteredProjects.length === 0 ? (
              <Card className="border-dashed">
                <CardContent className="p-6 text-sm text-text-muted">
                  {projects.length === 0
                    ? 'No projects are available yet. Create the first project to populate the dashboard.'
                    : 'No projects match the current search keyword.'}
                </CardContent>
              </Card>
            ) : (
              <div className="space-y-2">
                {filteredProjects.map((project) => {
                  const isActive = project.id === selectedProject?.id;

                  return (
                    <button
                      key={project.id}
                      type="button"
                      onClick={() => navigateToPreview(project.id)}
                      onMouseEnter={() => prefetchProjectPreview(project.id)}
                      onFocus={() => prefetchProjectPreview(project.id)}
                      onTouchStart={() => prefetchProjectPreview(project.id)}
                      aria-pressed={isActive}
                      className={`group w-full rounded-2xl border p-3 text-left transition-colors ${
                        isActive
                          ? 'border-primary/30 bg-primary/10 shadow-sm'
                          : 'border-transparent bg-background/60 hover:border-border/60 hover:bg-background'
                      }`}
                    >
                      <div className="flex items-start justify-between gap-3">
                        <div className="min-w-0">
                          <p className="truncate text-sm font-medium text-text-main">{project.name}</p>
                          <p className="truncate text-xs text-text-muted">{project.slug}</p>
                        </div>
                        <ProjectStatusBadge status={project.status} />
                      </div>
                      <div className="mt-3 flex flex-wrap gap-2 text-xs text-text-muted">
                        {isActive ? (
                          <Badge variant="outline" className="border-primary/20 bg-primary/10 text-primary">
                            Selected
                          </Badge>
                        ) : null}
                        {project.created_at ? (
                          <span>Created {formatDate(project.created_at, 'YYYY-MM-DD')}</span>
                        ) : null}
                        {!isActive ? (
                          <span className="transition-opacity group-hover:opacity-100 lg:opacity-0">
                            Preview
                          </span>
                        ) : null}
                      </div>
                    </button>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      </aside>

      <main className="min-h-0 min-w-0 flex-1 overflow-y-auto">
        <div className="space-y-6 p-4 md:p-6">
          {selectedProject ? (
            <ProjectPreviewPanel
              project={selectedProject}
              onEdit={() => openEditDialog(selectedProject)}
            />
          ) : (
            <ProjectDashboardWelcome
              projects={projects}
              onOpenProject={navigateToPreview}
              onCreateProject={openCreateDialog}
            />
          )}
        </div>
      </main>

      <ProjectFormDialog
        open={isFormOpen}
        mode={formMode}
        project={editingProject}
        isSubmitting={createProjectMutation.isPending || updateProjectMutation.isPending}
        onOpenChange={setIsFormOpen}
        onSubmit={handleProjectSubmit}
      />
    </div>
  );
}

function ProjectDashboardWelcome({
  projects,
  onOpenProject,
  onCreateProject,
}: {
  projects: ApiProject[];
  onOpenProject: (projectId?: number | null) => void;
  onCreateProject: () => void;
}) {
  const recentProjects = [...projects]
    .sort(sortProjectsByCreatedAtDesc)
    .slice(0, 5);

  return (
    <div className="space-y-6">
      <Card className="overflow-hidden border-primary/10 bg-linear-to-br from-primary/10 via-transparent to-transparent">
        <CardHeader>
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="space-y-2">
              <Badge variant="outline" className="border-primary/20 bg-primary/10 text-primary">
                Dashboard overview
              </Badge>
              <CardTitle className="text-2xl tracking-tight">
                Choose a project and continue where work happens
              </CardTitle>
              <CardDescription className="max-w-3xl">
                The left sidebar is your project inventory. The right panel is a launchpad:
                inspect the project, then jump into AI-assisted API Specs, Environments, or Test Cases.
              </CardDescription>
            </div>

            <Button type="button" onClick={onCreateProject}>
              <Plus className="h-4 w-4" />
              Create Project
            </Button>
          </div>
        </CardHeader>
      </Card>

      <div className="grid gap-6 xl:grid-cols-[1.1fr_0.9fr]">
        <Card className="border-border/60">
          <CardHeader>
            <CardTitle>Recent projects</CardTitle>
            <CardDescription>
              Start by previewing one of the projects below.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {recentProjects.length === 0 ? (
              <Alert>
                <AlertTitle>No projects yet</AlertTitle>
                <AlertDescription>
                  Use the create action to seed the dashboard with your first project.
                </AlertDescription>
              </Alert>
            ) : (
              recentProjects.map((project) => (
                <button
                  key={project.id}
                  type="button"
                  onClick={() => onOpenProject(project.id)}
                  className="flex w-full items-center justify-between rounded-2xl border border-border/60 bg-background/70 px-4 py-3 text-left transition-colors hover:border-primary/20 hover:bg-background"
                >
                  <div className="min-w-0">
                    <p className="truncate text-sm font-medium">{project.name}</p>
                    <p className="truncate text-xs text-text-muted">
                      {project.slug}
                    </p>
                  </div>
                  <ArrowRight className="h-4 w-4 text-text-muted" />
                </button>
              ))
            )}
          </CardContent>
        </Card>

        <Card className="border-dashed border-border/70">
          <CardHeader>
            <CardTitle>Recommended workflow</CardTitle>
            <CardDescription>
              Keep the user journey simple: define interfaces, set runtime context, then validate behavior.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-text-muted">
            <div className="rounded-2xl border border-border/60 bg-background/70 p-4">
              <p className="font-medium text-text-main">1. Define interfaces</p>
              <p className="mt-1">Start in API Specs and use AI draft to turn product intent into a structured endpoint.</p>
            </div>
            <div className="rounded-2xl border border-border/60 bg-background/70 p-4">
              <p className="font-medium text-text-main">2. Configure runtime context</p>
              <p className="mt-1">Add Environments and Collections once requests need real endpoints and reusable drafts.</p>
            </div>
            <div className="rounded-2xl border border-border/60 bg-background/70 p-4">
              <p className="font-medium text-text-main">3. Validate and iterate</p>
              <p className="mt-1">Move to Test Cases after the spec and environment baseline is clear.</p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function ProjectPreviewPanel({
  project,
  onEdit,
}: {
  project: ApiProject;
  onEdit: () => void;
}) {
  const [isSlowPreview, setIsSlowPreview] = useState(false);
  const statsQuery = useProjectStats(project.id);
  const apiSpecsQuery = useApiSpecs({
    projectId: project.id,
    page: 1,
    pageSize: MAX_PREVIEW_SPECS,
  });

  const projectDetail = project;
  const apiSpecs = apiSpecsQuery.data?.items ?? [];
  const stats = statsQuery.data;
  const apiSpecCount = stats?.api_spec_count ?? apiSpecs.length;
  const environmentCount = stats?.environment_count ?? 0;
  const memberCount = stats?.member_count ?? 0;
  const createdAtLabel = formatProjectTimestamp(projectDetail.created_at);
  const hasResolvedReadiness = Boolean(stats);
  const hasReadinessError =
    !hasResolvedReadiness &&
    Boolean(statsQuery.error);
  const nextStep = hasResolvedReadiness
    ? resolveDashboardNextStep({
        projectId: project.id,
        apiSpecCount,
        environmentCount,
      })
    : null;
  const readinessItems = hasResolvedReadiness
    ? ([
        {
          label: 'Source of truth',
          value:
            apiSpecCount > 0
              ? `${apiSpecCount} spec${apiSpecCount === 1 ? '' : 's'} ready`
              : 'Missing',
          tone: apiSpecCount > 0 ? 'ready' : 'pending',
          detail:
            apiSpecCount > 0
              ? 'The interface inventory exists and can drive docs and tests.'
              : 'Start in API Specs so the project has a stable interface inventory.',
        },
        {
          label: 'Runtime context',
          value:
            environmentCount > 0
              ? `${environmentCount} environment${environmentCount === 1 ? '' : 's'} configured`
              : 'Missing',
          tone: environmentCount > 0 ? 'ready' : 'pending',
          detail:
            environmentCount > 0
              ? 'Base URLs, headers, and variables are ready for execution.'
              : 'Add one environment before running requests against real targets.',
        },
        {
          label: 'Validation',
          value:
            apiSpecCount > 0 && environmentCount > 0 ? 'Ready to generate' : 'Waiting on setup',
          tone: apiSpecCount > 0 && environmentCount > 0 ? 'ready' : 'pending',
          detail:
            apiSpecCount > 0 && environmentCount > 0
              ? 'Move into Test Cases when you want the first runnable coverage.'
              : 'Test generation becomes useful after the spec and environment baseline exist.',
        },
      ] as const)
    : [];

  const handleRetryPreview = () => {
    void statsQuery.refetch();
    void apiSpecsQuery.refetch();
  };

  useEffect(() => {
    if (hasResolvedReadiness || hasReadinessError) {
      setIsSlowPreview(false);
      return;
    }

    const timeoutId = window.setTimeout(() => {
      setIsSlowPreview(true);
    }, 2500);

    return () => {
      window.clearTimeout(timeoutId);
    };
  }, [hasReadinessError, hasResolvedReadiness]);

  return (
    <div className="space-y-6">
      <Card className="overflow-hidden border-border/60 bg-linear-to-r from-background via-background to-primary/5">
        <CardHeader className="space-y-4">
          <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
            <div className="space-y-3">
              <div className="flex flex-wrap items-center gap-2">
                <Badge variant="outline" className="border-primary/20 bg-primary/10 text-primary">
                  Selected project
                </Badge>
                <ProjectStatusBadge status={projectDetail.status} />
              </div>
              <div>
                <CardTitle className="text-2xl tracking-tight">{projectDetail.name}</CardTitle>
                <CardDescription className="mt-2 max-w-3xl text-sm leading-6">
                  {hasReadinessError
                    ? 'Some project signals failed to load. Retry this preview or open the workspace directly.'
                    : isSlowPreview
                      ? 'Loading project status is taking longer than usual. You can retry or open the workspace directly.'
                    : nextStep?.summary ??
                      'Loading project status so the next step reflects real project data.'}
                </CardDescription>
              </div>
              <div className="flex flex-wrap items-center gap-5 text-sm text-text-muted">
                <span>Slug: {projectDetail.slug}</span>
                {createdAtLabel ? <span>Created {createdAtLabel}</span> : null}
                {typeof stats?.member_count === 'number' ? (
                  <span>
                    {memberCount} team member{memberCount === 1 ? '' : 's'}
                  </span>
                ) : null}
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={onEdit}
                  className="h-auto px-0 text-sm text-text-muted hover:bg-transparent hover:text-text-main"
                >
                  Edit details
                </Button>
              </div>
            </div>

            <div className="flex flex-wrap gap-2">
              <Button asChild variant="outline">
                <Link href={buildProjectDetailRoute(project.id)}>
                  Open workspace
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </CardHeader>
      </Card>

      <div className="grid gap-6 xl:grid-cols-[1.05fr_0.95fr]">
        <Card className="border-border/60">
          <CardHeader>
            <CardTitle>Progress</CardTitle>
            <CardDescription>Three checks are enough to understand project readiness.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {hasReadinessError ? (
              <Alert>
                <AlertTitle>Unable to load project readiness</AlertTitle>
                <AlertDescription className="mt-2">
                  The preview could not determine whether the project is ready for setup or testing.
                </AlertDescription>
                <Button type="button" variant="outline" size="sm" className="mt-4" onClick={handleRetryPreview}>
                  Retry preview
                </Button>
              </Alert>
            ) : !hasResolvedReadiness ? (
              <>
                <div className="space-y-3">
                  {Array.from({ length: 3 }).map((_, index) => (
                    <div
                      key={index}
                      className="h-18 animate-pulse rounded-2xl border border-border/60 bg-muted/40"
                    />
                  ))}
                </div>
                {isSlowPreview ? (
                  <div className="flex items-center justify-between gap-3 rounded-2xl border border-dashed border-border/70 bg-background/70 px-4 py-3 text-sm text-text-muted">
                    <span>Still loading readiness signals.</span>
                    <Button type="button" variant="outline" size="sm" onClick={handleRetryPreview}>
                      Retry
                    </Button>
                  </div>
                ) : null}
              </>
            ) : (
              readinessItems.map((item) => (
                <div
                  key={item.label}
                  className="rounded-2xl border border-border/60 bg-background/70 p-4"
                >
                  <div className="flex items-center justify-between gap-4">
                    <div>
                      <p className="text-sm font-medium text-text-main">{item.label}</p>
                      <p className="mt-1 text-sm text-text-muted">{item.detail}</p>
                    </div>
                    <Badge
                      variant="outline"
                      className={
                        item.tone === 'ready'
                          ? 'border-emerald-200 bg-emerald-500/10 text-emerald-700'
                          : 'border-amber-200 bg-amber-500/10 text-amber-700'
                      }
                    >
                      {item.value}
                    </Badge>
                  </div>
                </div>
              ))
            )}

            {apiSpecsQuery.isLoading ? (
              <div className="space-y-2 pt-2">
                {Array.from({ length: 3 }).map((_, index) => (
                  <div key={index} className="h-14 animate-pulse rounded-2xl bg-muted" />
                ))}
              </div>
            ) : apiSpecs.length > 0 ? (
              <div className="space-y-3 rounded-2xl border border-border/60 bg-muted/20 p-4">
                <div className="flex items-center justify-between gap-3">
                  <div>
                    <p className="text-sm font-medium text-text-main">Recent API specs</p>
                    <p className="mt-1 text-sm text-text-muted">
                      Keep this preview short. Inspect the rest inside the workspace.
                    </p>
                  </div>
                  <Button asChild variant="ghost" className="px-0">
                    <Link href={buildProjectApiSpecsRoute(project.id)}>
                      Open API Specs
                      <ArrowRight className="h-4 w-4" />
                    </Link>
                  </Button>
                </div>

                <div className="space-y-2">
                  {apiSpecs.slice(0, 3).map((spec) => (
                    <div
                      key={spec.id}
                      className="flex items-start justify-between gap-3 rounded-2xl border border-border/60 bg-background/80 px-4 py-3"
                    >
                      <div className="min-w-0">
                        <p className="truncate text-sm font-medium">
                          {spec.method} {spec.path}
                        </p>
                        <p className="truncate text-xs text-text-muted">
                          {spec.summary || spec.description || 'No summary provided'}
                        </p>
                      </div>
                      <Badge variant="outline">{spec.version}</Badge>
                    </div>
                  ))}
                </div>
              </div>
            ) : null}
          </CardContent>
        </Card>

        <Card className="border-border/60 bg-linear-to-br from-primary/5 via-background to-background">
          <CardHeader>
            <CardTitle>
              {hasReadinessError ? 'Preview needs attention' : nextStep?.title ?? 'Loading next step'}
            </CardTitle>
            <CardDescription>
              {hasReadinessError
                ? 'The dashboard cannot recommend the next step until the preview data loads successfully.'
                : nextStep?.description ??
                  'Waiting for project readiness so this recommendation is not guessed.'}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {hasReadinessError ? (
              <Alert>
                <AlertTitle>Recommendation unavailable</AlertTitle>
                <AlertDescription className="mt-2">
                  Retry the preview, or continue in the workspace if you already know the next action.
                </AlertDescription>
                <div className="mt-4 flex flex-wrap gap-2">
                  <Button type="button" variant="outline" size="sm" onClick={handleRetryPreview}>
                    Retry preview
                  </Button>
                  <Button asChild size="sm">
                    <Link href={buildProjectDetailRoute(project.id)}>
                      Open workspace
                      <ArrowRight className="h-4 w-4" />
                    </Link>
                  </Button>
                </div>
              </Alert>
            ) : !nextStep ? (
              <>
                <div className="space-y-3">
                  <div className="h-28 animate-pulse rounded-2xl border border-border/60 bg-muted/40" />
                  <div className="h-40 animate-pulse rounded-2xl border border-border/60 bg-muted/40" />
                </div>
                {isSlowPreview ? (
                  <div className="flex items-center justify-between gap-3 rounded-2xl border border-dashed border-border/70 bg-background/70 px-4 py-3 text-sm text-text-muted">
                    <span>Recommendation is taking longer than usual.</span>
                    <Button type="button" variant="outline" size="sm" onClick={handleRetryPreview}>
                      Retry
                    </Button>
                  </div>
                ) : null}
              </>
            ) : (
              <>
                <div className="rounded-2xl border border-border/60 bg-background/80 p-4">
                  <p className="text-xs font-medium uppercase tracking-[0.18em] text-text-muted">
                    Why now
                  </p>
                  <p className="mt-3 text-sm leading-6 text-text-muted">{nextStep.reason}</p>
                </div>

                <div className="rounded-2xl border border-border/60 bg-background/80 p-4">
                  <p className="text-xs font-medium uppercase tracking-[0.18em] text-text-muted">
                    To unlock this step
                  </p>
                  <div className="mt-3 space-y-3">
                    {nextStep.blockers.map((blocker) => (
                      <div key={blocker.label} className="flex items-start justify-between gap-4">
                        <div>
                          <p className="text-sm font-medium text-text-main">{blocker.label}</p>
                          <p className="mt-1 text-sm text-text-muted">{blocker.detail}</p>
                        </div>
                        <Badge
                          variant="outline"
                          className={
                            blocker.tone === 'ready'
                              ? 'border-emerald-200 bg-emerald-500/10 text-emerald-700'
                              : 'border-amber-200 bg-amber-500/10 text-amber-700'
                          }
                        >
                          {blocker.state}
                        </Badge>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="flex flex-wrap gap-2">
                  <Button asChild>
                    <Link href={nextStep.primaryHref}>
                      {nextStep.primaryLabel}
                      <ArrowRight className="h-4 w-4" />
                    </Link>
                  </Button>
                  {nextStep.secondaryHref && nextStep.secondaryLabel ? (
                    <Button asChild variant="outline">
                      <Link href={nextStep.secondaryHref}>
                        {nextStep.secondaryLabel}
                        <ArrowRight className="h-4 w-4" />
                      </Link>
                    </Button>
                  ) : null}
                </div>
              </>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function resolveDashboardNextStep({
  projectId,
  apiSpecCount,
  environmentCount,
}: {
  projectId: number;
  apiSpecCount: number;
  environmentCount: number;
}) {
  if (apiSpecCount === 0) {
    return {
      summary: 'No API spec exists yet. Start with AI Draft before setting up secondary surfaces.',
      title: 'Define the first API surface',
      description: 'Describe one endpoint, create the first spec, then come back for runtime setup.',
      reason: 'The project needs one concrete interface before environments or validation become useful.',
      primaryHref: `${buildProjectApiSpecsRoute(projectId)}?ai=create`,
      primaryLabel: 'AI Draft API',
      blockers: [
        {
          label: 'API source of truth',
          detail: 'No API spec exists yet, so documentation and tests have nothing stable to build from.',
          state: 'Missing',
          tone: 'pending' as const,
        },
        {
          label: 'Runtime setup',
          detail: 'Environments can wait until the first endpoint is defined.',
          state: 'Not started',
          tone: 'pending' as const,
        },
      ],
    };
  }

  if (environmentCount === 0) {
    return {
      summary: `This project already has ${apiSpecCount} API spec${apiSpecCount === 1 ? '' : 's'}, but execution still has no runtime target.`,
      title: 'Add the first environment',
      description: 'Add one development or staging target so requests and tests can run somewhere real.',
      reason: 'You already defined the surface. One environment unlocks requests, examples, and the first test runs.',
      primaryHref: buildProjectEnvironmentsRoute(projectId),
      primaryLabel: 'Configure Environment',
      secondaryHref: buildProjectApiSpecsRoute(projectId),
      secondaryLabel: 'Review API Specs',
      blockers: [
        {
          label: 'Execution target',
          detail: 'No base URL, variables, or shared headers are configured yet.',
          state: 'Missing',
          tone: 'pending' as const,
        },
        {
          label: 'Spec baseline',
          detail: `The project already has ${apiSpecCount} API spec${apiSpecCount === 1 ? '' : 's'} ready for downstream work.`,
          state: 'Ready',
          tone: 'ready' as const,
        },
      ],
    };
  }

  return {
    summary: `The project has ${apiSpecCount} API spec${apiSpecCount === 1 ? '' : 's'} and ${environmentCount} runtime environment${environmentCount === 1 ? '' : 's'}. Move into validation instead of adding more dashboard detail.`,
    title: 'Generate validation coverage',
    description: 'Generate test cases from the existing specs, then use Collections only for debugging.',
    reason: 'The spec and environment baseline exists. The next real value comes from runnable coverage.',
    primaryHref: buildProjectTestCasesRoute(projectId),
    primaryLabel: 'Open Test Cases',
    secondaryHref: buildProjectCollectionsRoute(projectId),
    secondaryLabel: 'Open Collections',
    blockers: [
      {
        label: 'API source of truth',
        detail: 'API Specs are already in place and can drive generated coverage.',
        state: 'Ready',
        tone: 'ready' as const,
      },
      {
        label: 'Runtime setup',
        detail: `At least ${environmentCount} environment${environmentCount === 1 ? '' : 's'} is configured for execution.`,
        state: 'Ready',
        tone: 'ready' as const,
      },
    ],
  };
}
