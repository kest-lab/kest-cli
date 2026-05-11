import type { MarketingLogoItem } from './types';

/**
 * @component LogoCloud
 * @category Feature
 * @status Stable
 * @description Shows muted logo placeholders to communicate trust and product maturity.
 * @usage Use below the hero section on the marketing homepage.
 * @example
 * <LogoCloud title="Built for modern API teams" logos={logos} />
 */
export interface LogoCloudProps {
  title: string;
  logos: MarketingLogoItem[];
}

export function LogoCloud({ title, logos }: LogoCloudProps) {
  return (
    <section className="border-y border-border-main bg-bg-canvas py-5 text-text-main">
      <div className="container">
        <div className="flex flex-col gap-5 lg:flex-row lg:items-center lg:gap-8">
          <p className="figma-caption shrink-0 text-text-muted">
            {title}
          </p>
          <div className="grid flex-1 grid-cols-2 gap-x-6 gap-y-4 sm:grid-cols-3 lg:grid-cols-5">
            {logos.map((logo) => (
              <div
                key={logo.name}
                className="figma-caption flex min-h-9 items-center justify-center rounded-full bg-bg-subtle px-3 text-center text-text-muted transition-colors duration-200 hover:text-text-main"
              >
                {logo.name}
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
