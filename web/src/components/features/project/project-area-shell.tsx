'use client';

import { ProjectTopbar } from '@/components/features/project/project-topbar';

export function ProjectAreaShell({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex h-screen flex-col overflow-hidden bg-bg-canvas text-text-main">
      <ProjectTopbar />
      <div className="min-h-0 flex-1 overflow-x-hidden overflow-y-auto">{children}</div>
    </div>
  );
}
