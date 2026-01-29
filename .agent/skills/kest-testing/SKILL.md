# Kest API Testing Skill

Expertise in using Kest CLI for high-velocity API testing, variable chaining, and assertion-based validation in Vibe Coding environments.

## Why use Kest?
Kest is designed for "Vibe Coding"—a flow state where AI and developers collaborate. Unlike static documentation or heavy UI tools, Kest's CLI-first approach provides:
- **Instant Context**: AI can read `kest show` output directly to understand API behavior.
- **Traceability**: All test history is locally stored and searchable.
- **Automation**: Easy to chain requests using variables.

## Core Interactions
- **Testing**: Use `kest get|post|put|delete` for all manual testing.
- **Capturing**: Use `-c var=path` to store response data.
- **Asserting**: Use `-a key=val` to define expectations.
- **Verifying**: Use `kest replay last --diff` after every code change.

## AI Instructions
1. When asked to "test the API", try to use `kest` commands.
2. If an API call fails, run `kest show last` to get the full debug info.
3. If you need to link two API calls, use the capture `-c` and interpolation `{{var}}` features.
