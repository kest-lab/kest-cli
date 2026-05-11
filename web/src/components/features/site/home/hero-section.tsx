import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { ProductPreviewMockup } from './product-preview-mockup';
import type { MarketingHeroContent } from './types';

/**
 * @component HeroSection
 * @category Feature
 * @status Stable
 * @description Renders the hero copy and the main product interface preview for the homepage.
 * @usage Place at the top of the marketing homepage.
 * @example
 * <HeroSection content={hero} />
 */
export interface HeroSectionProps {
  content: MarketingHeroContent;
}

export function HeroSection({ content }: HeroSectionProps) {
  return (
    <section id="product" className="relative overflow-hidden bg-bg-canvas py-[4.5rem] sm:py-24 lg:py-28">
      <div className="container relative">
        <div className="flex flex-col gap-12 lg:gap-14">
          <div className="max-w-5xl">
            <div className="figma-caption inline-flex items-center rounded-full bg-highlight px-3 py-1.5 text-text-main">
              {content.badge}
            </div>

            <h1 className="figma-display-xl mt-6 max-w-5xl text-balance text-text-main">
              {content.title}
            </h1>
            <p className="figma-body-lg mt-6 max-w-3xl text-text-subtle">
              {content.description}
            </p>
            <div className="mt-8 flex flex-wrap items-center gap-3">
              <Button asChild size="2xl" className="bg-primary text-primary-foreground hover:bg-primary/95">
                <Link href="/register" className="inline-flex items-center gap-2 whitespace-nowrap">
                  <span>{content.primaryCta}</span>
                </Link>
              </Button>
              <Button
                type="button"
                variant="outline"
                size="2xl"
                className="border-border-strong bg-bg-canvas text-text-main hover:bg-bg-soft"
              >
                {content.secondaryCta}
              </Button>
            </div>
            <p className="mt-6 max-w-2xl text-sm leading-7 text-text-muted">{content.supportingNote}</p>
          </div>

          <ProductPreviewMockup variant="hero" content={content.mockup} />
        </div>
      </div>
    </section>
  );
}
