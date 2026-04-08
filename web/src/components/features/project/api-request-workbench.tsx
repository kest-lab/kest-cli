'use client';

import { startTransition, useDeferredValue, useMemo, useState } from 'react';
import {
  ChevronDown,
  ChevronRight,
  Copy,
  FolderOpen,
  MoreHorizontal,
  Plus,
  Save,
  Search,
  SendHorizonal,
  Trash2,
} from 'lucide-react';
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import { useDeleteCollection, useUpdateCollection } from '@/hooks/use-collections';
import { cn } from '@/utils';

type RequestMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
type RequestSection = 'params' | 'authorization' | 'headers' | 'body' | 'scripts' | 'settings';
type BulkMode = 'table' | 'bulk';
type AuthorizationMode = 'none' | 'bearer' | 'basic' | 'api-key';
type BodyMode = 'json' | 'raw' | 'form-data';

interface KeyValueRow {
  id: string;
  key: string;
  value: string;
  description: string;
}

interface ResponseDraft {
  status: number | null;
  statusLabel: string;
  durationMs: number | null;
  sizeBytes: number | null;
  body: string;
  error: string | null;
}

interface RequestPageTab {
  id: string;
  title: string;
  collectionId: string | null;
  method: RequestMethod;
  url: string;
  activeSection: RequestSection;
  paramsMode: BulkMode;
  paramsRows: KeyValueRow[];
  paramsBulk: string;
  authorizationMode: AuthorizationMode;
  authorizationValue: string;
  headersMode: BulkMode;
  headersRows: KeyValueRow[];
  headersBulk: string;
  bodyMode: BodyMode;
  bodyContent: string;
  scripts: string;
  settings: {
    followRedirects: boolean;
    strictTls: boolean;
    persistCookies: boolean;
  };
  response: ResponseDraft;
  isSending: boolean;
}

interface CollectionNode {
  id: string;
  name: string;
  color: string;
  requestIds: string[];
}

interface InitialWorkbenchState {
  tabs: RequestPageTab[];
  collections: CollectionNode[];
  activeTabId: string;
  activeCollectionId: string | null;
  expandedCollectionIds: string[];
  nextTabIndex: number;
}

const METHOD_OPTIONS: RequestMethod[] = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH'];
const ENVIRONMENT_OPTIONS = ['development', 'staging', 'production'] as const;
const SECTION_ITEMS: Array<{ value: RequestSection; label: string }> = [
  { value: 'params', label: 'Params' },
  { value: 'authorization', label: 'Authorization' },
  { value: 'headers', label: 'Headers' },
  { value: 'body', label: 'Body' },
  { value: 'scripts', label: 'Scripts' },
  { value: 'settings', label: 'Settings' },
];
const BODY_MODE_OPTIONS: BodyMode[] = ['json', 'raw', 'form-data'];
const AUTHORIZATION_OPTIONS: AuthorizationMode[] = ['none', 'bearer', 'basic', 'api-key'];
const COLLECTION_COLORS = ['#2563eb', '#0f766e', '#ea580c', '#7c3aed', '#dc2626'];
const METHOD_BADGE_STYLES: Record<RequestMethod, string> = {
  GET: 'border-emerald-500/30 bg-emerald-500/10 text-emerald-700 dark:text-emerald-300',
  POST: 'border-sky-500/30 bg-sky-500/10 text-sky-700 dark:text-sky-300',
  PUT: 'border-amber-500/30 bg-amber-500/10 text-amber-700 dark:text-amber-300',
  PATCH: 'border-violet-500/30 bg-violet-500/10 text-violet-700 dark:text-violet-300',
  DELETE: 'border-rose-500/30 bg-rose-500/10 text-rose-700 dark:text-rose-300',
};

const createLocalId = (prefix: string) =>
  `${prefix}-${Math.random().toString(36).slice(2, 8)}-${Date.now().toString(36)}`;

const createKeyValueRow = (key = '', value = '', description = ''): KeyValueRow => ({
  id: createLocalId('kv'),
  key,
  value,
  description,
});

const createEmptyResponse = (): ResponseDraft => ({
  status: null,
  statusLabel: '',
  durationMs: null,
  sizeBytes: null,
  body: '',
  error: null,
});

const createRequestPageTab = (
  index: number,
  overrides: Partial<RequestPageTab> = {}
): RequestPageTab => ({
  id: overrides.id ?? createLocalId('request-tab'),
  title: overrides.title ?? (index === 1 ? 'New Request' : `New Request ${index}`),
  collectionId: overrides.collectionId ?? null,
  method: overrides.method ?? 'GET',
  url: overrides.url ?? 'https://localhost:3000/health',
  activeSection: overrides.activeSection ?? 'params',
  paramsMode: overrides.paramsMode ?? 'table',
  paramsRows: overrides.paramsRows ?? [createKeyValueRow()],
  paramsBulk: overrides.paramsBulk ?? '',
  authorizationMode: overrides.authorizationMode ?? 'none',
  authorizationValue: overrides.authorizationValue ?? '',
  headersMode: overrides.headersMode ?? 'table',
  headersRows:
    overrides.headersRows ?? [createKeyValueRow('Accept', 'application/json', 'Default mock header')],
  headersBulk: overrides.headersBulk ?? 'Accept: application/json',
  bodyMode: overrides.bodyMode ?? 'json',
  bodyContent: overrides.bodyContent ?? '{\n  "ping": "hello"\n}',
  scripts:
    overrides.scripts ??
    "// Inspect the mock response here\npm.test('status should be 200', () => true);",
  settings:
    overrides.settings ?? {
      followRedirects: true,
      strictTls: false,
      persistCookies: true,
    },
  response: overrides.response ?? createEmptyResponse(),
  isSending: overrides.isSending ?? false,
});

const rowsToBulkText = (rows: KeyValueRow[]) =>
  rows
    .filter((row) => row.key.trim() || row.value.trim() || row.description.trim())
    .map((row) => `${row.key}: ${row.value}${row.description ? ` # ${row.description}` : ''}`)
    .join('\n');

const bulkTextToRows = (value: string) => {
  const rows = value
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => {
      const [pairPart, descriptionPart = ''] = line.split('#');
      const separatorIndex = pairPart.includes(':') ? pairPart.indexOf(':') : pairPart.indexOf('=');
      const key = separatorIndex >= 0 ? pairPart.slice(0, separatorIndex).trim() : pairPart.trim();
      const fieldValue = separatorIndex >= 0 ? pairPart.slice(separatorIndex + 1).trim() : '';

      return createKeyValueRow(key, fieldValue, descriptionPart.trim());
    });

  return rows.length > 0 ? rows : [createKeyValueRow()];
};

