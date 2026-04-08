'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { Menu, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { LanguageSwitcher } from '@/components/common/locale-switcher';
import { cn } from '@/utils';
import type { MarketingNavItem } from './types';

/**
 * @component MarketingNavbar
 * @category Feature
 * @status Stable
 * @description Provides the sticky marketing navigation with mobile menu support.
 * @usage Use in the public site layout for the SaaS homepage.
 * @example
 * <MarketingNavbar brandName="KestFlow" navItems={items} />
 */
export interface MarketingNavbarProps {
  brandName: string;
  navItems: MarketingNavItem[];
  loginLabel: string;
  signUpLabel: string;
  docsSoonLabel: string;
  openMenuLabel: string;
  closeMenuLabel: string;
}

// 判断是否为外部链接。
// 作用：外部文档站链接需要新标签页打开，避免用户直接离开主站。
function isExternalHref(href: string) {
  return href.startsWith('http://') || href.startsWith('https://');
}

function NavItem({
  item,
  docsSoonLabel,
  mobile = false,
  onNavigate,
}: {
  item: MarketingNavItem;
  docsSoonLabel: string;
  mobile?: boolean;
  onNavigate?: () => void;
}) {
  if (item.href) {
    // 顶部导航既支持站内锚点，也支持外部 docs 链接。
    const external = isExternalHref(item.href);

    return (
      <Link
        href={item.href}
        onClick={onNavigate}
        target={external ? '_blank' : undefined}
        rel={external ? 'noreferrer' : undefined}
        className={cn(
          'text-sm font-medium text-text-subtle transition-colors duration-200 hover:text-text-main',
          mobile ? 'rounded-2xl border border-border/70 bg-white px-4 py-3' : ''
        )}
      >
        {item.label}
      </Link>
    );
  }

  return (
    <span
      className={cn(
        'inline-flex items-center gap-2 text-sm font-medium text-text-muted',
        mobile ? 'rounded-2xl border border-dashed border-border/80 px-4 py-3' : ''
      )}
    >
      {item.label}
      {item.placeholder ? (
        <span className="rounded-full border border-border/80 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-[0.22em]">
          {docsSoonLabel}
        </span>
      ) : null}
    </span>
  );
}

export function MarketingNavbar({
  brandName,
  navItems,
  loginLabel,
  signUpLabel,
  docsSoonLabel,
  openMenuLabel,
  closeMenuLabel,
}: MarketingNavbarProps) {
  const [open, setOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 12);

    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });

    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  return (
    <header
      className={cn(
        'sticky top-0 z-50 border-b border-transparent transition-all duration-300',
        scrolled
          ? 'bg-white/82 shadow-[0_18px_45px_-34px_rgba(15,23,42,0.28)] backdrop-blur-xl'
          : 'bg-white/72 backdrop-blur-md'
      )}
    >
      <div className="container">
        <div className="flex h-[4.5rem] items-center justify-between gap-6">
          <Link href="/" className="flex items-center gap-3">
            <p className="text-2xl font-semibold tracking-tight text-text-main [font-family:var(--font-space-grotesk)]">
              {brandName}
            </p>
          </Link>

          <nav className="hidden items-center gap-6 lg:flex">
            {navItems.map((item) => (
              <NavItem key={item.label} item={item} docsSoonLabel={docsSoonLabel} />
            ))}
          </nav>

          <div className="hidden items-center gap-2 lg:flex">
            <LanguageSwitcher />
            <Button asChild variant="ghost" size="lg" className="rounded-2xl">
              <Link href="/login">{loginLabel}</Link>
            </Button>
            <Button asChild size="lg" className="rounded-2xl bg-slate-950 text-white hover:bg-slate-900">
              <Link href="/register">{signUpLabel}</Link>
            </Button>
          </div>

          <div className="flex items-center gap-2 lg:hidden">
            <LanguageSwitcher />
            <Button
              type="button"
              variant="outline"
              size="lg"
              isIcon
              noScale
              className="size-11 rounded-2xl bg-white"
              aria-label={open ? closeMenuLabel : openMenuLabel}
              onClick={() => setOpen((value) => !value)}
            >
              {open ? <X className="size-5" /> : <Menu className="size-5" />}
            </Button>
          </div>
        </div>

        <div
          className={cn(
            'overflow-hidden transition-all duration-300 lg:hidden',
            open ? 'max-h-[28rem] pb-5' : 'max-h-0'
          )}
        >
          <div className="rounded-[1.75rem] border border-border/70 bg-white/95 p-4 shadow-[0_20px_55px_-40px_rgba(15,23,42,0.4)] backdrop-blur-xl">
            <div className="grid gap-3">
              {navItems.map((item) => (
                <NavItem
                  key={item.label}
                  item={item}
                  docsSoonLabel={docsSoonLabel}
                  mobile
                  onNavigate={() => setOpen(false)}
                />
              ))}
            </div>
            <div className="mt-4 grid gap-3 sm:grid-cols-2">
              <Button asChild variant="outline" size="lg" className="rounded-2xl">
                <Link href="/login" onClick={() => setOpen(false)}>
                  {loginLabel}
                </Link>
              </Button>
              <Button asChild size="lg" className="rounded-2xl bg-slate-950 text-white hover:bg-slate-900">
                <Link href="/register" onClick={() => setOpen(false)}>
                  {signUpLabel}
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
