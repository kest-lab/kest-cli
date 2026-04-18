'use client';

import Link from 'next/link';
import { useDeferredValue, useMemo, useState } from 'react';
import {
  ArrowLeft,
  Crown,
  Mail,
  Pencil,
  RefreshCw,
  Search,
  ShieldCheck,
  Trash2,
  UserPlus,
  Users,
} from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
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
import { StatCard, StatCardSkeleton } from '@/components/features/console/dashboard-stats';
import { ActionMenu, type ActionMenuItem } from '@/components/features/project/action-menu';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { buildApiPath } from '@/config/api';
import { buildProjectDetailRoute } from '@/constants/routes';
import {
  useCreateProjectMember,
  useDeleteProjectMember,
  useProjectMemberRole,
  useProjectMembers,
  useUpdateProjectMember,
} from '@/hooks/use-members';
import { useProject, useProjectStats } from '@/hooks/use-projects';
import { useUserSearch } from '@/hooks/use-users';
import type { ApiUser } from '@/types/auth';
import {
  PROJECT_MEMBER_ASSIGNABLE_ROLES,
  canEditProjectMember,
  canManageProjectMembers,
  canRemoveProjectMember,
  getProjectMemberRoleLabel,
  sortProjectMembers,
  type AssignableProjectMemberRole,
  type ProjectMember,
  type ProjectMemberRole,
} from '@/types/member';
import { formatDate } from '@/utils';

const EMPTY_MEMBERS: ProjectMember[] = [];
const ROLE_FILTER_OPTIONS: Array<{ value: 'all' | ProjectMemberRole; label: string }> = [
  { value: 'all', label: 'All roles' },
  { value: 'owner', label: 'Owner' },
  { value: 'admin', label: 'Admin' },
  { value: 'write', label: 'Write' },
  { value: 'read', label: 'Read' },
];

const getRoleBadgeVariant = (role: ProjectMemberRole) => {
  if (role === 'owner') {
    return 'default';
  }

  if (role === 'admin') {
    return 'secondary';
  }

  return 'outline';
};

const getMemberInitials = (member: Pick<ProjectMember, 'username' | 'email'>) => {
  const source = member.username.trim() || member.email.trim();
  const parts = source.split(/\s+/).filter(Boolean);

  if (parts.length >= 2) {
    return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
  }

  return source.slice(0, 2).toUpperCase();
};

const getAssignableRole = (role?: ProjectMemberRole): AssignableProjectMemberRole => {
  if (role === 'admin' || role === 'write' || role === 'read') {
    return role;
  }

  return 'read';
};

function RoleBadge({ role }: { role?: ProjectMemberRole }) {
  return (
    <Badge variant="outline" className="border-primary/20 bg-primary/10 text-primary">
      Role: {getProjectMemberRoleLabel(role)}
    </Badge>
  );
}

function MembersTableSkeleton() {
  return (
    <div className="space-y-3">
      {Array.from({ length: 5 }).map((_, index) => (
        <div key={index} className="h-14 animate-pulse rounded-xl border bg-muted/40" />
      ))}
    </div>
  );
}