const getTabSaveLabel = (tab: RequestPageTab) => {
  if (!tab.url.trim()) {
    return tab.title;
  }

  try {
    const parsed = new URL(tab.url);
    const path = parsed.pathname === '/' ? parsed.host : parsed.pathname;
    return `${tab.method} ${path}`;
  } catch {
    return `${tab.method} ${tab.url.trim()}`;
  }
};

const getRowsRecord = (rows: KeyValueRow[]) =>
  rows.reduce<Record<string, string>>((accumulator, row) => {
    if (row.key.trim()) {
      accumulator[row.key.trim()] = row.value.trim();
    }

    return accumulator;
  }, {});

const byteLength = (value: string) =>
  typeof TextEncoder !== 'undefined' ? new TextEncoder().encode(value).length : value.length;

const buildMockResponse = (
  tab: RequestPageTab,
  environment: (typeof ENVIRONMENT_OPTIONS)[number]
): ResponseDraft => {
  if (!tab.url.trim()) {
    return {
      ...createEmptyResponse(),
      error: 'Enter a request URL before sending.',
    };
  }

  let parsedUrl: URL;

  try {
    parsedUrl = new URL(tab.url);
  } catch {
    return {
      ...createEmptyResponse(),
      error: 'The URL is not valid. Try a value like https://localhost:3000/health.',
    };
  }

  const durationMs = 180 + Math.floor(Math.random() * 620);
  const status =
    parsedUrl.pathname.includes('error') ? 500 : tab.method === 'POST' ? 201 : 200;
  const statusLabel = status >= 400 ? 'Server Error' : status === 201 ? 'Created' : 'OK';

  const responseBody = JSON.stringify(
    {
      ok: status < 400,
      environment,
      request: {
        method: tab.method,
        url: tab.url,
        params: getRowsRecord(tab.paramsRows),
        headers: getRowsRecord(tab.headersRows),
        authorization: tab.authorizationMode === 'none' ? null : tab.authorizationMode,
        body_mode: tab.bodyMode,
      },
      data:
        parsedUrl.pathname === '/health'
          ? {
              service: 'mock-gateway',
              status: 'healthy',
              timestamp: new Date().toISOString(),
            }
          : {
              message: 'Mock response generated by the front-end workbench.',
              echoed_body: tab.bodyContent || null,
            },
      meta: {
        request_tab: tab.title,
        sent_at: new Date().toISOString(),
        follow_redirects: tab.settings.followRedirects,
        strict_tls: tab.settings.strictTls,
      },
    },
    null,
    2
  );

  return {
    status,
    statusLabel,
    durationMs,
    sizeBytes: byteLength(responseBody),
    body: responseBody,
    error: null,
  };
};

const getInitialWorkbenchState = (): InitialWorkbenchState => {
  const healthRequest = createRequestPageTab(1, {
    id: 'request-health',
    title: 'Health Check',
    collectionId: 'collection-core',
    method: 'GET',
    url: 'https://localhost:3000/health',
    paramsRows: [createKeyValueRow('verbose', 'true', 'Return dependency detail')],
    paramsBulk: 'verbose: true # Return dependency detail',
  });

  const createUserRequest = createRequestPageTab(2, {
    id: 'request-create-user',
    title: 'Create User',
    collectionId: 'collection-core',
    method: 'POST',
    url: 'https://localhost:3000/api/users',
    activeSection: 'body',
    bodyContent: '{\n  "name": "Ming",\n  "email": "ming@example.com"\n}',
  });

  const tokenRequest = createRequestPageTab(3, {
    id: 'request-issue-token',
    title: 'Issue Token',
    collectionId: 'collection-auth',
    method: 'POST',
    url: 'https://localhost:3000/api/auth/token',
    authorizationMode: 'basic',
    authorizationValue: 'admin:secret',
    bodyContent: '{\n  "scope": "read:all"\n}',
  });

  const tabs = [healthRequest, createUserRequest, tokenRequest];
  const collections: CollectionNode[] = [
    {
      id: 'collection-core',
      name: 'Core APIs',
      color: COLLECTION_COLORS[0],
      requestIds: [healthRequest.id, createUserRequest.id],
    },
    {
      id: 'collection-auth',
      name: 'Identity',
      color: COLLECTION_COLORS[1],
      requestIds: [tokenRequest.id],
    },
    {
      id: 'collection-playground',
      name: 'Playground',
      color: COLLECTION_COLORS[2],
      requestIds: [],
    },
  ];

  return {
    tabs,
    collections,
    activeTabId: healthRequest.id,
    activeCollectionId: 'collection-core',
    expandedCollectionIds: collections.map((collection) => collection.id),
    nextTabIndex: 4,
  };
};

