# Upgrade Guide: v2.0.0 → v2.1.0

This guide outlines the changes and upgrade steps for migrating from LlamaFront AI Scaffold v2.0.0 to v2.1.0.

## Overview

v2.1.0 focuses on keeping core dependencies up-to-date with the latest stable releases for improved performance, security, and stability.

**Release Date:** 2025-12-27

### Key Updates

- **Next.js**: 16.0.5 → **16.1.1**
- **React**: 19.2.0 → **19.2.3**
- **Tailwind CSS**: 4.1.17 → **4.1.18**
- **TypeScript & Tooling**: Various updates

## Breaking Changes

There are no major breaking changes in this release. All updates are minor/patch versions that should be backward compatible.

## Dependency Updates

### Core Frameworks
| Package | Old Version | New Version |
|---------|-------------|-------------|
| next | 16.0.5 | 16.1.1 |
| react | 19.2.0 | 19.2.3 |
| react-dom | 19.2.0 | 19.2.3 |

### UI & Styling
| Package | Old Version | New Version |
|---------|-------------|-------------|
| tailwindcss | 4.1.17 | 4.1.18 |
| lucide-react | 0.555.0 | 0.562.0 |
| @tabler/icons-react | 3.35.0 | 3.36.0 |

### Libraries
| Package | Old Version | New Version |
|---------|-------------|-------------|
| zod | 4.1.13 | 4.2.1 |
| react-hook-form | 7.67.0 | 7.69.0 |
| @tanstack/react-query | 5.90.11 | 5.90.12 |
| zustand | 5.0.8 | 5.0.9 |

## Verification

After upgrading, run the following commands to ensure everything is working correctly:

```bash
# 1. Update dependencies
pnpm install

# 2. Check types
pnpm type-check

# 3. Lint code
pnpm lint

# 4. Start dev server
pnpm dev
```

## Known Issues

- None reported.

---

For architectural changes, please refer to [ARCHITECTURE_GUIDE.md](./ARCHITECTURE_GUIDE.md).
