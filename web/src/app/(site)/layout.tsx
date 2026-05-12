import { PropsWithChildren } from 'react';
import { getT } from '@/i18n/server';
import {
  buildMarketingChromeContent,
  MarketingFooter,
  MarketingNavbar,
} from '@/components/features/site/home';

export default async function SiteLayout({ children }: PropsWithChildren) {
  const t = await getT('marketing');
  const content = buildMarketingChromeContent(t);

  return (
    <div className="marketing-shell flex min-h-screen flex-col bg-bg-canvas">
      <MarketingNavbar
        brandName={content.brandName}
        navItems={content.navItems}
        loginLabel={content.loginLabel}
        signUpLabel={content.signUpLabel}
        contactSalesLabel={content.contactSalesLabel}
        docsSoonLabel={content.docsSoonLabel}
        openMenuLabel={t('nav.mobileMenu')}
        closeMenuLabel={t('nav.closeMenu')}
      />
      <main className="flex-1">{children}</main>
      <MarketingFooter brandName={content.brandName} content={content.footer} />
    </div>
  );
}
