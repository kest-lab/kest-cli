import type { ScopedTranslations } from '@/i18n/shared';
import type {
  MarketingChromeContent,
  MarketingFeatureSectionContent,
  MarketingPageContent,
  MarketingStorySectionContent,
} from './types';

type MarketingTranslator = ScopedTranslations<'marketing'>;

// 统一维护对外文档站地址，避免导航和 footer 分别写死多个链接。
const API_DOCS_URL = 'https://kest-docs.vercel.app';

const logoNames = [
  'NORTHSTACK',
  'VECTORLAB',
  'AURORA API',
  'SHIPYARD',
  'LATTICE',
  'STACKPORT',
  'TRACEGRID',
  'MERIDIAN',
  'UNIT 47',
  'ORBITAL',
];

// 组装营销站顶部导航和 footer 文案。
// 作用：把 API Docs 统一指向外部文档站，页面层只消费结构化内容。
export function buildMarketingChromeContent(t: MarketingTranslator): MarketingChromeContent {
  return {
    brandName: t('brand.name'),
    navItems: [
      { label: t('nav.product'), href: '#product' },
      { label: t('nav.features'), href: '#features' },
      { label: t('nav.apiDocs'), href: API_DOCS_URL },
      { label: t('nav.resources'), href: '#resources' },
      { label: t('nav.pricing'), href: '#pricing' },
    ],
    loginLabel: t('nav.login'),
    signUpLabel: t('nav.signUp'),
    docsSoonLabel: t('nav.docsSoon'),
    footer: {
      tagline: t('brand.tagline'),
      socialsTitle: t('footer.socialTitle'),
      columns: [
        {
          title: t('footer.product'),
          links: [
            { label: t('footer.links.overview'), href: '#product' },
            { label: t('footer.links.features'), href: '#features' },
            { label: t('footer.links.flows'), href: '#resources' },
          ],
        },
        {
          title: t('footer.apiDocs'),
          links: [
            { label: t('footer.links.docsOverview'), href: API_DOCS_URL },
            { label: t('footer.links.examples'), href: API_DOCS_URL },
            { label: t('footer.links.schemas'), href: API_DOCS_URL },
          ],
        },
        {
          title: t('footer.resources'),
          links: [
            { label: t('footer.links.guides'), href: '#resources' },
            { label: t('footer.links.changelog'), placeholder: true },
            { label: t('footer.links.blog'), placeholder: true },
          ],
        },
        {
          title: t('footer.company'),
          links: [
            { label: t('footer.links.openSource'), href: '#pricing' },
            { label: t('footer.links.careers'), placeholder: true },
            { label: t('footer.links.contact'), placeholder: true },
          ],
        },
        {
          title: t('footer.legal'),
          links: [
            { label: t('footer.links.privacy'), placeholder: true },
            { label: t('footer.links.terms'), placeholder: true },
            { label: t('footer.links.security'), placeholder: true },
          ],
        },
      ],
      socials: [
        { label: t('footer.links.github'), placeholder: true },
        { label: t('footer.links.discord'), placeholder: true },
        { label: t('footer.links.x'), placeholder: true },
      ],
    },
  };
}

function buildFeatureSection(t: MarketingTranslator): MarketingFeatureSectionContent {
  return {
    eyebrow: t('features.eyebrow'),
    title: t('features.title'),
    description: t('features.description'),
    items: [
      {
        icon: 'flows',
        title: t('features.items.flows.title'),
        description: t('features.items.flows.description'),
      },
      {
        icon: 'context',
        title: t('features.items.context.title'),
        description: t('features.items.context.description'),
      },
      {
        icon: 'history',
        title: t('features.items.history.title'),
        description: t('features.items.history.description'),
      },
      {
        icon: 'collaboration',
        title: t('features.items.collaboration.title'),
        description: t('features.items.collaboration.description'),
      },
      {
        icon: 'workflow',
        title: t('features.items.workflow.title'),
        description: t('features.items.workflow.description'),
      },
      {
        icon: 'diagnosis',
        title: t('features.items.diagnosis.title'),
        description: t('features.items.diagnosis.description'),
      },
    ],
  };
}

