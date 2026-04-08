import type { Metadata } from 'next';
import { getT } from '@/i18n/server';
import { buildMarketingPageContent, MarketingHomePage } from '@/components/features/site/home';

export async function generateMetadata(): Promise<Metadata> {
  const t = await getT('marketing');

  return {
    title: t('meta.title'),
    description: t('meta.description'),
  };
}

export default async function HomePage() {
  const t = await getT('marketing');
  const content = buildMarketingPageContent(t);

  return <MarketingHomePage content={content} />;
}
