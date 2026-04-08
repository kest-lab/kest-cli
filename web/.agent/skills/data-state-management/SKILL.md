---
name: data-state-management
description: Specifications for Zustand stores, React Query hooks, and the Service-Hook-Type pattern with optimistic updates.
---

# data-state-management

## Overview

This skill outlines the state management architecture of the project. It focuses on the separation of concerns between stateless services, React Query hooks for server state, and Zustand for global UI state.

## Guidelines

### 1. State Categories
- **Auth State**: `src/store/auth-store.ts` (Zustand + persist). Stores user data and system features.
- **UI State**: `src/store/ui-store.ts` (Zustand). Temporary/global UI states.
- **Server State**: Managed by React Query for all API-driven data.

### 2. Service-Hook-Type Pattern
All data handling must follow this layered architecture:

#### Layer 1: Types (`src/types/*.ts`)
Strictly type all domain models, DTOs, and query schemas.

#### Layer 2: Stateless Services (`src/services/*.ts`)
Pure functional objects that wrap HTTP requests. They do not hold state or use hooks.

```typescript
import request from '@/http/request'; 
export const userService = {
  get: () => request.get('/user'),
  update: (id: string, data: UserUpdate) => request.patch(`/user/${id}`, data)
};
```

#### Layer 3: Encapsulated Hooks (`src/hooks/*.ts`)
Manage React Query state and side effects. Implement full CRUD with optimistic updates.

```typescript
export function useUpdateExample() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }) => exampleService.update(id, data),
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: exampleKeys.detail(id) });
      const prev = queryClient.getQueryData(exampleKeys.detail(id));
      if (prev) queryClient.setQueryData(exampleKeys.detail(id), { ...prev, ...data });
      return { prev };
    },
    onError: (err, { id }, context) => {
      queryClient.setQueryData(exampleKeys.detail(id), context.prev);
      toast.error(err.message);
    },
    onSettled: (data, err, { id }) => {
      queryClient.invalidateQueries({ queryKey: exampleKeys.detail(id) });
    }
  });
}
```

### 3. Key Strategies
- **Optimistic Updates**: Provide immediate feedback.
- **Refetch-on-Failure**: Roll back state and refetch to ensure alignment.
- **Key Factories**: Use constant key objects for all query keys.

> [!TIP]
> **Stateless Services**: Always use the appropriate request instance (e.g., `request` vs `fileRequest`) from `src/http/request.ts`.