export function ProjectMemberManagementPage({ projectId }: { projectId: number }) {
  const [searchQuery, setSearchQuery] = useState('');
  const [roleFilter, setRoleFilter] = useState<'all' | ProjectMemberRole>('all');
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [candidateQuery, setCandidateQuery] = useState('');
  const [selectedCandidate, setSelectedCandidate] = useState<ApiUser | null>(null);
  const [newMemberRole, setNewMemberRole] = useState<AssignableProjectMemberRole>('read');
  const [addDialogError, setAddDialogError] = useState<string | null>(null);
  const [editingMember, setEditingMember] = useState<ProjectMember | null>(null);
  const [editingRole, setEditingRole] = useState<AssignableProjectMemberRole>('read');
  const [deleteTarget, setDeleteTarget] = useState<ProjectMember | null>(null);

  const deferredSearchQuery = useDeferredValue(searchQuery.trim().toLowerCase());
  const deferredCandidateQuery = useDeferredValue(candidateQuery.trim());

  const projectQuery = useProject(projectId);
  const projectStatsQuery = useProjectStats(projectId);
  const membersQuery = useProjectMembers(projectId);
  const memberRoleQuery = useProjectMemberRole(projectId);
  const userSearchQuery = useUserSearch(deferredCandidateQuery, 20);
  const createMemberMutation = useCreateProjectMember(projectId);
  const updateMemberMutation = useUpdateProjectMember(projectId);
  const deleteMemberMutation = useDeleteProjectMember(projectId);

  const project = projectQuery.data;
  const currentRole = memberRoleQuery.data?.role;
  const currentUserId = memberRoleQuery.data?.user_id;
  const canManageMembers = canManageProjectMembers(currentRole);
  const members = useMemo(
    () => sortProjectMembers(membersQuery.data ?? EMPTY_MEMBERS),
    [membersQuery.data]
  );
  const memberUserIds = useMemo(() => new Set(members.map((member) => member.user_id)), [members]);
  const filteredMembers = useMemo(() => {
    return members.filter((member) => {
      const matchesRole = roleFilter === 'all' || member.role === roleFilter;
      const matchesKeyword =
        !deferredSearchQuery ||
        member.username.toLowerCase().includes(deferredSearchQuery) ||
        member.email.toLowerCase().includes(deferredSearchQuery);

      return matchesRole && matchesKeyword;
    });
  }, [deferredSearchQuery, members, roleFilter]);
  const candidateResults = useMemo(
    () =>
      (userSearchQuery.data ?? []).filter((candidate) => !memberUserIds.has(candidate.id)),
    [memberUserIds, userSearchQuery.data]
  );
  const ownerCount = members.filter((member) => member.role === 'owner').length;
  const adminCount = members.filter((member) => member.role === 'admin').length;
  const writeCount = members.filter((member) => member.role === 'write').length;
  const readCount = members.filter((member) => member.role === 'read').length;
  const membersPath = buildApiPath('/projects/:id/members');
  const membersMePath = buildApiPath('/projects/:id/members/me');
  const isRefreshing =
    projectQuery.isFetching ||
    projectStatsQuery.isFetching ||
    membersQuery.isFetching ||
    memberRoleQuery.isFetching;
  const hasLoadError =
    !membersQuery.isLoading &&
    (Boolean(membersQuery.error) || Boolean(memberRoleQuery.error));

  const resetAddDialog = () => {
    setCandidateQuery('');
    setSelectedCandidate(null);
    setNewMemberRole('read');
    setAddDialogError(null);
  };

  const handleRefresh = async () => {
    await Promise.all([
      projectQuery.refetch(),
      projectStatsQuery.refetch(),
      membersQuery.refetch(),
      memberRoleQuery.refetch(),
    ]);
  };

  const handleOpenAddDialog = () => {
    resetAddDialog();
    setIsAddDialogOpen(true);
  };

  const handleAddMember = async () => {
    if (!selectedCandidate) {
      setAddDialogError('Select a user before adding a member.');
      return;
    }

    try {
      await createMemberMutation.mutateAsync({
        user_id: selectedCandidate.id,
        role: newMemberRole,
      });
      setIsAddDialogOpen(false);
      resetAddDialog();
    } catch {
      // Global HTTP error handling already surfaces failure feedback.
    }
  };

  const handleOpenEditDialog = (member: ProjectMember) => {
    setEditingMember(member);
    setEditingRole(getAssignableRole(member.role));
  };

  const handleUpdateMember = async () => {
    if (!editingMember) {
      return;
    }

    try {
      await updateMemberMutation.mutateAsync({
        userId: editingMember.user_id,
        data: {
          role: editingRole,
        },
      });
      setEditingMember(null);
    } catch {
      // Global HTTP error handling already surfaces failure feedback.
    }
  };

  const handleDeleteMember = async () => {
    if (!deleteTarget) {
      return;
    }

    try {
      await deleteMemberMutation.mutateAsync(deleteTarget.user_id);
      setDeleteTarget(null);
    } catch {
      // Global HTTP error handling already surfaces failure feedback.
    }
  };

  const headerActionItems: ActionMenuItem[] = [
    {
      key: 'members-refresh',
      label: isRefreshing ? 'Refreshing...' : 'Refresh',
      icon: RefreshCw,
      disabled: isRefreshing,
      onSelect: () => {
        void handleRefresh();
      },
    },
    {
      key: 'members-add',
      label: 'Add Member',
      icon: UserPlus,
      disabled: !canManageMembers,
      onSelect: handleOpenAddDialog,
    },
  ];

  return (
    <>
      <main className="h-full min-h-0 overflow-y-auto">
        <div className="space-y-8 p-6 pt-6">
          <div className="relative overflow-hidden rounded-xl border border-primary/10 bg-linear-to-r from-primary/10 via-cyan-500/5 to-transparent p-6 transition-colors duration-500">
            <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxwYXRoIGQ9Ik0xOCAxOGgyNHYyNEgxOHoiIHN0cm9rZT0iY3VycmVudENvbG9yIiBzdHJva2Utb3BhY2l0eT0iLjA1Ii8+PC9nPjwvc3ZnPg==')] opacity-50" />
            <div className="relative flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
              <div className="space-y-3">
                <Button asChild variant="link" className="h-auto px-0 text-sm text-muted-foreground">
                  <Link href={buildProjectDetailRoute(projectId)}>
                    <ArrowLeft className="h-4 w-4" />
                    Back to Project Overview
                  </Link>
                </Button>

                <div className="space-y-2">
                  <div className="flex flex-wrap items-center gap-2">
                    <h1 className="text-3xl font-bold tracking-tight">Members</h1>
                    <Users className="h-6 w-6 text-primary" />
                    <RoleBadge role={currentRole} />
                  </div>
                  <p className="max-w-4xl text-sm text-text-muted">
                    Manage project access through
                    {' '}
                    <code>{membersPath}</code>
                    {' '}
                    and resolve the current user role through
                    {' '}
                    <code>{membersMePath}</code>
                    .
                  </p>
                </div>

                {project ? (
                  <div className="flex flex-wrap items-center gap-2">
                    <Badge variant="outline" className="font-mono">
                      {project.slug}
                    </Badge>
                    <Badge variant="outline">{project.name}</Badge>
                    <Badge variant="outline">{projectStatsQuery.data?.member_count ?? members.length} members</Badge>
                  </div>
                ) : null}
              </div>

              <div className="flex flex-wrap items-center gap-3">
                <Button type="button" onClick={handleOpenAddDialog} disabled={!canManageMembers}>
                  <UserPlus className="h-4 w-4" />
                  Add Member
                </Button>
                <ActionMenu
                  items={headerActionItems}
                  ariaLabel="Open member management actions"
                  triggerVariant="outline"
                />
              </div>
            </div>
          </div>

          {!canManageMembers && memberRoleQuery.isSuccess ? (
            <Alert>
              <ShieldCheck className="h-4 w-4" />
              <AlertTitle>Read-only member access</AlertTitle>
              <AlertDescription>
                当前角色是
                {' '}
                <strong>{getProjectMemberRoleLabel(currentRole)}</strong>
                ，可以查看成员列表，但只有
                {' '}
                <strong>admin</strong>
                {' '}
                和
                {' '}
                <strong>owner</strong>
                {' '}
                可以新增、调整或移除成员。
              </AlertDescription>
            </Alert>
          ) : null}

          <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
            {membersQuery.isLoading ? (
              <>
                <StatCardSkeleton />
                <StatCardSkeleton />
                <StatCardSkeleton />
                <StatCardSkeleton />
              </>
            ) : (
              <>
                <StatCard
                  title="Total Members"
                  value={projectStatsQuery.data?.member_count ?? members.length}
                  description="All users with project access"
                  icon={Users}
                />
                <StatCard
                  title="Admins & Owners"
                  value={ownerCount + adminCount}
                  description={`${ownerCount} owner${ownerCount === 1 ? '' : 's'}, ${adminCount} admin${adminCount === 1 ? '' : 's'}`}
                  icon={ShieldCheck}
                  variant="warning"
                />
                <StatCard
                  title="Writers"
                  value={writeCount}
                  description="Can edit project resources"
                  icon={Pencil}
                  variant="success"
                />
                <StatCard
                  title="Readers"
                  value={readCount}
                  description="Browse-only project access"
                  icon={Mail}
                  variant="default"
                />
              </>
            )}
          </div>

          <Card>
            <CardHeader className="space-y-4">
              <div className="flex flex-col gap-4 xl:flex-row xl:items-start xl:justify-between">
                <div>
                  <CardTitle>Project Members</CardTitle>
                  <CardDescription>
                    Search by username or email, then adjust operational roles for non-owner members.
                  </CardDescription>
                </div>
                <div className="flex flex-col gap-3 sm:flex-row">
                  <div className="relative min-w-[240px]">
                    <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted" />
                    <Input
                      value={searchQuery}
                      onChange={(event) => setSearchQuery(event.target.value)}
                      placeholder="Filter by username or email"
                      className="pl-9"
                    />
                  </div>
                  <Select
                    value={roleFilter}
                    onValueChange={(value) => setRoleFilter(value as 'all' | ProjectMemberRole)}
                  >
                    <SelectTrigger className="w-[180px]">
                      <SelectValue placeholder="Filter by role" />
                    </SelectTrigger>
                    <SelectContent>
                      {ROLE_FILTER_OPTIONS.map((option) => (
                        <SelectItem key={option.value} value={option.value}>
                          {option.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              {membersQuery.isLoading ? (
                <MembersTableSkeleton />
              ) : hasLoadError ? (
                <Alert>
                  <AlertTitle>Unable to load members</AlertTitle>
                  <AlertDescription>
                    The project members page could not load its data. Retry the request or confirm the current user still has access to this project.
                  </AlertDescription>
                </Alert>
              ) : (
                <div className="overflow-hidden rounded-xl border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>User</TableHead>
                        <TableHead>Role</TableHead>
                        <TableHead>Joined</TableHead>
                        <TableHead>Updated</TableHead>
                        <TableHead className="text-right">Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {filteredMembers.map((member) => {
                        const isCurrentUser = currentUserId !== undefined && member.user_id === currentUserId;
                        const canEdit = canEditProjectMember(member, currentRole, currentUserId);
                        const canRemove = canRemoveProjectMember(member, currentRole, currentUserId);
                        const rowActionItems: ActionMenuItem[] = [
                          {
                            key: `edit-${member.user_id}`,
                            label: 'Edit Role',
                            icon: Pencil,
                            disabled: !canEdit,
                            onSelect: () => handleOpenEditDialog(member),
                          },
                          {
                            key: `delete-${member.user_id}`,
                            label: 'Remove',
                            icon: Trash2,
                            destructive: true,
                            separatorBefore: true,
                            disabled: !canRemove,
                            onSelect: () => setDeleteTarget(member),
                          },
                        ];

                        return (
                          <TableRow key={member.user_id}>
                            <TableCell>
                              <div className="flex items-center gap-3">
                                <Avatar className="h-10 w-10">
                                  <AvatarFallback>{getMemberInitials(member)}</AvatarFallback>
                                </Avatar>
                                <div className="space-y-1">
                                  <div className="flex flex-wrap items-center gap-2">
                                    <span className="font-medium">{member.username}</span>
                                    {isCurrentUser ? <Badge variant="outline">You</Badge> : null}
                                    {member.role === 'owner' ? (
                                      <Badge variant="outline" className="gap-1">
                                        <Crown className="h-3 w-3" />
                                        Owner
                                      </Badge>
                                    ) : null}
                                  </div>
                                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                                    <Mail className="h-3.5 w-3.5" />
                                    <span>{member.email}</span>
                                  </div>
                                </div>
                              </div>
                            </TableCell>
                            <TableCell>
                              <Badge variant={getRoleBadgeVariant(member.role)}>
                                {getProjectMemberRoleLabel(member.role)}
                              </Badge>
                            </TableCell>
                            <TableCell>{formatDate(member.created_at, 'YYYY-MM-DD HH:mm')}</TableCell>
                            <TableCell>{formatDate(member.updated_at, 'YYYY-MM-DD HH:mm')}</TableCell>
                            <TableCell className="text-right">
                              {!canManageMembers ? (
                                <span className="text-sm text-muted-foreground">Admin required</span>
                              ) : canEdit || canRemove ? (
                                <ActionMenu
                                  items={rowActionItems}
                                  ariaLabel={`Open actions for ${member.username}`}
                                />
                              ) : (
                                <span className="text-sm text-muted-foreground">Protected</span>
                              )}
                            </TableCell>
                          </TableRow>
                        );
                      })}
                      {filteredMembers.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={5} className="py-10 text-center text-muted-foreground">
                            No members match the current filter.
                          </TableCell>
                        </TableRow>
                      ) : null}
                    </TableBody>
                  </Table>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </main>

      <Dialog
        open={isAddDialogOpen}
        onOpenChange={(open) => {
          setIsAddDialogOpen(open);
          if (!open) {
            resetAddDialog();
          }
        }}
      >
        <DialogContent size="lg">
          <DialogHeader>
            <DialogTitle>Add Member</DialogTitle>
            <DialogDescription>
              Search existing users, then grant them project access with an assignable role.
            </DialogDescription>
          </DialogHeader>
          <DialogBody>
            <div className="space-y-5">
              {addDialogError ? (
                <Alert variant="destructive">
                  <AlertTitle>Cannot add member</AlertTitle>
                  <AlertDescription>{addDialogError}</AlertDescription>
                </Alert>
              ) : null}

              <div className="space-y-2">
                <Label htmlFor="member-search">Find user</Label>
                <Input
                  id="member-search"
                  value={candidateQuery}
                  onChange={(event) => {
                    setCandidateQuery(event.target.value);
                    setAddDialogError(null);
                  }}
                  placeholder="Search by username or email"
                />
              </div>

              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <Label>Matching users</Label>
                  <span className="text-sm text-muted-foreground">
                    {userSearchQuery.isFetching ? 'Searching…' : `${candidateResults.length} available`}
                  </span>
                </div>
                <div className="max-h-72 space-y-2 overflow-y-auto rounded-xl border p-3">
                  {deferredCandidateQuery.length === 0 ? (
                    <p className="text-sm text-muted-foreground">
                      Start typing to search existing users.
                    </p>
                  ) : userSearchQuery.isFetching ? (
                    <div className="space-y-2">
                      {Array.from({ length: 3 }).map((_, index) => (
                        <div key={index} className="h-14 animate-pulse rounded-xl bg-muted/50" />
                      ))}
                    </div>
                  ) : candidateResults.length === 0 ? (
                    <p className="text-sm text-muted-foreground">
                      No eligible users matched this query. Matching users who are already project members are hidden.
                    </p>
                  ) : (
                    candidateResults.map((candidate) => {
                      const isSelected = selectedCandidate?.id === candidate.id;

                      return (
                        <button
                          key={candidate.id}
                          type="button"
                          className="flex w-full items-center justify-between rounded-xl border px-3 py-3 text-left transition-colors hover:border-primary/40 hover:bg-primary/5"
                          onClick={() => {
                            setSelectedCandidate(candidate);
                            setAddDialogError(null);
                          }}
                        >
                          <div className="min-w-0">
                            <div className="font-medium">{candidate.username}</div>
                            <div className="truncate text-sm text-muted-foreground">{candidate.email}</div>
                          </div>
                          {isSelected ? <Badge>Selected</Badge> : <Badge variant="outline">Select</Badge>}
                        </button>
                      );
                    })
                  )}
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="member-role">Role</Label>
                <Select
                  value={newMemberRole}
                  onValueChange={(value) => setNewMemberRole(value as AssignableProjectMemberRole)}
                >
                  <SelectTrigger id="member-role">
                    <SelectValue placeholder="Select role" />
                  </SelectTrigger>
                  <SelectContent>
                    {PROJECT_MEMBER_ASSIGNABLE_ROLES.map((role) => (
                      <SelectItem key={role} value={role}>
                        {getProjectMemberRoleLabel(role)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {selectedCandidate ? (
                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4">
                  <div className="text-sm font-medium">Selected user</div>
                  <div className="mt-1 text-sm text-muted-foreground">
                    {selectedCandidate.username}
                    {' '}
                    ·
                    {' '}
                    {selectedCandidate.email}
                  </div>
                </div>
              ) : null}
            </div>
          </DialogBody>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setIsAddDialogOpen(false)}>
              Cancel
            </Button>
            <Button
              type="button"
              onClick={() => {
                void handleAddMember();
              }}
              disabled={createMemberMutation.isPending || !canManageMembers}
            >
              <UserPlus className="h-4 w-4" />
              Add Member
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={Boolean(editingMember)} onOpenChange={(open) => !open && setEditingMember(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Member Role</DialogTitle>
            <DialogDescription>
              Update the project role for
              {' '}
              <strong>{editingMember?.username}</strong>
              .
            </DialogDescription>
          </DialogHeader>
          <DialogBody>
            <div className="space-y-4">
              <div className="rounded-xl border p-4">
                <div className="font-medium">{editingMember?.username}</div>
                <div className="text-sm text-muted-foreground">{editingMember?.email}</div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="edit-member-role">Role</Label>
                <Select
                  value={editingRole}
                  onValueChange={(value) => setEditingRole(value as AssignableProjectMemberRole)}
                >
                  <SelectTrigger id="edit-member-role">
                    <SelectValue placeholder="Select role" />
                  </SelectTrigger>
                  <SelectContent>
                    {PROJECT_MEMBER_ASSIGNABLE_ROLES.map((role) => (
                      <SelectItem key={role} value={role}>
                        {getProjectMemberRoleLabel(role)}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </DialogBody>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setEditingMember(null)}>
              Cancel
            </Button>
            <Button
              type="button"
              onClick={() => {
                void handleUpdateMember();
              }}
              disabled={updateMemberMutation.isPending || !editingMember}
            >
              Save Role
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={Boolean(deleteTarget)} onOpenChange={(open) => !open && setDeleteTarget(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Remove Member</DialogTitle>
            <DialogDescription>
              Remove
              {' '}
              <strong>{deleteTarget?.username}</strong>
              {' '}
              from this project.
            </DialogDescription>
          </DialogHeader>
          <DialogBody>
            <Alert>
              <AlertTitle>Access will be revoked immediately</AlertTitle>
              <AlertDescription>
                {deleteTarget?.username}
                {' '}
                will no longer be able to access project resources after this action completes.
              </AlertDescription>
            </Alert>
            {deleteTarget ? (
              <div className="mt-4 rounded-xl border p-4">
                <div className="font-medium">{deleteTarget.username}</div>
                <div className="text-sm text-muted-foreground">{deleteTarget.email}</div>
              </div>
            ) : null}
          </DialogBody>
          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => setDeleteTarget(null)}>
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={() => {
                void handleDeleteMember();
              }}
              disabled={deleteMemberMutation.isPending || !deleteTarget}
            >
              <Trash2 className="h-4 w-4" />
              Remove Member
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}