export function ApiRequestWorkbench({
  projectId,
  projectName,
}: {
  projectId: number;
  projectName: string;
}) {
  const initialState = useMemo(() => getInitialWorkbenchState(), []);
  const [tabs, setTabs] = useState<RequestPageTab[]>(initialState.tabs);
  const [collections, setCollections] = useState<CollectionNode[]>(initialState.collections);
  const [activeTabId, setActiveTabId] = useState(initialState.activeTabId);
  const [activeCollectionId, setActiveCollectionId] = useState<string | null>(
    initialState.activeCollectionId
  );
  const [expandedCollectionIds, setExpandedCollectionIds] = useState<string[]>(
    initialState.expandedCollectionIds
  );
  const [nextTabIndex, setNextTabIndex] = useState(initialState.nextTabIndex);
  const [environment, setEnvironment] =
    useState<(typeof ENVIRONMENT_OPTIONS)[number]>('development');
  const [sidebarQuery, setSidebarQuery] = useState('');
  const [deletingCollectionId, setDeletingCollectionId] = useState<string | null>(null);
  const [renamingCollectionId, setRenamingCollectionId] = useState<string | null>(null);
  const [renameDialogCollectionId, setRenameDialogCollectionId] = useState<string | null>(null);
  const [renameDraftName, setRenameDraftName] = useState('');
  const deleteCollectionMutation = useDeleteCollection(projectId);
  const updateCollectionMutation = useUpdateCollection(projectId);

  const deferredSidebarQuery = useDeferredValue(sidebarQuery);

  const tabMap = useMemo(() => new Map(tabs.map((tab) => [tab.id, tab])), [tabs]);
  const activeTab = useMemo(
    () => tabs.find((tab) => tab.id === activeTabId) ?? tabs[0] ?? null,
    [activeTabId, tabs]
  );
  const scratchpadTabs = useMemo(
    () => tabs.filter((tab) => !tab.collectionId),
    [tabs]
  );

  const collectionViews = useMemo(() => {
    const normalizedQuery = deferredSidebarQuery.trim().toLowerCase();

    return collections.reduce<Array<{ collection: CollectionNode; requests: RequestPageTab[] }>>(
      (accumulator, collection) => {
        const requests = collection.requestIds
          .map((requestId) => tabMap.get(requestId))
          .filter((request): request is RequestPageTab => Boolean(request));

        if (!normalizedQuery) {
          accumulator.push({ collection, requests });
          return accumulator;
        }

        const collectionMatches = collection.name.toLowerCase().includes(normalizedQuery);
        const requestMatches = requests.filter((request) =>
          [request.title, request.url, request.method]
            .some((value) => value.toLowerCase().includes(normalizedQuery))
        );

        if (collectionMatches || requestMatches.length > 0) {
          accumulator.push({
            collection,
            requests: collectionMatches ? requests : requestMatches,
          });
        }

        return accumulator;
      },
      []
    );
  }, [collections, deferredSidebarQuery, tabMap]);

  const visibleScratchpadTabs = useMemo(() => {
    const normalizedQuery = deferredSidebarQuery.trim().toLowerCase();

    if (!normalizedQuery) {
      return scratchpadTabs;
    }

    return scratchpadTabs.filter((tab) =>
      [tab.title, tab.url, tab.method].some((value) =>
        value.toLowerCase().includes(normalizedQuery)
      )
    );
  }, [deferredSidebarQuery, scratchpadTabs]);

  const updateTab = (tabId: string, updater: (tab: RequestPageTab) => RequestPageTab) => {
    setTabs((current) => current.map((tab) => (tab.id === tabId ? updater(tab) : tab)));
  };

  const updateActiveTab = (updater: (tab: RequestPageTab) => RequestPageTab) => {
    if (!activeTab) {
      return;
    }

    updateTab(activeTab.id, updater);
  };

  const createStandaloneTab = () => {
    const nextTab = createRequestPageTab(nextTabIndex);

    startTransition(() => {
      setTabs((current) => [...current, nextTab]);
      setActiveTabId(nextTab.id);
      setNextTabIndex((current) => current + 1);
    });
  };

  const createCollection = () => {
    const collectionNumber = collections.length + 1;
    const nextCollection: CollectionNode = {
      id: createLocalId('collection'),
      name: `New Collection ${collectionNumber}`,
      color: COLLECTION_COLORS[(collectionNumber - 1) % COLLECTION_COLORS.length],
      requestIds: [],
    };

    startTransition(() => {
      setCollections((current) => [nextCollection, ...current]);
      setExpandedCollectionIds((current) => [nextCollection.id, ...current]);
      setActiveCollectionId(nextCollection.id);
    });
  };

  const removeCollectionFromWorkbench = (collectionId: string) => {
    const targetCollection = collections.find((collection) => collection.id === collectionId);
    if (!targetCollection) {
      return;
    }

    const remainingCollections = collections.filter((collection) => collection.id !== collectionId);
    let remainingTabs = tabs.filter((tab) => tab.collectionId !== collectionId);
    const activeTabRemoved = remainingTabs.every((tab) => tab.id !== activeTabId);

    let nextActiveTabId = activeTabId;

    if (remainingTabs.length === 0) {
      const nextTab = createRequestPageTab(nextTabIndex);
      remainingTabs = [nextTab];
      nextActiveTabId = nextTab.id;
      setNextTabIndex((current) => current + 1);
    } else if (activeTabRemoved) {
      nextActiveTabId = remainingTabs[0].id;
    }

    startTransition(() => {
      setCollections(remainingCollections);
      setTabs(remainingTabs);
      setExpandedCollectionIds((current) => current.filter((id) => id !== collectionId));
      setActiveCollectionId((current) =>
        current === collectionId ? remainingCollections[0]?.id ?? null : current
      );
      setActiveTabId(nextActiveTabId);
    });
  };

  const handleDeleteCollection = async (collection: CollectionNode) => {
    if (deletingCollectionId) {
      return;
    }

    setDeletingCollectionId(collection.id);

    try {
      const persistedCollectionId = Number(collection.id);

      if (Number.isInteger(persistedCollectionId) && persistedCollectionId > 0) {
        await deleteCollectionMutation.mutateAsync(persistedCollectionId);
      }

      removeCollectionFromWorkbench(collection.id);
    } finally {
      setDeletingCollectionId(null);
    }
  };

  const openRenameCollectionDialog = (collection: CollectionNode) => {
    setRenameDialogCollectionId(collection.id);
    setRenameDraftName(collection.name);
  };

  const closeRenameCollectionDialog = (open: boolean) => {
    if (!open) {
      setRenameDialogCollectionId(null);
      setRenameDraftName('');
    }
  };

  const handleRenameCollection = async () => {
    if (!renameDialogCollectionId || renamingCollectionId) {
      return;
    }

    const nextName = renameDraftName.trim();
    if (!nextName) {
      return;
    }

    const targetCollection = collections.find(
      (collection) => collection.id === renameDialogCollectionId
    );
    if (!targetCollection) {
      closeRenameCollectionDialog(false);
      return;
    }

    setRenamingCollectionId(targetCollection.id);

    try {
      const persistedCollectionId = Number(targetCollection.id);

      if (Number.isInteger(persistedCollectionId) && persistedCollectionId > 0) {
        await updateCollectionMutation.mutateAsync({
          collectionId: persistedCollectionId,
          data: { name: nextName },
        });
      }

      startTransition(() => {
        setCollections((current) =>
          current.map((collection) =>
            collection.id === targetCollection.id ? { ...collection, name: nextName } : collection
          )
        );
      });
      closeRenameCollectionDialog(false);
    } finally {
      setRenamingCollectionId(null);
    }
  };

  const toggleCollection = (collectionId: string) => {
    setActiveCollectionId(collectionId);
    setExpandedCollectionIds((current) =>
      current.includes(collectionId)
        ? current.filter((id) => id !== collectionId)
        : [...current, collectionId]
    );
  };

  const selectRequest = (tabId: string, collectionId: string | null) => {
    setActiveTabId(tabId);

    if (collectionId) {
      setActiveCollectionId(collectionId);
      setExpandedCollectionIds((current) =>
        current.includes(collectionId) ? current : [...current, collectionId]
      );
    }
  };

  const handleDuplicateTab = () => {
    if (!activeTab) {
      return;
    }

    const duplicatedTab: RequestPageTab = {
      ...activeTab,
      id: createLocalId('request-tab'),
      title: `${activeTab.title} Copy`,
      response: createEmptyResponse(),
      isSending: false,
      paramsRows: activeTab.paramsRows.map((row) => ({ ...row, id: createLocalId('kv') })),
      headersRows: activeTab.headersRows.map((row) => ({ ...row, id: createLocalId('kv') })),
    };

    startTransition(() => {
      setTabs((current) => [...current, duplicatedTab]);

      if (duplicatedTab.collectionId) {
        setCollections((current) =>
          current.map((collection) =>
            collection.id === duplicatedTab.collectionId
              ? {
                  ...collection,
                  requestIds: [duplicatedTab.id, ...collection.requestIds],
                }
              : collection
          )
        );
      }

      setActiveTabId(duplicatedTab.id);
      setNextTabIndex((current) => current + 1);
    });
  };

  const handleSaveTab = () => {
    updateActiveTab((tab) => ({
      ...tab,
      title: getTabSaveLabel(tab),
    }));
  };

  const handleSend = () => {
    if (!activeTab) {
      return;
    }

    const tabSnapshot = activeTab;
    const tabId = activeTab.id;

    updateTab(tabId, (tab) => ({
      ...tab,
      isSending: true,
    }));

    window.setTimeout(() => {
      const nextResponse = buildMockResponse(tabSnapshot, environment);

      startTransition(() => {
        updateTab(tabId, (tab) => ({
          ...tab,
          isSending: false,
          response: nextResponse,
        }));
      });
    }, 700);
  };

  if (!activeTab) {
    return null;
  }

  return (
    <main className="flex h-full min-h-0 flex-col overflow-hidden bg-[radial-gradient(circle_at_top_left,_rgba(59,130,246,0.10),_transparent_28%),linear-gradient(180deg,_rgba(255,255,255,0.98),_rgba(244,247,251,0.98))]">
      <div className="border-b border-border/60 bg-white/80 backdrop-blur">
        <div className="flex flex-col gap-3 px-4 py-4 md:px-6">
          <div className="flex items-center gap-3">
            <div className="min-w-0 flex-1 overflow-hidden">
              <RequestTabs
                tabs={tabs}
                activeTabId={activeTabId}
                onSelectTab={(tabId) => selectRequest(tabId, tabMap.get(tabId)?.collectionId ?? null)}
              />
            </div>
            <Button type="button" variant="outline" size="sm" isIcon onClick={createStandaloneTab}>
              <Plus className="h-4 w-4" />
            </Button>
            <EnvironmentSwitcher
              environment={environment}
              onEnvironmentChange={(value) =>
                setEnvironment(value as (typeof ENVIRONMENT_OPTIONS)[number])
              }
            />
          </div>

          <div className="flex flex-wrap items-center gap-2 text-xs text-text-muted">
            <span className="rounded-full border border-border/60 bg-background/80 px-2.5 py-1">
              {projectName}
            </span>
            <span className="rounded-full border border-border/60 bg-background/80 px-2.5 py-1">
              Environment: {environment}
            </span>
          </div>
        </div>
      </div>

      <div className="min-h-0 flex-1 overflow-hidden xl:flex-row xl:flex">
        <aside className="w-full shrink-0 border-b border-border/60 bg-white/82 backdrop-blur xl:w-[320px] xl:border-b-0 xl:border-r">
          <CollectionsSidebar
            collections={collectionViews}
            activeCollectionId={activeCollectionId}
            activeTabId={activeTabId}
            deletingCollectionId={deletingCollectionId}
            renamingCollectionId={renamingCollectionId}
            expandedCollectionIds={expandedCollectionIds}
            scratchpadTabs={visibleScratchpadTabs}
            query={sidebarQuery}
            onQueryChange={setSidebarQuery}
            onCreateCollection={createCollection}
            onDeleteCollection={handleDeleteCollection}
            onRenameCollection={openRenameCollectionDialog}
            onToggleCollection={toggleCollection}
            onSelectRequest={selectRequest}
          />
        </aside>

        <div className="min-h-0 min-w-0 flex-1 overflow-auto p-4 md:p-6">
          <div className="mx-auto flex min-h-full max-w-[1600px] flex-col gap-4">
            <Card className="gap-0 rounded-[28px] border-border/60 bg-white/90 py-0 shadow-[0_12px_44px_rgba(15,23,42,0.08)]">
              <CardHeader className="gap-4 border-b border-border/60 py-5">
                <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
                  <div className="space-y-2">
                    <div className="flex items-center gap-2">
                      <Badge variant="outline" className="border-primary/20 bg-primary/10 text-primary">
                        API Request
                      </Badge>
                      {activeTab.collectionId ? (
                        <Badge variant="secondary">
                          {collections.find((collection) => collection.id === activeTab.collectionId)?.name || 'Collection'}
                        </Badge>
                      ) : (
                        <Badge variant="secondary">Scratchpad</Badge>
                      )}
                    </div>
                    <div>
                      <CardTitle className="text-xl tracking-tight">{activeTab.title}</CardTitle>
                      <CardDescription className="mt-1">
                        Lightweight request workbench with local mock responses and per-tab draft state.
                      </CardDescription>
                    </div>
                  </div>

                  <div className="flex flex-wrap gap-2">
                    <Button type="button" variant="outline" onClick={handleDuplicateTab}>
                      <Copy className="h-4 w-4" />
                      Duplicate
                    </Button>
                    <Button type="button" variant="outline" onClick={handleSaveTab}>
                      <Save className="h-4 w-4" />
                      Save tab
                    </Button>
                  </div>
                </div>
              </CardHeader>

              <CardContent className="space-y-5 px-4 py-5 md:px-6">
                <RequestToolbar
                  tab={activeTab}
                  onMethodChange={(method) => updateActiveTab((tab) => ({ ...tab, method }))}
                  onUrlChange={(url) => updateActiveTab((tab) => ({ ...tab, url }))}
                  onSend={handleSend}
                  onSave={handleSaveTab}
                  onDuplicate={handleDuplicateTab}
                />

                <RequestSectionTabs
                  activeSection={activeTab.activeSection}
                  onSelectSection={(section) =>
                    updateActiveTab((tab) => ({ ...tab, activeSection: section }))
                  }
                />

                <RequestSectionPanel
                  tab={activeTab}
                  onTabChange={updateActiveTab}
                />
              </CardContent>
            </Card>

            <ResponsePanel response={activeTab.response} isSending={activeTab.isSending} />
          </div>
        </div>
      </div>

      <RenameCollectionDialog
        open={renameDialogCollectionId !== null}
        value={renameDraftName}
        isSubmitting={renamingCollectionId !== null}
        onOpenChange={closeRenameCollectionDialog}
        onValueChange={setRenameDraftName}
        onConfirm={handleRenameCollection}
      />
    </main>
  );
}

