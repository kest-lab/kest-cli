import type { MarketingMessages } from './zh-Hans';

const messages: MarketingMessages = {
  meta: {
    title: 'kest',
    description: 'Build smarter API test workflows with context, history, and AI-powered diagnosis.',
  },
  brand: {
    name: 'kest',
    tagline: 'Open-source API testing and collaboration for modern product teams.',
  },
  nav: {
    product: 'Product',
    features: 'Features',
    apiDocs: 'API Docs',
    resources: 'Resources',
    pricing: 'Pricing',
    login: 'Login',
    signUp: 'Sign Up',
    docsSoon: 'Soon',
    mobileMenu: 'Open navigation',
    closeMenu: 'Close navigation',
  },
  hero: {
    badge: 'AI-NATIVE API TEST FLOWS',
    title: 'Test APIs with context, history, and AI-powered diagnosis',
    description:
      'Build readable test flows, reuse tokens and variables across requests, inspect historical results, and help your team understand why an API test failed.',
    primaryCta: 'Get Started',
    secondaryCta: 'View API Docs',
    supportingNote: 'Open-source core, collaborative workspaces, and enterprise-grade diagnostics.',
    mockup: {
      sidebarTitle: 'Workspace',
      projectsLabel: 'Projects',
      flowsLabel: 'Flows',
      environmentsLabel: 'Environments',
      teamspacesLabel: 'Team spaces',
      activeProject: 'Payments Platform',
      flowOne: 'Auth chain',
      flowTwo: 'Checkout regression',
      environmentValue: 'Staging EU',
      teamValue: 'Core API',
      workspaceTitle: 'Context-aware test flow',
      workspaceSubtitle: 'Treat every request as one step in the same engineered workflow.',
      requestOne: 'POST /auth/login',
      requestTwo: 'GET /me',
      requestThree: 'POST /billing/preview',
      tokenForwarded: 'Bearer token forwarded automatically from the login step',
      sessionForwarded: 'Session cookie injected into downstream requests',
      variableForwarded: 'tenantId variable captured from the previous response',
      headersForwarded: 'x-trace-context propagated through the entire flow',
      resultsTitle: 'Execution result',
      statusLabel: 'Status',
      failedCheck: 'Missing invoice_session in POST /billing/preview',
      failedHint: 'The failure appears after token refresh while the upstream context still points at an outdated session.',
      aiTitle: 'AI diagnosis',
      aiReason:
        'The login step returned a new session_id, but step three still references a cached header from the old session. Remap the session variable after step two.',
      aiAction: 'Suggested fix: update session.current right after login succeeds, then let the billing step read from shared context.',
      historyTitle: 'Recent runs',
      historyOne: '2 min ago · Failed · 812ms',
      historyTwo: '18 min ago · Passed · 768ms',
      historyThree: '1 hour ago · Passed · 790ms',
    },
  },
  logos: {
    title: 'Built for modern API teams',
  },
  features: {
    eyebrow: 'Core capabilities',
    title: 'Put API testing, historical insight, and team collaboration in one workspace',
    description:
      'kest is not just a request sender. It unifies execution chains, reusable context, historical results, and team signals into one engineered testing system.',
    items: {
      flows: {
        title: 'Visual Test Flows',
        description: 'See the full request chain, step health, and dependency relationships in one visual flow.',
      },
      context: {
        title: 'Context-Aware Requests',
        description: 'Carry tokens, cookies, headers, and variables forward without manual copy-paste.',
      },
      history: {
        title: 'Historical Results',
        description: 'Review execution history, timing, state changes, and regression trends over time.',
      },
      collaboration: {
        title: 'Team Collaboration',
        description: 'Share flows, annotations, status signals, and workspaces across the whole team.',
      },
      workflow: {
        title: '.flow.md Workflow Files',
        description: 'Keep test definitions readable like docs while remaining structured enough for AI.',
      },
      diagnosis: {
        title: 'AI Failure Diagnosis',
        description: 'Explain failures with request context, response data, and the surrounding execution chain.',
      },
    },
  },
  sections: {
    flow: {
      eyebrow: 'Visualize the chain',
      title: 'See every test as a connected flow',
      description:
        'Stop reasoning about one request at a time. Inspect the full test chain, request dependencies, variable propagation, and execution status step by step.',
      cta: 'Explore visual flows',
      points: {
        one: 'Flow-based request organization',
        two: 'Dependency visibility',
        three: 'Step-by-step execution tracing',
        four: 'Token and variable propagation',
      },
      mockup: {
        title: 'Flow canvas',
        laneOne: 'Login',
        laneTwo: 'Identity',
        laneThree: 'Billing',
        laneFour: 'Regression report',
        detailOne: 'Create auth.token and session.current',
        detailTwo: 'Read profile.id and assemble tenant context',
        detailThree: 'Inject invoice_session from shared variables',
        detailFour: 'Publish failure summary and change impact',
      },
    },
    history: {
      eyebrow: 'History and collaboration',
      title: 'Track results over time and keep teams aligned',
      description:
        'Let the team inspect historical runs, failure records, release impact, and live collaboration state in one place so regressions are easier to triage together.',
      cta: 'Open shared workspaces',
      points: {
        one: 'Execution history',
        two: 'Shared workspaces',
        three: 'Team status visibility',
        four: 'Comments and collaboration signals',
      },
      mockup: {
        title: 'Team timeline',
        feedOne: 'Today 09:42 · Billing regression failed',
        feedTwo: 'Lina added a comment: session refresh timing looks inconsistent',
        feedThree: 'Marco confirmed the change came from auth-service@v2.18',
        feedFour: 'Impacted flows: checkout-flow, billing-preview',
      },
    },
    ai: {
      eyebrow: 'AI + .flow.md',
      title: 'Readable workflows for humans, diagnosable by AI',
      description:
        '.flow.md files keep test definitions readable like living documentation while exposing enough structure for AI to interpret context, explain failures, and speed up onboarding.',
      cta: 'Review workflow examples',
      points: {
        one: 'Human-readable test definitions',
        two: 'AI-readable structured workflows',
        three: 'Failure explanation with context',
        four: 'Faster debugging and onboarding',
      },
      mockup: {
        title: '.flow.md snapshot',
        lineOne: 'flow "billing-preview" uses auth.login -> user.profile -> billing.preview',
        lineTwo: 'capture response.token as session.current.token',
        lineThree: 'replay headers.authorization from session.current.token',
        lineFour: 'ai note: compare failure with last green run and auth refresh timing',
      },
    },
  },
  stats: {
    eyebrow: 'Why teams switch',
    title: 'Designed for execution-heavy API teams that need clarity at scale',
    description:
      'From context propagation to failure diagnosis, every layer is tuned for engineering speed, shared visibility, and trustworthy debugging.',
    items: {
      runs: {
        value: '10K+',
        label: 'test runs visualized',
        detail: 'Keep every execution in one place, from urgent incident traces to long-lived regressions.',
      },
      teams: {
        value: '500+',
        label: 'teams collaborating',
        detail: 'Shared workspaces, comments, and status streams keep engineering and QA aligned.',
      },
      debugging: {
        value: '90%',
        label: 'faster debugging',
        detail: 'Use context-aware history and AI explanations instead of manually replaying the chain.',
      },
      readable: {
        value: '100%',
        label: 'readable workflow files',
        detail: '.flow.md serves engineers, reviewers, and AI analyzers without extra translation layers.',
      },
    },
  },
  cta: {
    eyebrow: 'Start building',
    title: 'Build smarter API test workflows with your team',
    description:
      'Start with an open-source workflow foundation, then scale into team collaboration, execution history, and AI-powered diagnostics without changing tools.',
    primaryCta: 'Start Free',
    secondaryCta: 'Read API Docs',
    pricingHint: 'Open-source at the core, with room for stronger governance as API teams grow.',
  },
  footer: {
    product: 'Product',
    apiDocs: 'API Docs',
    resources: 'Resources',
    company: 'Company',
    legal: 'Legal',
    socialTitle: 'Social',
    links: {
      overview: 'Overview',
      features: 'Features',
      flows: 'Flows',
      docsOverview: 'Overview',
      examples: 'Examples',
      schemas: 'Schemas',
      changelog: 'Changelog',
      guides: 'Guides',
      blog: 'Blog',
      openSource: 'Open Source',
      careers: 'Careers',
      contact: 'Contact',
      privacy: 'Privacy',
      terms: 'Terms',
      security: 'Security',
      github: 'GitHub',
      discord: 'Discord',
      x: 'X',
    },
  },
};

export default messages;
