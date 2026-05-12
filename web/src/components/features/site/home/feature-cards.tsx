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

const cardToneClasses = [
  'bg-block-pink',
  'bg-block-mint',
  'bg-block-cream',
  'bg-block-lilac',
  'bg-block-coral',
  'bg-block-lime',
];

export function FeatureCards({ content }: FeatureCardsProps) {
  return (
    <section id="features" className="bg-bg-canvas py-20 sm:py-24 lg:py-28">
      <div className="container">
        <div className="max-w-4xl">
          <p className="figma-eyebrow text-text-main">
            {content.eyebrow}
          </p>
          <h2 className="figma-display-lg mt-4 text-text-main">
            {content.title}
          </h2>
          <p className="figma-body-lg mt-5 max-w-3xl text-text-subtle">{content.description}</p>
        </div>

        <div className="mt-14 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {content.items.map((item, index) => {
            const Icon = iconMap[item.icon];

            return (
              <article
                key={item.title}
                className={`rounded-[1.75rem] border border-border-subtle p-8 ${cardToneClasses[index % cardToneClasses.length]}`}
              >
                <div className="flex size-11 items-center justify-center rounded-full border border-border-main bg-bg-canvas text-text-main">
                  <Icon className="size-5" />
                </div>
                <h3 className="figma-headline mt-5 text-text-main">{item.title}</h3>
                <p className="mt-3 text-base leading-7 text-text-subtle">{item.description}</p>
              </article>
            );
          })}
        </div>
      </div>
    </section>
  );
}