function buildStorySections(t: MarketingTranslator): MarketingStorySectionContent[] {
  return [
    {
      id: 'product',
      eyebrow: t('sections.flow.eyebrow'),
      title: t('sections.flow.title'),
      description: t('sections.flow.description'),
      points: [
        t('sections.flow.points.one'),
        t('sections.flow.points.two'),
        t('sections.flow.points.three'),
        t('sections.flow.points.four'),
      ],
      cta: t('sections.flow.cta'),
      ctaHref: '#features',
      variant: 'flow',
      mockup: {
        title: t('sections.flow.mockup.title'),
        lines: [
          t('sections.flow.mockup.laneOne'),
          t('sections.flow.mockup.detailOne'),
          t('sections.flow.mockup.laneTwo'),
          t('sections.flow.mockup.detailTwo'),
          t('sections.flow.mockup.laneThree'),
          t('sections.flow.mockup.detailThree'),
          t('sections.flow.mockup.laneFour'),
          t('sections.flow.mockup.detailFour'),
        ],
      },
    },
    {
      id: 'resources',
      eyebrow: t('sections.history.eyebrow'),
      title: t('sections.history.title'),
      description: t('sections.history.description'),
      points: [
        t('sections.history.points.one'),
        t('sections.history.points.two'),
        t('sections.history.points.three'),
        t('sections.history.points.four'),
      ],
      cta: t('sections.history.cta'),
      ctaHref: '/register',
      variant: 'history',
      mockup: {
        title: t('sections.history.mockup.title'),
        lines: [
          t('sections.history.mockup.feedOne'),
          t('sections.history.mockup.feedTwo'),
          t('sections.history.mockup.feedThree'),
          t('sections.history.mockup.feedFour'),
        ],
      },
    },
    {
      id: 'workflow-files',
      eyebrow: t('sections.ai.eyebrow'),
      title: t('sections.ai.title'),
      description: t('sections.ai.description'),
      points: [
        t('sections.ai.points.one'),
        t('sections.ai.points.two'),
        t('sections.ai.points.three'),
        t('sections.ai.points.four'),
      ],
      cta: t('sections.ai.cta'),
      ctaHref: '/register',
      variant: 'ai',
      mockup: {
        title: t('sections.ai.mockup.title'),
        lines: [
          t('sections.ai.mockup.lineOne'),
          t('sections.ai.mockup.lineTwo'),
          t('sections.ai.mockup.lineThree'),
          t('sections.ai.mockup.lineFour'),
        ],
      },
    },
  ];
}

export function buildMarketingPageContent(t: MarketingTranslator): MarketingPageContent {
  return {
    hero: {
      badge: t('hero.badge'),
      title: t('hero.title'),
      description: t('hero.description'),
      primaryCta: t('hero.primaryCta'),
      secondaryCta: t('hero.secondaryCta'),
      supportingNote: t('hero.supportingNote'),
      mockup: {
        sidebarTitle: t('hero.mockup.sidebarTitle'),
        projectsLabel: t('hero.mockup.projectsLabel'),
        flowsLabel: t('hero.mockup.flowsLabel'),
        environmentsLabel: t('hero.mockup.environmentsLabel'),
        teamspacesLabel: t('hero.mockup.teamspacesLabel'),
        activeProject: t('hero.mockup.activeProject'),
        flowOne: t('hero.mockup.flowOne'),
        flowTwo: t('hero.mockup.flowTwo'),
        environmentValue: t('hero.mockup.environmentValue'),
        teamValue: t('hero.mockup.teamValue'),
        workspaceTitle: t('hero.mockup.workspaceTitle'),
        workspaceSubtitle: t('hero.mockup.workspaceSubtitle'),
        requestOne: t('hero.mockup.requestOne'),
        requestTwo: t('hero.mockup.requestTwo'),
        requestThree: t('hero.mockup.requestThree'),
        tokenForwarded: t('hero.mockup.tokenForwarded'),
        sessionForwarded: t('hero.mockup.sessionForwarded'),
        variableForwarded: t('hero.mockup.variableForwarded'),
        headersForwarded: t('hero.mockup.headersForwarded'),
        resultsTitle: t('hero.mockup.resultsTitle'),
        statusLabel: t('hero.mockup.statusLabel'),
        failedCheck: t('hero.mockup.failedCheck'),
        failedHint: t('hero.mockup.failedHint'),
        aiTitle: t('hero.mockup.aiTitle'),
        aiReason: t('hero.mockup.aiReason'),
        aiAction: t('hero.mockup.aiAction'),
        historyTitle: t('hero.mockup.historyTitle'),
        historyOne: t('hero.mockup.historyOne'),
        historyTwo: t('hero.mockup.historyTwo'),
        historyThree: t('hero.mockup.historyThree'),
      },
    },
    logosTitle: t('logos.title'),
    logos: logoNames.map((name) => ({ name })),
    features: buildFeatureSection(t),
    sections: buildStorySections(t),
    stats: {
      eyebrow: t('stats.eyebrow'),
      title: t('stats.title'),
      description: t('stats.description'),
      items: [
        {
          value: t('stats.items.runs.value'),
          label: t('stats.items.runs.label'),
          detail: t('stats.items.runs.detail'),
        },
        {
          value: t('stats.items.teams.value'),
          label: t('stats.items.teams.label'),
          detail: t('stats.items.teams.detail'),
        },
        {
          value: t('stats.items.debugging.value'),
          label: t('stats.items.debugging.label'),
          detail: t('stats.items.debugging.detail'),
        },
        {
          value: t('stats.items.readable.value'),
          label: t('stats.items.readable.label'),
          detail: t('stats.items.readable.detail'),
        },
      ],
    },
    finalCta: {
      eyebrow: t('cta.eyebrow'),
      title: t('cta.title'),
      description: t('cta.description'),
      primaryCta: t('cta.primaryCta'),
      secondaryCta: t('cta.secondaryCta'),
      pricingHint: t('cta.pricingHint'),
    },
  };
}
