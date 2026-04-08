# LlamaFront Architectural Optimization Roadmap

## Phase 1: Engineering Foundation (🛡️ Critical)
The highest priority is establishing a safety net for development.

- [ ] **Setup Vitest**: Fast unit testing framework compatible with Vite/Next.js.
  - Configure `vitest.config.ts`
  - Setup `jsdom` environment
  - Add test scripts to `package.json`
- [ ] **Setup Testing Library**: Component testing utilities.
  - Install `@testing-library/react`
  - Create standard test setup/render utilities
- [ ] **CI Pipeline**: Automated checks.
  - Create `.github/workflows/ci.yml`
  - Configure Lint, Type-Check, and Test jobs

## Phase 2: Architecture Standardization (🧩 High)
Refining data patterns for consistency.

- [ ] **React Query Factories**:
  - Implement `QueryKeyFactory` pattern to centralized query keys.
  - Create typed custom hooks for API integration.
- [ ] **State Management Review**:
  - Audit `zustand` stores to ensure separation of server/client state.

## Phase 3: Developer Experience & AI (🚀 Innovative)
Enhancing the workflow for both humans and AI agents.

- [ ] **AI Context Rules**:
  - Create `.cursorrules` to guide AI code generation behavior.
  - Document architectural patterns for AI consumption.
- [ ] **Type Generation**:
  - Setup `openapi-typescript` for auto-generating API types (if backend Swagger available).

## Phase 4: Polish (🌍 Medium)
- [ ] **Design Tokens**: Formalize the theme system.
- [ ] **Strict I18n**: Type-safe translation keys.
