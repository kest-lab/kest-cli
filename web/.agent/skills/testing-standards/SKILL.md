---
name: testing-standards
description: Unit and integration testing standards using Vitest and Testing Library. Use when writing tests, reviewing test coverage, setting up test infrastructure, or asking about testing patterns.
---

# testing-standards

## Overview

This skill defines testing standards for the project using Vitest and React Testing Library. It ensures consistent, maintainable tests that provide confidence in code changes.

## Guidelines

### 1. Test File Organization

- **Location**: All tests go in `src/test/` directory
- **Naming**: `[feature].test.ts` or `[component].test.tsx`
- **Setup**: Global setup in `src/test/setup.ts`

```
src/test/
├── setup.ts              # Global test configuration
├── utils.test.ts         # Utility function tests
├── hooks/                # Hook tests
│   └── use-auth.test.ts
└── components/           # Component tests
    └── auth-form.test.tsx
```

### 2. Test Commands

| Command | Purpose |
|---------|---------|
| `pnpm test` | Run tests in watch mode |
| `pnpm test:coverage` | Run tests with coverage report |

### 3. Component Testing Pattern

Use React Testing Library for component tests:

```tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect } from 'vitest';
import { LoginForm } from '@/components/features/auth/login-form';

describe('LoginForm', () => {
  it('should display error on invalid email', async () => {
    const user = userEvent.setup();
    render(<LoginForm />);
    
    await user.type(screen.getByLabelText(/email/i), 'invalid');
    await user.click(screen.getByRole('button', { name: /submit/i }));
    
    expect(screen.getByText(/invalid email/i)).toBeInTheDocument();
  });
});
```

### 4. Hook Testing Pattern

Test hooks with `renderHook`:

```tsx
import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { useLogin } from '@/hooks/use-auth';

describe('useLogin', () => {
  it('should call login API with credentials', async () => {
    const { result } = renderHook(() => useLogin());
    
    result.current.mutate({ email: 'test@example.com', password: 'pass' });
    
    await waitFor(() => {
      expect(result.current.isSuccess).toBe(true);
    });
  });
});
```

### 5. Mocking Guidelines

- **Next.js Router**: Already mocked in `setup.ts`
- **API Calls**: Use `vi.mock()` or MSW for HTTP mocking
- **Zustand Stores**: Reset state between tests

```tsx
import { vi, beforeEach } from 'vitest';
import { useAuthStore } from '@/store/auth-store';

beforeEach(() => {
  // Reset store state
  useAuthStore.getState().reset();
});

// Mock API module
vi.mock('@/services/auth', () => ({
  authApi: {
    login: vi.fn().mockResolvedValue({ id: '1', name: 'Test' }),
  },
}));
```

### 6. Naming Conventions

- **Describe blocks**: Use component/function name
- **It blocks**: Start with "should" + expected behavior
- **Test IDs**: Use `data-testid` sparingly, prefer accessible queries

```tsx
// ✅ Good - accessible query
screen.getByRole('button', { name: /submit/i });

// ⚠️ Fallback - test ID
screen.getByTestId('submit-button');
```

> [!IMPORTANT]
> **Do NOT** place test files alongside source files. All tests must be in `src/test/` directory.
