export interface MarketingNavItem {
  label: string;
  href?: string;
  placeholder?: boolean;
}

export interface MarketingHeroMockupContent {
  sidebarTitle: string;
  projectsLabel: string;
  flowsLabel: string;
  environmentsLabel: string;
  teamspacesLabel: string;
  activeProject: string;
  flowOne: string;
  flowTwo: string;
  environmentValue: string;
  teamValue: string;
  workspaceTitle: string;
  workspaceSubtitle: string;
  requestOne: string;
  requestTwo: string;
  requestThree: string;
  tokenForwarded: string;
  sessionForwarded: string;
  variableForwarded: string;
  headersForwarded: string;
  resultsTitle: string;
  statusLabel: string;
  failedCheck: string;
  failedHint: string;
  aiTitle: string;
  aiReason: string;
  aiAction: string;
  historyTitle: string;
  historyOne: string;
  historyTwo: string;
  historyThree: string;
}

export interface MarketingHeroContent {
  badge: string;
  title: string;
  description: string;
  primaryCta: string;
  secondaryCta: string;
  supportingNote: string;
  mockup: MarketingHeroMockupContent;
}

export interface MarketingLogoItem {
  name: string;
}

export type MarketingFeatureIconKey =
  | 'flows'
  | 'context'
  | 'history'
  | 'collaboration'
  | 'workflow'
  | 'diagnosis';

export interface MarketingFeatureItem {
  icon: MarketingFeatureIconKey;
  title: string;
  description: string;
}

export interface MarketingFeatureSectionContent {
  eyebrow: string;
  title: string;
  description: string;
  items: MarketingFeatureItem[];
}

export type MarketingStoryVariant = 'flow' | 'history' | 'ai';

export interface MarketingStoryMockupContent {
  title: string;
  lines: string[];
}

export interface MarketingStorySectionContent {
  id: string;
  eyebrow: string;
  title: string;
  description: string;
  points: string[];
  cta: string;
  ctaHref: string;
  variant: MarketingStoryVariant;
  mockup: MarketingStoryMockupContent;
}

export interface MarketingStatItem {
  value: string;
  label: string;
  detail: string;
}

export interface MarketingStatsContent {
  eyebrow: string;
  title: string;
  description: string;
  items: MarketingStatItem[];
}

export interface MarketingFinalCtaContent {
  eyebrow: string;
  title: string;
  description: string;
  primaryCta: string;
  secondaryCta: string;
  pricingHint: string;
}

export interface MarketingFooterLink {
  label: string;
  href?: string;
  placeholder?: boolean;
}

export interface MarketingFooterColumn {
  title: string;
  links: MarketingFooterLink[];
}

export interface MarketingFooterContent {
  tagline: string;
  columns: MarketingFooterColumn[];
  socialsTitle: string;
  socials: MarketingFooterLink[];
}

export interface MarketingChromeContent {
  brandName: string;
  navItems: MarketingNavItem[];
  loginLabel: string;
  signUpLabel: string;
  docsSoonLabel: string;
  footer: MarketingFooterContent;
}

export interface MarketingPageContent {
  hero: MarketingHeroContent;
  logosTitle: string;
  logos: MarketingLogoItem[];
  features: MarketingFeatureSectionContent;
  sections: MarketingStorySectionContent[];
  stats: MarketingStatsContent;
  finalCta: MarketingFinalCtaContent;
}