function CollectionsSidebar({
  collections,
  activeCollectionId,
  activeTabId,
  deletingCollectionId,
  renamingCollectionId,
  expandedCollectionIds,
  scratchpadTabs,
  query,
  onQueryChange,
  onCreateCollection,
  onDeleteCollection,
  onRenameCollection,
  onToggleCollection,
  onSelectRequest,
}: {
  collections: Array<{ collection: CollectionNode; requests: RequestPageTab[] }>;
  activeCollectionId: string | null;
  activeTabId: string;
  deletingCollectionId: string | null;
  renamingCollectionId: string | null;
  expandedCollectionIds: string[];
  scratchpadTabs: RequestPageTab[];
  query: string;
  onQueryChange: (value: string) => void;
  onCreateCollection: () => void;
  onDeleteCollection: (collection: CollectionNode) => void;
  onRenameCollection: (collection: CollectionNode) => void;
  onToggleCollection: (collectionId: string) => void;
  onSelectRequest: (tabId: string, collectionId: string | null) => void;
}) {
  return (
    <div className="flex h-full min-h-0 flex-col overflow-hidden">
      <div className="space-y-4 p-4">
        <div className="flex items-center gap-2">
          <div className="relative min-w-0 flex-1">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted" />
            <Input
              value={query}
              onChange={(event) => onQueryChange(event.target.value)}
              placeholder="Filter collections or requests"
              className="pl-9"
            />
          </div>
          <Button type="button" variant="outline" size="sm" isIcon onClick={onCreateCollection}>
            <Plus className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <div className="min-h-0 flex-1 overflow-y-auto px-3 pb-4">
        <div className="space-y-2">
          {collections.map(({ collection, requests }) => {
            const isExpanded = expandedCollectionIds.includes(collection.id);
            const isActiveCollection = activeCollectionId === collection.id;

            return (
              <div
                key={collection.id}
                className={cn(
                  'group/collection rounded-[18px] border border-transparent bg-white/55 p-1.5 transition-colors',
                  isActiveCollection ? 'bg-primary/5' : 'hover:bg-white/80'
                )}
              >
                <div className="flex items-start gap-2">
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    isIcon
                    className="mt-0.5 h-8 w-8 rounded-xl"
                    onClick={() => onToggleCollection(collection.id)}
                  >
                    {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
                  </Button>

                  <button
                    type="button"
                    onClick={() => onToggleCollection(collection.id)}
                    className="min-w-0 flex-1 rounded-xl px-1 py-0.5 text-left"
                  >
                    <div className="flex items-center gap-2">
                      <span
                        className="h-2.5 w-2.5 rounded-full"
                        style={{ backgroundColor: collection.color }}
                        aria-hidden="true"
                      />
                      <p className="truncate text-sm font-medium text-text-main">{collection.name}</p>
                    </div>
                    <p className="mt-0.5 text-[11px] text-text-muted">
                      {collection.requestIds.length} requests
                    </p>
                  </button>

                  <CollectionActionsMenu
                    isDeleting={deletingCollectionId === collection.id}
                    isRenaming={renamingCollectionId === collection.id}
                    onRename={() => onRenameCollection(collection)}
                    onDelete={() => void onDeleteCollection(collection)}
                  />
                </div>

                {isExpanded ? (
                  <div className="mt-1.5 space-y-1 pl-10">
                    {requests.map((request) => (
                      <button
                        key={request.id}
                        type="button"
                        onClick={() => onSelectRequest(request.id, collection.id)}
                        className={cn(
                          'w-full rounded-xl px-3 py-1.5 text-left transition-colors',
                          activeTabId === request.id
                            ? 'bg-primary/10 text-text-main'
                            : 'hover:bg-white/80'
                        )}
                      >
                        <div className="flex items-center gap-2">
                          <MethodBadge method={request.method} compact />
                          <p className="truncate text-sm font-medium">{request.title}</p>
                        </div>
                      </button>
                    ))}
                  </div>
                ) : null}
              </div>
            );
          })}

          {scratchpadTabs.length > 0 ? (
            <div className="pt-4">
              <div className="mb-2 px-2 text-xs font-medium uppercase tracking-[0.16em] text-text-muted">
                Scratchpad
              </div>
              <div className="space-y-1.5">
                {scratchpadTabs.map((tab) => (
                  <button
                    key={tab.id}
                    type="button"
                    onClick={() => onSelectRequest(tab.id, null)}
                    className={cn(
                      'w-full rounded-xl px-3 py-1.5 text-left transition-colors',
                      activeTabId === tab.id ? 'bg-primary/10 text-text-main' : 'hover:bg-white/80'
                    )}
                  >
                    <div className="flex items-center gap-2">
                      <FolderOpen className="h-4 w-4 text-text-muted" />
                      <p className="truncate text-sm font-medium">{tab.title}</p>
                    </div>
                  </button>
                ))}
              </div>
            </div>
          ) : null}
        </div>
      </div>
    </div>
  );
}

function CollectionActionsMenu({
  isDeleting,
  isRenaming,
  onRename,
  onDelete,
}: {
  isDeleting: boolean;
  isRenaming: boolean;
  onRename: () => void;
  onDelete: () => void;
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          type="button"
          variant="ghost"
          size="sm"
          isIcon
          className="h-8 w-8 rounded-xl opacity-0 transition-opacity group-hover/collection:opacity-100 focus-visible:opacity-100 data-[state=open]:opacity-100"
          aria-label="Open collection actions"
        >
          <MoreHorizontal className="h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-44 rounded-xl">
        <DropdownMenuItem>New request</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem>Import</DropdownMenuItem>
        <DropdownMenuItem>Export</DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem disabled={isRenaming} onSelect={onRename}>
          Rename
        </DropdownMenuItem>
        <DropdownMenuItem variant="destructive" disabled={isDeleting} onSelect={onDelete}>
          Delete
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function RenameCollectionDialog({
  open,
  value,
  isSubmitting,
  onOpenChange,
  onValueChange,
  onConfirm,
}: {
  open: boolean;
  value: string;
  isSubmitting: boolean;
  onOpenChange: (open: boolean) => void;
  onValueChange: (value: string) => void;
  onConfirm: () => Promise<void>;
}) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent size="sm">
        <DialogHeader>
          <DialogTitle>Rename Collection</DialogTitle>
          <DialogDescription>
            Update the collection name and sync it to the backend when a persisted collection ID is
            available.
          </DialogDescription>
        </DialogHeader>

        <DialogBody>
          <div className="space-y-2">
            <Label htmlFor="rename-collection-name">Collection name</Label>
            <Input
              id="rename-collection-name"
              value={value}
              onChange={(event) => onValueChange(event.target.value)}
              placeholder="Enter collection name"
              className="rounded-2xl"
            />
          </div>
        </DialogBody>

        <DialogFooter>
          <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            type="button"
            loading={isSubmitting}
            disabled={!value.trim()}
            onClick={() => void onConfirm()}
          >
            Save
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function RequestTabs({
  tabs,
  activeTabId,
  onSelectTab,
}: {
  tabs: RequestPageTab[];
  activeTabId: string;
  onSelectTab: (tabId: string) => void;
}) {
  return (
    <div className="overflow-x-auto">
      <div className="flex min-w-max items-center gap-2 pr-2">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            type="button"
            onClick={() => onSelectTab(tab.id)}
            className={cn(
              'group inline-flex items-center gap-2 rounded-2xl border px-3 py-2 text-sm transition-all',
              tab.id === activeTabId
                ? 'border-primary/30 bg-primary/10 text-text-main shadow-sm'
                : 'border-border/60 bg-white/75 text-text-muted hover:border-border hover:bg-white hover:text-text-main'
            )}
          >
            <span
              className={cn(
                'h-2 w-2 rounded-full',
                tab.id === activeTabId ? 'bg-primary' : 'bg-text-muted/40'
              )}
            />
            <span className="truncate font-medium">{tab.title}</span>
          </button>
        ))}
      </div>
    </div>
  );
}

