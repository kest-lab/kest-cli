---
name: i18n-handler
description: Guidelines for managing internationalization (i18n) in the project using next-intl and unified translation patterns.
---

# i18n-handler

## Overview

This skill provides comprehensive instructions for implementing and maintaining internationalization (i18n) within the LlamaFront AI Scaffold. It covers client and server component usage, type-safe translation keys, module organization, and scoped translations.

## Key Concepts

### Type System

The i18n system provides full TypeScript type safety through:

- **`AllTranslationKeys`**: Union type of all valid dot-notation keys (e.g., `'common.save' | 'auth.login' | ...`)
- **`Messages`**: Root type containing all module namespaces and their translations
- **`ScopedTranslations<P>`**: Type-safe scoped translator for a specific prefix path

## Guidelines

### 1. Unified Translation Pattern

The project uses `next-intl` with a unified pattern that supports both dot notation (recommended) and namespace-based access.

#### Client Components

Use the `useT` hook for translations in client-side components.

```tsx
'use client';

import { useT } from '@/i18n';

function MyComponent() {
  const t = useT();
  return (
    <div>
      {/* Dot notation (recommended) */}
      <button>{t('common.save')}</button>
      
      {/* Namespace-based (backward compatible) */}
      <button>{t.common('save')}</button>
    </div>
  );
}
```

#### Server Components

Use the asynchronous `getT` function for translations in server-side components.

```tsx
// app/page.tsx (Server Component)
import { getT } from '@/i18n';

export default async function Page() {
  const t = await getT();
  return <h1>{t('common.loading')}</h1>;
}
```

#### Scoped Translations

For components with many translations from a single namespace, use the scoped pattern:

```tsx
import { useT } from '@/i18n';

function SettingsPage() {
  // Scoped to 'settings' namespace - only settings keys are valid
  const t = useT('settings');
  return (
    <div>
      <h1>{t('title')}</h1>        {/* settings.title */}
      <p>{t('description')}</p>    {/* settings.description */}
    </div>
  );
}
```

### 2. Module Organization

Translations are organized into functional modules located in `src/i18n/modules/[module]/`.

#### Available Modules

| Module | Description |
|--------|-------------|
| `common` | Common UI text (buttons, labels, messages) |
| `auth` | Authentication-related text |
| `nav` | Navigation labels |
| `settings` | Settings page translations |
| `errors` | Error messages |
| `metadata` | Page titles and SEO metadata |
| `dashboard` | Dashboard-specific translations |
| `test` | Testing translations |

#### Supported Locales

- `zh-Hans` (Simplified Chinese) - **Base locale, defines types**
- `en-US` (English US) - Implements base type

Configuration in `src/i18n/config.ts`.

### 3. Module Structure

Each module contains three files:

```
src/i18n/modules/[module]/
├── zh-Hans.ts    # Base file (defines types)
├── en-US.ts      # English translations (implements type)
└── index.ts      # Barrel export
```

#### Base File (`zh-Hans.ts`)

```typescript
const messages = {
  title: '标题',
  description: '描述',
  nested: {
    item: '嵌套项目',
  },
};

export default messages;
export type ModuleNameMessages = typeof messages;
```

#### Translation File (`en-US.ts`)

```typescript
import type { ModuleNameMessages } from './zh-Hans';

const messages: ModuleNameMessages = {
  title: 'Title',
  description: 'Description',
  nested: {
    item: 'Nested Item',
  },
};

export default messages;
```

### 4. Adding New Translations

When adding a new page or feature:

1. Identify the relevant module in `src/i18n/modules/`.
2. Add the translation key and its values to `zh-Hans.ts` first (this defines the type).
3. Add the same key to `en-US.ts` (TypeScript will enforce this).
4. Ensure the keys are consistent across all locale files to maintain type safety.

### 5. Adding a New Module

1. Create folder: `modules/[module-name]/`
2. Create `zh-Hans.ts` with type export
3. Create `en-US.ts` implementing the type
4. Create `index.ts` barrel export
5. Register in `modules/index.ts`
6. Add to `AVAILABLE_MODULES` in `src/i18n/loader.ts`
7. Add dynamic imports to `moduleRegistry` in `loader.ts`

### 6. Using Variables (ICU Format)

```tsx
// In zh-Hans.ts:
// greeting: '你好，{name}！欢迎回来。'

const t = useT();
t('common.greeting', { name: '张三' }); // -> "你好，张三！欢迎回来。"
// or
t.common('greeting', { name: '张三' }); // -> "你好，张三！欢迎回来。"
```

## Usage Scenarios

| Task | Action |
|------|--------|
| New Page | Add translations to `src/i18n/modules/[module]/[locale].ts` and update `src/constants/routes.ts`. |
| New Component | Use `useT` (Client) or `getT` (Server) for all user-facing text. |
| Error Messages | Use the `errors` namespace and follow the standardized error handling patterns. |
| Scoped Component | Use `useT('namespace')` to get type-safe scoped translations. |

## File Reference

| File | Purpose |
|------|---------|
| `src/i18n/config.ts` | Locale configuration and settings |
| `src/i18n/index.ts` | Barrel exports for all i18n utilities |
| `src/i18n/translations.ts` | `useT` and `getT` implementation with types |
| `src/i18n/loader.ts` | Dynamic module loading and `AVAILABLE_MODULES` |
| `src/i18n/modules/index.ts` | Static type generation from modules |
| `src/i18n/README.md` | Detailed documentation with examples |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_DEFAULT_LOCALE` | Default locale | `zh-Hans` |
| `NEXT_PUBLIC_LOCALE_SWITCHER_ENABLED` | Show language switcher | `true` |

> [!IMPORTANT]
> **Zero Hardcoded Strings**: All user-facing text MUST use i18n hooks. Never hardcode text directly in JSX.

> [!TIP]
> For comprehensive examples and detailed API documentation, refer to `src/i18n/README.md`.
