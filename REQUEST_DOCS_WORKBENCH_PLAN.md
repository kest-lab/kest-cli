# Request Docs Workbench Plan

## Goal

Move AI-generated request documentation into the collection request workflow, with docs stored on `requests` and the `Docs` section shown before `Params` in the request workbench.

## Branch And Worktree Safety

- Current task branch: `codex/request-docs-in-workbench`.
- Original branch `codex/unified-runs` was behind `origin/codex/unified-runs` by 3 commits.
- Treat unrelated dirty files as user-owned unless proven otherwise.
- Suspect unrelated files at start:
  - `web/src/hooks/use-runs.ts`
  - `WORKSPACE_FIRST_V1_PLUS_6_PR_BREAKDOWN.md`

## Module Plan

1. Backend request docs module
   - Add request doc fields and migration.
   - Add request-level `gen-doc` endpoint.
   - Persist manual and AI doc source/timestamps.
   - Run Go tests for request module.
   - Commit and push.

2. Frontend request docs module
   - Add request doc fields to frontend types/service/hooks.
   - Add `Docs` section before `Params`.
   - Support manual doc edits and AI generation from saved request.
   - Run `pnpm type-check` and targeted frontend tests if needed.
   - Commit and push.

3. Integration verification
   - Rebase or otherwise align with latest upstream before final verification.
   - Use local frontend/backend.
   - Login with test account provided in the thread.
   - Verify request docs creation, AI generation, save, refresh persistence.
   - If failing, fix and repeat.
   - Capture screenshot after passing.
   - Commit and push screenshot.

4. PR process
   - Before PR creation, explicitly request user `wsdoreview`.
   - After approval, create PR with `gh`.
   - PR description must include:
     - PR Summary
     - Tests run
     - Frontend verification screenshot

## Risks To Watch

- Dirty user-owned files may be mixed in the worktree.
- Existing local servers may run different code until refreshed.
- AI generation depends on backend OpenAI configuration.
- Migration state must match local database before UI verification.
- PR screenshot should be committed into the repo so the PR body can reference it.