function EnvironmentSwitcher({
  environment,
  onEnvironmentChange,
}: {
  environment: (typeof ENVIRONMENT_OPTIONS)[number];
  onEnvironmentChange: (value: string) => void;
}) {
  return (
    <div className="flex items-center gap-2 rounded-2xl border border-border/60 bg-white/85 px-3 py-1.5 shadow-sm">
      <span className="text-xs font-medium uppercase tracking-[0.16em] text-text-muted">
        Environment
      </span>
      <Select value={environment} onValueChange={onEnvironmentChange}>
        <SelectTrigger className="h-8 min-w-[156px] border-0 bg-transparent px-2 shadow-none">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {ENVIRONMENT_OPTIONS.map((option) => (
            <SelectItem key={option} value={option}>
              {option}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}

function RequestToolbar({
  tab,
  onMethodChange,
  onUrlChange,
  onSend,
  onSave,
  onDuplicate,
}: {
  tab: RequestPageTab;
  onMethodChange: (method: RequestMethod) => void;
  onUrlChange: (url: string) => void;
  onSend: () => void;
  onSave: () => void;
  onDuplicate: () => void;
}) {
  return (
    <div className="rounded-[24px] border border-border/60 bg-slate-50/90 p-3 shadow-[inset_0_1px_0_rgba(255,255,255,0.9)]">
      <div className="grid gap-3 xl:grid-cols-[140px_minmax(0,1fr)_auto]">
        <Select value={tab.method} onValueChange={(value) => onMethodChange(value as RequestMethod)}>
          <SelectTrigger className="h-11 w-full rounded-2xl border-border/70 bg-white font-semibold">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {METHOD_OPTIONS.map((method) => (
              <SelectItem key={method} value={method}>
                {method}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>

        <Input
          value={tab.url}
          onChange={(event) => onUrlChange(event.target.value)}
          placeholder="https://localhost:3000/health"
          className="h-11 rounded-2xl border-border/70 bg-white px-4 text-sm shadow-none"
        />

        <div className="flex flex-wrap gap-2">
          <Button type="button" variant="outline" className="h-11 rounded-2xl" onClick={onSave}>
            <Save className="h-4 w-4" />
            Save
          </Button>
          <Button type="button" variant="outline" className="h-11 rounded-2xl" onClick={onDuplicate}>
            <Copy className="h-4 w-4" />
          </Button>
          <Button type="button" className="h-11 rounded-2xl px-5" onClick={onSend} loading={tab.isSending}>
            <SendHorizonal className="h-4 w-4" />
            Send
          </Button>
        </div>
      </div>
    </div>
  );
}

function RequestSectionTabs({
  activeSection,
  onSelectSection,
}: {
  activeSection: RequestSection;
  onSelectSection: (section: RequestSection) => void;
}) {
  return (
    <div className="flex flex-wrap items-center gap-2">
      {SECTION_ITEMS.map((item) => (
        <button
          key={item.value}
          type="button"
          onClick={() => onSelectSection(item.value)}
          className={cn(
            'rounded-full border px-3 py-2 text-sm font-medium transition-colors',
            item.value === activeSection
              ? 'border-primary/30 bg-primary/10 text-primary shadow-sm'
              : 'border-border/60 bg-white/70 text-text-muted hover:border-border hover:bg-white hover:text-text-main'
          )}
        >
          {item.label}
        </button>
      ))}
    </div>
  );
}

function RequestSectionPanel({
  tab,
  onTabChange,
}: {
  tab: RequestPageTab;
  onTabChange: (updater: (tab: RequestPageTab) => RequestPageTab) => void;
}) {
  switch (tab.activeSection) {
    case 'params':
      return (
        <KeyValueEditor
          title="Query Params"
          description="Edit structured query parameters or switch to bulk mode for quick pasting."
          mode={tab.paramsMode}
          rows={tab.paramsRows}
          bulkValue={tab.paramsBulk}
          onModeChange={(mode) =>
            onTabChange((current) =>
              mode === 'bulk'
                ? {
                    ...current,
                    paramsMode: mode,
                    paramsBulk: rowsToBulkText(current.paramsRows),
                  }
                : {
                    ...current,
                    paramsMode: mode,
                    paramsRows: bulkTextToRows(current.paramsBulk),
                  }
            )
          }
          onRowsChange={(rows) =>
            onTabChange((current) => ({
              ...current,
              paramsRows: rows,
              paramsBulk: rowsToBulkText(rows),
            }))
          }
          onBulkChange={(bulkValue) =>
            onTabChange((current) => ({
              ...current,
              paramsBulk: bulkValue,
            }))
          }
        />
      );
    case 'authorization':
      return (
        <AuthorizationPanel
          mode={tab.authorizationMode}
          value={tab.authorizationValue}
          onModeChange={(mode) =>
            onTabChange((current) => ({
              ...current,
              authorizationMode: mode,
            }))
          }
          onValueChange={(value) =>
            onTabChange((current) => ({
              ...current,
              authorizationValue: value,
            }))
          }
        />
      );
    case 'headers':
      return (
        <KeyValueEditor
          title="Headers"
          description="Manage request headers with a table view or bulk input."
          mode={tab.headersMode}
          rows={tab.headersRows}
          bulkValue={tab.headersBulk}
          onModeChange={(mode) =>
            onTabChange((current) =>
              mode === 'bulk'
                ? {
                    ...current,
                    headersMode: mode,
                    headersBulk: rowsToBulkText(current.headersRows),
                  }
                : {
                    ...current,
                    headersMode: mode,
                    headersRows: bulkTextToRows(current.headersBulk),
                  }
            )
          }
          onRowsChange={(rows) =>
            onTabChange((current) => ({
              ...current,
              headersRows: rows,
              headersBulk: rowsToBulkText(rows),
            }))
          }
          onBulkChange={(bulkValue) =>
            onTabChange((current) => ({
              ...current,
              headersBulk: bulkValue,
            }))
          }
        />
      );
    case 'body':
      return (
        <BodyEditor
          mode={tab.bodyMode}
          value={tab.bodyContent}
          onModeChange={(mode) =>
            onTabChange((current) => ({
              ...current,
              bodyMode: mode,
            }))
          }
          onValueChange={(value) =>
            onTabChange((current) => ({
              ...current,
              bodyContent: value,
            }))
          }
        />
      );
    case 'scripts':
      return (
        <ScriptsPanel
          value={tab.scripts}
          onValueChange={(value) =>
            onTabChange((current) => ({
              ...current,
              scripts: value,
            }))
          }
        />
      );
    case 'settings':
      return (
        <SettingsPanel
          settings={tab.settings}
          onSettingChange={(key, value) =>
            onTabChange((current) => ({
              ...current,
              settings: {
                ...current.settings,
                [key]: value,
              },
            }))
          }
        />
      );
    default:
      return null;
  }
}

function KeyValueEditor({
  title,
  description,
  mode,
  rows,
  bulkValue,
  onModeChange,
  onRowsChange,
  onBulkChange,
}: {
  title: string;
  description: string;
  mode: BulkMode;
  rows: KeyValueRow[];
  bulkValue: string;
  onModeChange: (mode: BulkMode) => void;
  onRowsChange: (rows: KeyValueRow[]) => void;
  onBulkChange: (value: string) => void;
}) {
  const updateRow = (rowId: string, patch: Partial<KeyValueRow>) => {
    onRowsChange(rows.map((row) => (row.id === rowId ? { ...row, ...patch } : row)));
  };

  const removeRow = (rowId: string) => {
    const nextRows = rows.filter((row) => row.id !== rowId);
    onRowsChange(nextRows.length > 0 ? nextRows : [createKeyValueRow()]);
  };

  return (
    <div className="rounded-[24px] border border-border/60 bg-white/85 shadow-sm">
      <div className="flex flex-col gap-4 border-b border-border/60 px-5 py-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h3 className="text-base font-semibold text-text-main">{title}</h3>
          <p className="mt-1 text-sm text-text-muted">{description}</p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <div className="inline-flex rounded-full border border-border/60 bg-slate-50/80 p-1">
            <button
              type="button"
              onClick={() => onModeChange('table')}
              className={cn(
                'rounded-full px-3 py-1.5 text-xs font-medium transition-colors',
                mode === 'table'
                  ? 'bg-white text-text-main shadow-sm'
                  : 'text-text-muted hover:text-text-main'
              )}
            >
              Table
            </button>
            <button
              type="button"
              onClick={() => onModeChange('bulk')}
              className={cn(
                'rounded-full px-3 py-1.5 text-xs font-medium transition-colors',
                mode === 'bulk'
                  ? 'bg-white text-text-main shadow-sm'
                  : 'text-text-muted hover:text-text-main'
              )}
            >
              Bulk Edit
            </button>
          </div>

          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => onRowsChange([...rows, createKeyValueRow()])}
          >
            <Plus className="h-4 w-4" />
            Add row
          </Button>
        </div>
      </div>

      {mode === 'bulk' ? (
        <div className="px-5 py-5">
          <Textarea
            value={bulkValue}
            onChange={(event) => onBulkChange(event.target.value)}
            rows={10}
            className="min-h-[220px] rounded-2xl font-mono text-sm"
            placeholder="key: value # description"
          />
        </div>
      ) : (
        <div className="overflow-x-auto px-5 py-5">
          <div className="min-w-[760px] space-y-3">
            <div className="grid grid-cols-[1.05fr_1.25fr_1fr_56px] gap-3 px-3 text-xs font-medium uppercase tracking-[0.16em] text-text-muted">
              <span>Key</span>
              <span>Value</span>
              <span>Description</span>
              <span />
            </div>

            {rows.map((row) => (
              <div key={row.id} className="grid grid-cols-[1.05fr_1.25fr_1fr_56px] gap-3">
                <Input
                  value={row.key}
                  onChange={(event) => updateRow(row.id, { key: event.target.value })}
                  placeholder="page"
                  className="rounded-2xl"
                />
                <Input
                  value={row.value}
                  onChange={(event) => updateRow(row.id, { value: event.target.value })}
                  placeholder="1"
                  className="rounded-2xl"
                />
                <Input
                  value={row.description}
                  onChange={(event) => updateRow(row.id, { description: event.target.value })}
                  placeholder="Optional note"
                  className="rounded-2xl"
                />
                <Button
                  type="button"
                  variant="ghost"
                  isIcon
                  className="h-9 w-9 rounded-2xl"
                  onClick={() => removeRow(row.id)}
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

function AuthorizationPanel({
  mode,
  value,
  onModeChange,
  onValueChange,
}: {
  mode: AuthorizationMode;
  value: string;
  onModeChange: (mode: AuthorizationMode) => void;
  onValueChange: (value: string) => void;
}) {
  return (
    <div className="grid gap-4 lg:grid-cols-[220px_1fr]">
      <Card className="border-border/60 bg-white/85 py-0 shadow-sm">
        <CardHeader className="border-b border-border/60 py-5">
          <CardTitle>Authorization</CardTitle>
          <CardDescription>Choose how this request should authenticate.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4 px-5 py-5">
          <div className="space-y-2">
            <Label htmlFor="request-auth-mode">Auth type</Label>
            <Select value={mode} onValueChange={(nextValue) => onModeChange(nextValue as AuthorizationMode)}>
              <SelectTrigger id="request-auth-mode" className="rounded-2xl">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {AUTHORIZATION_OPTIONS.map((option) => (
                  <SelectItem key={option} value={option}>
                    {option}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      <Card className="border-border/60 bg-white/85 py-0 shadow-sm">
        <CardHeader className="border-b border-border/60 py-5">
          <CardTitle>Credentials</CardTitle>
          <CardDescription>Provide the mock secret or token value for the selected auth scheme.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4 px-5 py-5">
          {mode === 'none' ? (
            <div className="rounded-2xl border border-dashed border-border/70 bg-slate-50/80 p-5 text-sm text-text-muted">
              This request currently sends without authentication.
            </div>
          ) : (
            <div className="space-y-2">
              <Label htmlFor="request-auth-value">
                {mode === 'basic' ? 'Username:Password' : mode === 'api-key' ? 'API key' : 'Token'}
              </Label>
              <Input
                id="request-auth-value"
                value={value}
                onChange={(event) => onValueChange(event.target.value)}
                placeholder={mode === 'basic' ? 'user:secret' : 'Paste credential value'}
                className="rounded-2xl"
              />
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

function BodyEditor({
  mode,
  value,
  onModeChange,
  onValueChange,
}: {
  mode: BodyMode;
  value: string;
  onModeChange: (mode: BodyMode) => void;
  onValueChange: (value: string) => void;
}) {
  return (
    <div className="rounded-[24px] border border-border/60 bg-white/85 shadow-sm">
      <div className="flex flex-col gap-4 border-b border-border/60 px-5 py-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h3 className="text-base font-semibold text-text-main">Body</h3>
          <p className="mt-1 text-sm text-text-muted">
            Choose a body mode and edit the payload in a large code-friendly area.
          </p>
        </div>

        <div className="inline-flex rounded-full border border-border/60 bg-slate-50/80 p-1">
          {BODY_MODE_OPTIONS.map((option) => (
            <button
              key={option}
              type="button"
              onClick={() => onModeChange(option)}
              className={cn(
                'rounded-full px-3 py-1.5 text-xs font-medium capitalize transition-colors',
                option === mode
                  ? 'bg-white text-text-main shadow-sm'
                  : 'text-text-muted hover:text-text-main'
              )}
            >
              {option}
            </button>
          ))}
        </div>
      </div>

      <div className="px-5 py-5">
        <Textarea
          value={value}
          onChange={(event) => onValueChange(event.target.value)}
          rows={14}
          className="min-h-[280px] rounded-2xl font-mono text-sm"
          placeholder={mode === 'form-data' ? 'field=value' : '{\n  \n}'}
        />
      </div>
    </div>
  );
}

function ScriptsPanel({
  value,
  onValueChange,
}: {
  value: string;
  onValueChange: (value: string) => void;
}) {
  return (
    <Card className="border-border/60 bg-white/85 py-0 shadow-sm">
      <CardHeader className="border-b border-border/60 py-5">
        <CardTitle>Scripts</CardTitle>
        <CardDescription>Use this area for pre-request or post-response scripting logic.</CardDescription>
      </CardHeader>
      <CardContent className="px-5 py-5">
        <Textarea
          value={value}
          onChange={(event) => onValueChange(event.target.value)}
          rows={14}
          className="min-h-[280px] rounded-2xl font-mono text-sm"
          placeholder="// Write mock scripts here"
        />
      </CardContent>
    </Card>
  );
}

function SettingsPanel({
  settings,
  onSettingChange,
}: {
  settings: RequestPageTab['settings'];
  onSettingChange: (key: keyof RequestPageTab['settings'], value: boolean) => void;
}) {
  const settingItems: Array<{
    key: keyof RequestPageTab['settings'];
    title: string;
    description: string;
  }> = [
    {
      key: 'followRedirects',
      title: 'Follow redirects',
      description: 'Keeps the mock request aligned with browser-like navigation behavior.',
    },
    {
      key: 'strictTls',
      title: 'Strict TLS validation',
      description: 'Simulate certificate enforcement before wiring a real network client.',
    },
    {
      key: 'persistCookies',
      title: 'Persist cookies',
      description: 'Store cookies between sends for session-driven flows.',
    },
  ];

  return (
    <div className="grid gap-4 lg:grid-cols-3">
      {settingItems.map((item) => (
        <Card key={item.key} className="border-border/60 bg-white/85 py-0 shadow-sm">
          <CardHeader className="border-b border-border/60 py-5">
            <div className="flex items-center justify-between gap-3">
              <div>
                <CardTitle>{item.title}</CardTitle>
                <CardDescription className="mt-1">{item.description}</CardDescription>
              </div>
              <Switch
                checked={settings[item.key]}
                onCheckedChange={(checked) => onSettingChange(item.key, checked)}
              />
            </div>
          </CardHeader>
        </Card>
      ))}
    </div>
  );
}

function ResponsePanel({
  response,
  isSending,
}: {
  response: ResponseDraft;
  isSending: boolean;
}) {
  return (
    <Card className="min-h-[320px] gap-0 rounded-[28px] border-border/60 bg-white/90 py-0 shadow-[0_12px_44px_rgba(15,23,42,0.06)]">
      <CardHeader className="gap-4 border-b border-border/60 py-5">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <CardTitle className="text-xl tracking-tight">Response</CardTitle>
            <CardDescription className="mt-1">
              Inspect the latest mock response payload, timing, and status details.
            </CardDescription>
          </div>

          <div className="flex flex-wrap gap-2">
            <MetricBadge label="Status" value={response.status ? `${response.status} ${response.statusLabel}` : '-'} />
            <MetricBadge label="Time" value={response.durationMs ? `${response.durationMs} ms` : '-'} />
            <MetricBadge label="Size" value={response.sizeBytes ? `${response.sizeBytes} B` : '-'} />
          </div>
        </div>
      </CardHeader>

      <CardContent className="flex min-h-[260px] flex-1 flex-col px-5 py-5">
        {isSending ? (
          <div className="flex flex-1 flex-col items-center justify-center rounded-[24px] border border-dashed border-border/70 bg-slate-50/80 text-center">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-primary/20 border-t-primary" />
            <p className="mt-4 text-sm font-medium text-text-main">Sending mock request...</p>
            <p className="mt-1 text-sm text-text-muted">
              The response panel updates when the simulated request completes.
            </p>
          </div>
        ) : response.error ? (
          <div className="flex flex-1 flex-col justify-center rounded-[24px] border border-rose-200 bg-rose-50/70 p-6">
            <p className="text-sm font-semibold text-rose-700">Unable to send request</p>
            <p className="mt-2 text-sm leading-6 text-rose-600">{response.error}</p>
          </div>
        ) : response.status === null ? (
          <div className="flex flex-1 flex-col items-center justify-center rounded-[24px] border border-dashed border-border/70 bg-slate-50/80 text-center">
            <p className="text-base font-semibold text-text-main">Click Send to get a response</p>
            <p className="mt-2 max-w-xl text-sm leading-6 text-text-muted">
              Once you trigger the mock request, this panel will render formatted JSON, status metadata, and any validation hints.
            </p>
          </div>
        ) : (
          <pre className="flex-1 overflow-auto rounded-[24px] border border-border/60 bg-slate-950/95 p-5 text-sm leading-6 text-slate-100">
            {response.body}
          </pre>
        )}
      </CardContent>
    </Card>
  );
}

function MetricBadge({
  label,
  value,
}: {
  label: string;
  value: string;
}) {
  return (
    <div className="rounded-full border border-border/60 bg-slate-50/80 px-3 py-1.5 text-sm">
      <span className="text-text-muted">{label}: </span>
      <span className="font-medium text-text-main">{value}</span>
    </div>
  );
}

function MethodBadge({
  method,
  compact = false,
}: {
  method: RequestMethod;
  compact?: boolean;
}) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full border font-semibold tracking-[0.14em]',
        compact ? 'px-2 py-0.5 text-[10px]' : 'px-2.5 py-1 text-[11px]',
        METHOD_BADGE_STYLES[method]
      )}
    >
      {method}
    </span>
  );
}
