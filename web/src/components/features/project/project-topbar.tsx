'use client';

import Link from 'next/link';
import { useMemo } from 'react';
import { usePathname, useRouter } from 'next/navigation';
import {
  Bell,
  FolderKanban,
  LayoutPanelTop,
  LogOut,
  Settings,
  Users,
} from 'lucide-react';
import { LanguageSwitcher } from '@/components/common';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { ROUTES } from '@/constants/routes';
import { useLogout } from '@/hooks/use-auth';
import { useAuthStore } from '@/store/auth-store';
import { cn } from '@/utils';

interface NavItem {
  href: string;
  label: string;
}

const NAV_ITEMS: NavItem[] = [
  {
    href: ROUTES.CONSOLE.PROJECTS,
    label: 'Projects',
  },
  {
    href: ROUTES.CONSOLE.HOME,
    label: 'Console',
  },
  {
    href: ROUTES.CONSOLE.USERS,
    label: 'Users',
  },
  {
    href: ROUTES.CONSOLE.SETTINGS,
    label: 'Settings',
  },
];

const buildInitials = (name: string) =>
  name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() || '')
    .join('') || 'U';

export function ProjectTopbar() {
  const pathname = usePathname();
  const router = useRouter();
  const logout = useLogout();
  const user = useAuthStore.use.user();

  const displayName = user?.nickname || user?.username || user?.email || 'User';
  const initials = useMemo(() => buildInitials(displayName), [displayName]);

  const isRouteActive = (href: string) => {
    if (href === ROUTES.CONSOLE.HOME) {
      return pathname === href;
    }

    return pathname === href || pathname.startsWith(`${href}/`);
  };

  const handleLogout = () => {
    logout();
    router.replace(ROUTES.AUTH.LOGIN);
  };

  return (
    <header className="z-40 flex h-16 shrink-0 items-center justify-between border-b border-border/60 bg-bg-surface/95 px-4 backdrop-blur md:px-6">
      <div className="flex min-w-0 items-center gap-4">
        <Link href={ROUTES.CONSOLE.PROJECTS} className="group flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-primary text-primary-foreground shadow-button-primary transition-transform group-hover:scale-105">
            <LayoutPanelTop className="h-4.5 w-4.5" />
          </div>
          <div className="min-w-0">
            <p className="truncate text-sm font-semibold tracking-tight text-text-main">KEST</p>
            <p className="truncate text-xs text-text-muted">Project workspace</p>
          </div>
        </Link>

        <nav className="hidden items-center gap-1 rounded-full border border-border/60 bg-muted/40 p-1 lg:flex">
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                'rounded-full px-3 py-1.5 text-sm transition-colors',
                isRouteActive(item.href)
                  ? 'bg-background text-text-main shadow-sm'
                  : 'text-text-muted hover:text-text-main'
              )}
            >
              {item.label}
            </Link>
          ))}
        </nav>
      </div>

      <div className="flex items-center gap-2">
        <LanguageSwitcher />

        <Button variant="ghost" isIcon className="relative h-9 w-9 rounded-full">
          <Bell className="h-4 w-4 text-text-muted" />
          <span className="absolute right-2 top-2 h-2 w-2 rounded-full border-2 border-bg-surface bg-primary" />
          <span className="sr-only">Notifications</span>
        </Button>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              isIcon
              noScale
              className="h-9 w-9 overflow-hidden rounded-full border border-border/60"
            >
              <Avatar className="h-full w-full">
                <AvatarImage src={user?.avatar} />
                <AvatarFallback className="bg-primary/10 text-primary">
                  {initials}
                </AvatarFallback>
              </Avatar>
              <span className="sr-only">Profile</span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56 rounded-xl p-1 shadow-premium">
            <div className="px-2 py-1.5 text-xs font-medium uppercase tracking-wider text-text-muted">
              {displayName}
            </div>
            <DropdownMenuItem asChild className="cursor-pointer rounded-lg">
              <Link href={ROUTES.CONSOLE.PROJECTS}>
                <FolderKanban className="mr-2 h-4 w-4" />
                Projects
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild className="cursor-pointer rounded-lg">
              <Link href={ROUTES.CONSOLE.USERS}>
                <Users className="mr-2 h-4 w-4" />
                Users
              </Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild className="cursor-pointer rounded-lg">
              <Link href={ROUTES.CONSOLE.SETTINGS}>
                <Settings className="mr-2 h-4 w-4" />
                Settings
              </Link>
            </DropdownMenuItem>
            <div className="my-1 h-px bg-border/60" />
            <DropdownMenuItem
              className="cursor-pointer rounded-lg text-destructive focus:bg-destructive/10"
              onClick={handleLogout}
            >
              <LogOut className="mr-2 h-4 w-4" />
              Logout
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
