import { Activity, Bot, CheckCircle2, Clock3, FileText, Layers3, Sparkles } from 'lucide-react';
import { cn } from '@/utils';
import type {
  MarketingHeroMockupContent,
  MarketingStoryMockupContent,
  MarketingStoryVariant,
} from './types';

/**
 * @component ProductPreviewMockup
 * @category Feature
 * @status Stable
 * @description Renders the software-like mockups used across the marketing homepage.
 * @usage Use in the hero and alternating product sections to show context-aware flows, history, and AI diagnosis.
 * @example
 * <ProductPreviewMockup variant="hero" content={mockupContent} />
 */
export interface ProductPreviewMockupProps {
  variant: 'hero' | MarketingStoryVariant;
  content: MarketingHeroMockupContent | MarketingStoryMockupContent;
  className?: string;
}

function HeroMockup({ content }: { content: MarketingHeroMockupContent }) {
  return (
    <div className="grid gap-4 xl:grid-cols-[0.92fr_1.35fr_1fr]">
      <aside className="marketing-panel rounded-[1.75rem] p-4">
        <div className="mb-5 flex items-center justify-between">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.28em] text-text-muted">
              {content.sidebarTitle}
            </p>
            <p className="mt-2 text-sm font-semibold text-text-main">{content.activeProject}</p>
          </div>
          <div className="size-2.5 rounded-full bg-[color:var(--marketing-accent-strong)] shadow-[0_0_0_6px_rgba(251,191,36,0.12)]" />
        </div>

        <div className="space-y-4 text-sm text-text-subtle">
          <div>
            <p className="mb-2 text-[11px] font-semibold uppercase tracking-[0.25em] text-text-muted">
              {content.projectsLabel}
            </p>
            <div className="rounded-2xl border border-border/80 bg-white px-3 py-2.5 font-medium text-text-main shadow-sm">
              {content.activeProject}
            </div>
          </div>
          <div>
            <p className="mb-2 text-[11px] font-semibold uppercase tracking-[0.25em] text-text-muted">
              {content.flowsLabel}
            </p>
            <div className="space-y-2">
              <div className="rounded-2xl border border-border/80 bg-white px-3 py-2.5 shadow-sm">
                {content.flowOne}
              </div>
              <div className="rounded-2xl border border-dashed border-border px-3 py-2.5">{content.flowTwo}</div>
            </div>
          </div>
          <div className="grid gap-2 sm:grid-cols-2 xl:grid-cols-1">
            <div className="rounded-2xl border border-border/80 bg-white px-3 py-3 shadow-sm">
              <p className="text-[11px] font-semibold uppercase tracking-[0.22em] text-text-muted">
                {content.environmentsLabel}
              </p>
              <p className="mt-1 text-sm font-medium text-text-main">{content.environmentValue}</p>
            </div>
            <div className="rounded-2xl border border-border/80 bg-white px-3 py-3 shadow-sm">
              <p className="text-[11px] font-semibold uppercase tracking-[0.22em] text-text-muted">
                {content.teamspacesLabel}
              </p>
              <p className="mt-1 text-sm font-medium text-text-main">{content.teamValue}</p>
            </div>
          </div>
        </div>
      </aside>

      <section className="marketing-panel marketing-grid rounded-[1.75rem] p-5">
        <div className="mb-5 flex items-start justify-between gap-4">
          <div>
            <h3 className="text-lg font-semibold text-text-main">{content.workspaceTitle}</h3>
            <p className="mt-1 text-sm leading-6 text-text-subtle">{content.workspaceSubtitle}</p>
          </div>
          <div className="h-9 w-20 rounded-full bg-[color:var(--marketing-accent-soft)]" />
        </div>

        <div className="relative space-y-4">
          <div className="absolute left-5 top-7 hidden h-[70%] w-px bg-[linear-gradient(180deg,var(--marketing-accent-border),transparent)] md:block" />
          {[
            {
              title: content.requestOne,
              note: content.tokenForwarded,
              accent: 'bg-[color:var(--marketing-accent-strong)]',
            },
            {
              title: content.requestTwo,
              note: content.sessionForwarded,
              accent: 'bg-emerald-500',
            },
            {
              title: content.requestThree,
              note: content.variableForwarded,
              accent: 'bg-slate-900',
            },
          ].map((item, index) => (
            <div key={item.title} className="relative md:pl-8">
              <div className={cn('absolute left-0 top-5 hidden size-3 rounded-full ring-4 ring-white md:block', item.accent)} />
              <div className="rounded-[1.5rem] border border-border/80 bg-white/95 p-4 shadow-[0_18px_40px_-28px_rgba(15,23,42,0.35)] transition-transform duration-300 hover:-translate-y-1">
                <div className="flex items-center justify-between gap-3">
                  <p className="font-semibold text-text-main">{item.title}</p>
                  <span className="rounded-full bg-bg-subtle px-2.5 py-1 text-[10px] font-semibold uppercase tracking-[0.22em] text-text-muted">
                    0{index + 1}
                  </span>
                </div>
                <p className="mt-3 text-sm leading-6 text-text-subtle">{item.note}</p>
              </div>
            </div>
          ))}
          <div className="rounded-[1.5rem] border border-dashed border-[color:var(--marketing-accent-border)] bg-[color:var(--marketing-accent-soft)]/70 p-4">
            <p className="flex items-center gap-2 text-sm font-medium text-text-main">
              <Sparkles className="size-4 text-[color:var(--marketing-accent-strong)]" />
              {content.headersForwarded}
            </p>
          </div>
        </div>
      </section>

      <aside className="marketing-panel rounded-[1.75rem] p-4">
        <div className="rounded-[1.4rem] border border-border/80 bg-white p-4 shadow-sm">
          <div className="flex items-center justify-between gap-2">
            <p className="text-sm font-semibold text-text-main">{content.resultsTitle}</p>
            <span className="rounded-full bg-rose-50 px-2.5 py-1 text-[10px] font-semibold uppercase tracking-[0.24em] text-rose-600">
              {content.statusLabel}
            </span>
          </div>
          <p className="mt-3 text-sm font-semibold text-rose-700">{content.failedCheck}</p>
          <p className="mt-2 text-sm leading-6 text-text-subtle">{content.failedHint}</p>
        </div>

        <div className="mt-4 rounded-[1.4rem] border border-[color:var(--marketing-accent-border)] bg-[linear-gradient(180deg,rgba(255,255,255,0.98),rgba(255,247,237,0.98))] p-4">
          <p className="flex items-center gap-2 text-sm font-semibold text-text-main">
            <Bot className="size-4 text-[color:var(--marketing-accent-strong)]" />
            {content.aiTitle}
          </p>
          <p className="mt-3 text-sm leading-6 text-text-subtle">{content.aiReason}</p>
          <div className="mt-4 rounded-2xl bg-slate-950 px-3.5 py-3 text-sm leading-6 text-slate-100">
            {content.aiAction}
          </div>
        </div>

        <div className="mt-4 rounded-[1.4rem] border border-border/80 bg-white p-4 shadow-sm">
          <p className="text-sm font-semibold text-text-main">{content.historyTitle}</p>
          <div className="mt-3 space-y-2.5 text-sm text-text-subtle">
            {[content.historyOne, content.historyTwo, content.historyThree].map((item) => (
              <div key={item} className="flex items-center gap-2 rounded-2xl bg-bg-subtle px-3 py-2.5">
                <Clock3 className="size-4 text-text-muted" />
                <span>{item}</span>
              </div>
            ))}
          </div>
        </div>
      </aside>
    </div>
  );
}

