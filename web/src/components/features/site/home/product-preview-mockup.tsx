import { Activity, Bot, CheckCircle2, FileText, Layers3, Sparkles } from 'lucide-react';
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
  inverse?: boolean;
}

function HeroMockup({ content }: { content: MarketingHeroMockupContent }) {
  return (
    <div className="whiteboard-mockup overflow-hidden rounded-xl border border-border-subtle bg-bg-canvas shadow-soft">
      <div className="flex items-center justify-between border-b border-border-subtle bg-bg-soft px-4 py-3">
        <div className="flex items-center gap-2">
          <span className="size-3 rounded-full bg-block-coral" />
          <span className="size-3 rounded-full bg-block-lime" />
          <span className="size-3 rounded-full bg-block-mint" />
        </div>
        <p className="figma-caption text-text-muted">{content.workspaceTitle}</p>
        <div className="h-7 w-24 rounded-full bg-primary" />
      </div>

      <div className="relative min-h-[560px] overflow-hidden bg-bg-canvas p-5 sm:p-8">
        <div className="absolute inset-0 bg-[linear-gradient(to_right,color-mix(in_srgb,var(--border-subtle),transparent_35%)_1px,transparent_1px),linear-gradient(to_bottom,color-mix(in_srgb,var(--border-subtle),transparent_35%)_1px,transparent_1px)] bg-[length:32px_32px]" />
        <div className="relative grid gap-5 lg:grid-cols-[0.78fr_1.35fr_0.9fr]">
          <aside className="rounded-xl border border-border-subtle bg-bg-canvas p-4">
            <p className="figma-caption text-text-muted">{content.sidebarTitle}</p>
            <h3 className="mt-2 text-base font-medium text-text-main">{content.activeProject}</h3>
            <div className="mt-5 space-y-3">
              {[content.flowOne, content.flowTwo].map((flow, index) => (
                <div
                  key={flow}
                  className={cn(
                    'rounded-xl border px-3 py-3 text-sm',
                    index === 0 ? 'border-border-strong bg-[var(--miro-surface-yellow)] text-text-main' : 'border-dashed border-border-main text-text-subtle'
                  )}
                >
                  {flow}
                </div>
              ))}
            </div>
            <div className="mt-5 grid gap-3">
              <div className="rounded-xl bg-bg-surface px-3 py-3">
                <p className="figma-caption text-text-muted">{content.environmentsLabel}</p>
                <p className="mt-1 text-sm font-medium text-text-main">{content.environmentValue}</p>
              </div>
              <div className="rounded-xl bg-bg-surface px-3 py-3">
                <p className="figma-caption text-text-muted">{content.teamspacesLabel}</p>
                <p className="mt-1 text-sm font-medium text-text-main">{content.teamValue}</p>
              </div>
            </div>
          </aside>

          <section className="relative min-h-[430px] rounded-xl border border-border-subtle bg-bg-canvas p-5">
            <p className="figma-caption text-text-muted">{content.workspaceSubtitle}</p>
            <div className="mt-6 grid gap-4 sm:grid-cols-3">
              {[
                { title: content.requestOne, note: content.tokenForwarded, tone: 'bg-block-lime' },
                { title: content.requestTwo, note: content.sessionForwarded, tone: 'bg-block-pink' },
                { title: content.requestThree, note: content.variableForwarded, tone: 'bg-block-mint' },
              ].map((item) => (
                <article key={item.title} className={cn('min-h-40 rounded-[1.75rem] p-4 text-text-main', item.tone)}>
                  <p className="text-sm font-medium">{item.title}</p>
                  <p className="mt-4 text-sm leading-6 text-text-subtle">{item.note}</p>
                </article>
              ))}
            </div>
            <div className="mt-5 rounded-[1.75rem] bg-block-coral p-4">
              <p className="flex items-center gap-2 text-sm font-medium text-text-main">
                <Sparkles className="size-4" />
                {content.headersForwarded}
              </p>
            </div>
            <div className="absolute bottom-6 left-10 hidden h-20 w-40 rotate-[-4deg] rounded-[1.75rem] bg-[var(--miro-orange-light)] p-4 text-sm leading-5 text-text-main md:block">
              {content.statusLabel}: {content.failedCheck}
            </div>
          </section>

          <aside className="space-y-4">
            <div className="rounded-[1.75rem] border border-border-subtle bg-bg-canvas p-4">
              <p className="figma-caption text-text-muted">{content.resultsTitle}</p>
              <p className="mt-3 text-sm font-medium leading-6 text-text-main">{content.failedHint}</p>
            </div>
            <div className="rounded-[1.75rem] bg-block-cream p-4">
              <p className="flex items-center gap-2 text-sm font-medium text-text-main">
                <Bot className="size-4" />
                {content.aiTitle}
              </p>
              <p className="mt-3 text-sm leading-6 text-text-subtle">{content.aiReason}</p>
            </div>
            <div className="rounded-[1.75rem] bg-primary p-4 text-primary-foreground">
              <p className="text-sm leading-6">{content.aiAction}</p>
            </div>
          </aside>
        </div>
      </div>
    </div>
  );
}

