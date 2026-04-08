import {
  Bot,
  FileText,
  GitBranch,
  History,
  KeyRound,
  Users,
} from 'lucide-react';
import type { MarketingFeatureIconKey, MarketingFeatureSectionContent } from './types';

/**
 * @component FeatureCards
 * @category Feature
 * @status Stable
 * @description Displays the core product capabilities in an animated card grid.
 * @usage Use on the marketing homepage to summarize the platform value proposition.
 * @example
 * <FeatureCards content={features} />
 */
export interface FeatureCardsProps {
  content: MarketingFeatureSectionContent;
}

const iconMap: Record<MarketingFeatureIconKey, typeof GitBranch> = {
  flows: GitBranch,
  context: KeyRound,
  history: History,
  collaboration: Users,
  workflow: FileText,
  diagnosis: Bot,
};

export function FeatureCards({ content }: FeatureCardsProps) {
  return (
    <section id="features" className="py-20 sm:py-24 lg:py-28">
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

        <div className="mt-14 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {content.items.map((item) => {
            const Icon = iconMap[item.icon];

            return (
              <article
                key={item.title}
                className="group marketing-panel rounded-[1.75rem] p-6 transition-all duration-300 hover:-translate-y-1.5 hover:shadow-[0_25px_65px_-32px_rgba(15,23,42,0.35)]"
              >
                <div className="flex size-12 items-center justify-center rounded-2xl bg-[color:var(--marketing-accent-soft)] text-[color:var(--marketing-accent-strong)] transition-transform duration-300 group-hover:scale-105">
                  <Icon className="size-5" />
                </div>
                <h3 className="mt-5 text-xl font-semibold text-text-main">{item.title}</h3>
                <p className="mt-3 text-sm leading-7 text-text-subtle">{item.description}</p>
              </article>
            );
          })}
        </div>
      </div>
    </section>
  );
}
