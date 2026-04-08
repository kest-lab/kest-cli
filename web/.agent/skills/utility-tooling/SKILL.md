---
name: utility-tooling
description: The "Search First" rule and standards for creating pure utility functions and hooks with discovery tags.
---

# utility-tooling

## Overview

This skill ensures a lean and consistent codebase by enforcing the "Search First" rule before any new utility or hook is implemented. It provides standards for creating and documenting internal tools.

## Guidelines

### 1. The "Search First" Rule
Before writing a new utility function or React hook, you **MUST**:
1. **Check `src/utils/index.ts`**: Scan exports for existing utilities.
2. **Check `src/hooks/`**: Browse file names and signatures for existing logic.
3. **Check Approved Libraries**:
   - `date-fns`: All date manipulations.
   - `lodash-es`: Complex object/array operations (use sparingly).
   - `validator`: Complex string validation.

### 2. Implementation Priority
1. **Native Web APIs**: `Intl`, `URL`, `Crypto`, etc.
2. **Existing Project Utils/Hooks**: Reuse what's already available.
3. **Approved Third-Party Libraries**: Use dependencies from `package.json`.
4. **Custom Implementation**: Only if the above options are exhausted.

### 3. Documentation & Discovery
Use discovery tags in JSDoc headers:
- `@util`: Marks a pure utility function.
- `@hook`: Marks a reusable React hook.

```typescript
/**
 * @util
 * @description Format a currency value using Intl.NumberFormat.
 */
export function formatCurrency(value: number) { ... }
```

### 4. Contract for New Additions
- **Utils**: Must be pure functions in `src/utils/[category].ts`, exported via `index.ts`, with tests in `__tests__/`.
- **Hooks**: Must be in `src/hooks/use-[purpose].ts` and follow React Hook rules.

> [!TIP]
> **Minimalism**: Favor native Web APIs over external libraries whenever possible.
