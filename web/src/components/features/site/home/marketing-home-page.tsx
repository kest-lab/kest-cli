import { FeatureCards } from './feature-cards';
import { FinalCta } from './final-cta';
import { HeroSection } from './hero-section';
import { LogoCloud } from './logo-cloud';
import { ProductStorySection } from './product-story-section';
import { StatsSection } from './stats-section';
import type { MarketingPageContent } from './types';

/**
 * @component MarketingHomePage
 * @category Feature
 * @status Stable
 * @description Composes all homepage sections for the KestFlow marketing site.
 * @usage Render from the public root page after resolving translated marketing content.
 * @example
 * <MarketingHomePage content={content} />
 */
export interface MarketingHomePageProps {
  content: MarketingPageContent;
}

export function MarketingHomePage({ content }: MarketingHomePageProps) {
  return (
    <>
      <HeroSection content={content.hero} />
      <LogoCloud title={content.logosTitle} logos={content.logos} />
      <FeatureCards content={content.features} />
      {content.sections.map((section, index) => (
        <ProductStorySection key={section.id} content={section} reverse={index % 2 === 1} />
      ))}
      <StatsSection content={content.stats} />
      <FinalCta content={content.finalCta} />
    </>
  );
}
