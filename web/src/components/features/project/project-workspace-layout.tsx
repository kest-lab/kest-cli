'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import {
  PROJECT_WORKSPACE_MODULES,
  buildProjectWorkspaceRoute,
} from '@/components/features/project/project-navigation';
import { cn } from '@/utils';

export function ProjectWorkspaceLayout({
  projectId,
  children,
}: {
  projectId: number;
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  return (
    <div className="flex h-full min-h-0 flex-col overflow-hidden lg:flex-row">
      <aside className="w-full shrink-0 border-b border-border/60 bg-bg-surface/70 lg:w-[92px] lg:border-b-0 lg:border-r">
        <div className="flex h-full flex-col overflow-hidden">
          <div className="p-3">
            <Button
              asChild
              variant="ghost"
              size="sm"
              isIcon
              className="h-8 w-8 rounded-full text-text-muted"
            >
              <Link href="/project">
                <ArrowLeft className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>

          <Separator />

          <div className="min-h-0 flex-1 overflow-y-auto p-3">
            <nav className="space-y-2">
              {PROJECT_WORKSPACE_MODULES.map((item) => {
                const Icon = item.icon;
                const href = buildProjectWorkspaceRoute(projectId, item.value);
                const isActive = pathname === href || pathname.startsWith(`${href}?`) || pathname.startsWith(`${href}/`);

                return (
                  <Link
                    key={item.value}
                    href={href}
                    className={cn(
                      'group flex flex-col items-center justify-center gap-1.5 rounded-xl border border-transparent px-2 py-2.5 text-center transition-colors',
                      isActive
                        ? ' text-text-main '
                        : 'text-text-muted hover:bg-background/70 hover:text-black'
                    )}
                  >
                    <div
                      className={cn(
                        'flex h-6 w-6 shrink-0 items-center justify-center rounded-md',
                        isActive ? 'bg-primary/10 text-primary' : 'bg-muted text-text-muted'
                      )}
                    >
                      <Icon className="h-3 w-3" />
                    </div>
                    <div className="min-w-0">
                      <p className="text-[11px] font-medium leading-4">{item.label}</p>
                    </div>
                  </Link>
                );
              })}
            </nav>
          </div>
        </div>
      </aside>

      <div className="min-h-0 min-w-0 flex-1 overflow-hidden">{children}</div>
    </div>
  );
}
