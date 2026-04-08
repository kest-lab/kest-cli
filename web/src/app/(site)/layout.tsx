import { PropsWithChildren } from 'react';
import { Plus_Jakarta_Sans } from 'next/font/google';
import { getT } from '@/i18n/server';
import {
  buildMarketingChromeContent,
  MarketingFooter,
  MarketingNavbar,
} from '@/components/features/site/home';

const displayFont = Plus_Jakarta_Sans({
  subsets: ['latin'],
  variable: '--font-space-grotesk',
});

export default async function SiteLayout({ children }: PropsWithChildren) {
  const t = await getT('marketing');
  const content = buildMarketingChromeContent(t);

  return (
    <div className={`${displayFont.variable} marketing-shell flex min-h-screen flex-col bg-bg-canvas`}>
      <MarketingNavbar
        brandName={content.brandName}
        navItems={content.navItems}
        loginLabel={content.loginLabel}
        signUpLabel={content.signUpLabel}
        docsSoonLabel={content.docsSoonLabel}
        openMenuLabel={t('nav.mobileMenu')}
        closeMenuLabel={t('nav.closeMenu')}
      />
      <main className="flex-1">{children}</main>
      <MarketingFooter brandName={content.brandName} content={content.footer} />
    </div>
  );
}