function StoryMockup({
  content,
  variant,
}: {
  content: MarketingStoryMockupContent;
  variant: MarketingStoryVariant;
}) {
  const TitleIcon = {
    flow: Layers3,
    history: Activity,
    ai: FileText,
  }[variant];

  return (
    <div className="marketing-panel marketing-grid rounded-[2rem] p-5">
      <div className="mb-5 flex items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <div className="flex size-10 items-center justify-center rounded-2xl bg-[color:var(--marketing-accent-soft)] text-[color:var(--marketing-accent-strong)]">
            {TitleIcon ? <TitleIcon className="size-5" /> : null}
          </div>
          <p className="text-sm font-semibold text-text-main">{content.title}</p>
        </div>
        <div className="flex items-center gap-2 rounded-full border border-border/80 bg-white px-3 py-1 text-[11px] font-medium text-text-subtle shadow-sm">
          <CheckCircle2 className="size-3.5 text-emerald-500" />
          <span className="size-1.5 rounded-full bg-emerald-500" />
        </div>
      </div>

      {variant === 'flow' ? (
        <div className="grid gap-3">
          {[0, 2, 4, 6].map((start) => (
            <div key={content.lines[start]} className="rounded-[1.4rem] border border-border/80 bg-white p-4 shadow-sm">
              <div className="flex items-center gap-3">
                <div className="size-2.5 rounded-full bg-[color:var(--marketing-accent-strong)]" />
                <p className="font-medium text-text-main">{content.lines[start]}</p>
              </div>
              <p className="mt-2 pl-5 text-sm leading-6 text-text-subtle">{content.lines[start + 1]}</p>
            </div>
          ))}
        </div>
      ) : null}

      {variant === 'history' ? (
        <div className="space-y-3">
          {content.lines.map((line, index) => (
            <div
              key={line}
              className={cn(
                'rounded-[1.35rem] border px-4 py-3.5',
                index === 0
                  ? 'border-rose-200 bg-rose-50'
                  : 'border-border/80 bg-white'
              )}
            >
              <p className="text-sm leading-6 text-text-main">{line}</p>
            </div>
          ))}
        </div>
      ) : null}

      {variant === 'ai' ? (
        <div className="overflow-hidden rounded-[1.5rem] border border-slate-900/80 bg-slate-950">
          <div className="flex items-center gap-2 border-b border-slate-800 px-4 py-3 text-xs uppercase tracking-[0.24em] text-slate-400">
            <Bot className="size-4 text-[color:var(--marketing-accent)]" />
            {content.title}
          </div>
          <div className="space-y-3 px-4 py-4 font-mono text-sm leading-6 text-slate-100">
            {content.lines.map((line, index) => (
              <div key={line} className="flex gap-3">
                <span className="w-5 shrink-0 text-slate-500">{index + 1}</span>
                <span>{line}</span>
              </div>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
}

export function ProductPreviewMockup({ variant, content, className }: ProductPreviewMockupProps) {
  return (
    <div className={cn('relative', className)}>
      <div className="pointer-events-none absolute -right-8 top-8 size-24 rounded-full bg-[color:var(--marketing-accent-glow)] blur-3xl" />
      <div className="pointer-events-none absolute -left-6 bottom-8 size-24 rounded-full bg-[color:var(--marketing-accent-soft)] blur-3xl" />
      {variant === 'hero' ? (
        <HeroMockup content={content as MarketingHeroMockupContent} />
      ) : (
        <StoryMockup content={content as MarketingStoryMockupContent} variant={variant} />
      )}
    </div>
  );
}
