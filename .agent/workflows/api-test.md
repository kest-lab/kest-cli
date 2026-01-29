---
description: How to test APIs using Kest CLI
---

When you modify backend code or add new API endpoints, follow these steps to verify your changes:

1. **New Endpoint**: If you've created a new endpoint, use `kest` to perform an initial test.
   - Example: `kest post /api/users -d '{"name": "test"}'`

2. **Capturing Auth**: If the endpoint requires authentication, look for a login/token endpoint and capture it.
   - Example: `kest post /login -d '{"u":"admin"}' -c "tk=token"`
   - Then use it: `kest get /protected -H "Authorization: Bearer {{tk}}"`

3. **Regression Testing**: After modifying logic in an existing endpoint, replay the last known good request to check for regressions.
// turbo
   - Run: `kest replay last --diff`
   - Examine the diff output. If changes are expected, acknowledge them. If not, fix the code.

4. **Logical Assertions**: Always add assertions to your test commands to ensure future runs detect failures.
   - Example: `kest get /health -a "status=200" -a "body.status=UP"`

5. **Contextual Debugging**: If a test fails, run `kest show last` to see the full request/response context before attempting a fix.
