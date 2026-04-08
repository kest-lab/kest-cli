import type { MarketingStatsContent } from './types';

/**
 * @component StatsSection
 * @category Feature
 * @status Stable
 * @description Highlights the value metrics and engineering outcomes of the platform.
 * @usage Use near the end of the marketing homepage before the final CTA.
 * @example
 * <StatsSection content={stats} />
 */
export interface StatsSectionProps {
  content: MarketingStatsContent;
}

export function StatsSection({ content }: StatsSectionProps) {
  return (
    <section className="py-20 sm:py-24">
      <div className="container">
        <div className="mx-auto max-w-3xl text-center">
          <p className="text-xs font-semibold uppercase tracking-[0.3em] text-[color:var(--marketing-accent-strong)]">
            {content.eyebrow}
          </p>
          <h2 className="mt-4 text-4xl font-semibold tracking-tight text-text-main sm:text-5xl [font-family:var(--font-space-grotesk)]">
            {content.title}
          </h2>
          <p className="mt-5 text-lg leading-8 text-text-subtle">{content.description}</p>
        </div>

        <div className="mt-14 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
          {content.items.map((item) => (
            <article key={item.label} className="marketing-panel rounded-[1.75rem] p-6">
              <p className="text-4xl font-semibold tracking-tight text-text-main [font-family:var(--font-space-grotesk)]">
                {item.value}
              </p>
              <p className="mt-3 text-sm font-semibold uppercase tracking-[0.22em] text-[color:var(--marketing-accent-strong)]">
                {item.label}
              </p>
              <p className="mt-4 text-sm leading-7 text-text-subtle">{item.detail}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
