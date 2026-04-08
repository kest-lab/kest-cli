---
name: project-strategy
description: High-level overview of project structure, mock API architecture, and authentication flow.
---

# project-strategy

## Overview

This skill provides a high-level strategic overview of the LlamaFront AI Scaffold. It covers the overall directory structure, the mock API architecture, and the authentication flow, ensuring that developers understand the core foundations of the project.

## Project Structure

The project follows a modular structure under `src/`:
- `app/`: Next.js App Router, routes, and API handlers.
- `components/`: UI primitives (`ui/`), feature components (`features/`), and common components (`common/`).
- `services/`: Stateless API wrappers.
- `store/`: Zustand stores for global state.
- `hooks/`: Custom React hooks for business logic and server state.
- `i18n/`: Internationalization configuration and modules.
- `themes/`: Design tokens and global styles.

## Mock API Architecture
All API calls initially go through Next.js Route Handlers under `src/app/api/*`, which mock backend responses. This allows for fast development without an external backend dependency.

**Flow**: `Browser → /api/* → Next.js API Routes (Mock)`

## Authentication Flow
The project supports two token modes via `NEXT_PUBLIC_AUTH_TOKEN_MODE`:
- `basic`: Single access token.
- `refresh`: Access + refresh token pair (Default).

Protected routes are managed via `AuthGuard` in layout components.

```typescript
// app/(normal)/layout.tsx
import { AuthGuard } from '@/components/auth-guard';
export default function NormalLayout({ children }) {
  return <AuthGuard>{children}</AuthGuard>;
}
```

## Route Groups
- `(auth)`: Public auth pages (login, register).
- `(normal)`: Protected dashboard and user console.
- `(site)`: Public marketing and information pages.

> [!NOTE]
> **Production Ready**: This scaffold is designed to be easily switched from mock APIs to real backend services by updating the service layer and environment variables.