function StoryMockup({
  content,
  variant,
  inverse = false,
}: {
  content: MarketingStoryMockupContent;
  variant: MarketingStoryVariant;
  inverse?: boolean;
}) {
  const TitleIcon = {
    flow: Layers3,
    history: Activity,
    ai: FileText,
  }[variant];

  return (
    <div className={cn('marketing-grid rounded-xl border p-5', inverse ? 'border-text-inverse/20 bg-text-inverse/10' : 'border-border-main bg-bg-canvas')}>
      <div className="mb-5 flex items-center justify-between gap-3">
        <div className="flex items-center gap-3">
        <div className={cn('flex size-10 items-center justify-center rounded-full', inverse ? 'bg-text-inverse/15 text-text-inverse' : 'bg-bg-soft text-text-main')}>
          {TitleIcon ? <TitleIcon className="size-5" /> : null}
        </div>
        <p className={cn('text-sm font-medium', inverse ? 'text-text-inverse' : 'text-text-main')}>{content.title}</p>
      </div>
        <div className={cn('flex items-center gap-2 rounded-full border px-3 py-1 text-[11px] font-medium', inverse ? 'border-text-inverse/20 bg-text-inverse/10 text-text-inverse/75' : 'border-border-main bg-bg-canvas text-text-subtle')}>
          <CheckCircle2 className="size-3.5 text-success" />
          <span className="size-1.5 rounded-full bg-success" />
        </div>
      </div>

      {variant === 'flow' ? (
        <div className="grid gap-3">
          {[0, 2, 4, 6].map((start) => (
            <div key={content.lines[start]} className={cn('rounded-xl border p-4', inverse ? 'border-text-inverse/20 bg-text-inverse/10' : 'border-border-main bg-bg-canvas')}>
              <div className="flex items-center gap-3">
                <div className={cn('size-2.5 rounded-full', inverse ? 'bg-text-inverse' : 'bg-primary')} />
                <p className={cn('font-medium', inverse ? 'text-text-inverse' : 'text-text-main')}>{content.lines[start]}</p>
              </div>
              <p className={cn('mt-2 pl-5 text-sm leading-6', inverse ? 'text-text-inverse/75' : 'text-text-subtle')}>{content.lines[start + 1]}</p>
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
                'rounded-xl border px-4 py-3.5',
                inverse
                  ? 'border-text-inverse/20 bg-text-inverse/10'
                  : index === 0
                    ? 'border-border-main bg-block-pink'
                    : 'border-border-main bg-bg-canvas'
              )}
            >
              <p className={cn('text-sm leading-6', inverse ? 'text-text-inverse' : 'text-text-main')}>{line}</p>
            </div>
          ))}
        </div>
      ) : null}

      {variant === 'ai' ? (
        <div className="overflow-hidden rounded-xl border border-border-strong bg-primary">
          <div className="flex items-center gap-2 border-b border-text-inverse/20 px-4 py-3 text-xs uppercase tracking-[0.24em] text-text-inverse/65">
            <Bot className="size-4 text-block-lime" />
            {content.title}
          </div>
          <div className="space-y-3 px-4 py-4 font-mono text-sm leading-6 text-text-inverse">
            {content.lines.map((line, index) => (
              <div key={line} className="flex gap-3">
                <span className="w-5 shrink-0 text-text-inverse/45">{index + 1}</span>
                <span>{line}</span>
              </div>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
}

export function ProductPreviewMockup({ variant, content, className, inverse = false }: ProductPreviewMockupProps) {
  return (
    <div className={cn('relative', className)}>
      {variant === 'hero' ? (
        <HeroMockup content={content as MarketingHeroMockupContent} />
      ) : (
        <StoryMockup content={content as MarketingStoryMockupContent} variant={variant} inverse={inverse} />
      )}
    </div>
  );
}
