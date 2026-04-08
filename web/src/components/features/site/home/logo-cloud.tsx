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
    <section className="py-8 sm:py-12">
      <div className="container">
        <div className="rounded-[2rem] border border-border/70 bg-white/70 px-6 py-8 shadow-[0_18px_45px_-35px_rgba(15,23,42,0.28)] backdrop-blur-sm">
          <p className="text-center text-xs font-semibold uppercase tracking-[0.3em] text-text-muted">
            {title}
          </p>
          <div className="mt-7 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
            {logos.map((logo) => (
              <div
                key={logo.name}
                className="flex min-h-14 items-center justify-center rounded-2xl border border-border/70 bg-bg-canvas px-4 text-center text-sm font-semibold tracking-[0.22em] text-text-muted transition-colors duration-300 hover:text-text-subtle"
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
