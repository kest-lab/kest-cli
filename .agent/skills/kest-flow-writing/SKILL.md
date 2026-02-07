# Kest Flow Writing Skill

Expertise in writing structured, chained API tests using Kest Flow (.flow.md). This skill ensures flows are written according to the latest Kest syntax standards.

## Standard Syntax

Always use ` ```kest ` code blocks with Markdown headings for step names.

### 1. Step Structure

```markdown
## 1. Step Name

` ` `kest
METHOD /path/to/api
Header: value

{
  "json": "body"
}

[Captures]
var_name: json.path

[Asserts]
status == 200
body.field exists
` ` `
```

### 2. Variable Chaining
- **Capture**: `var_name: data.path` (colon separator, extracts from response)
- **Injection**: `{{var_name}}` (injects into URL, Headers, or Body)
- **Built-in**: `{{$randomInt}}` generates a random integer

### 3. Supported Assertions
```
status == 200                    # Status code
body.data.id exists              # Field existence
body.data.password !exists       # Field non-existence
body.data.username != ""         # Not-equal
body.name == "Expected"          # String equality
duration < 500ms                 # Response time (always include ms)
```

## AI Writing Rules
1. **Always use ` ```kest `**: Never use ` ```step `, ` ```http `, or ` ```flow `.
2. **Markdown headings for steps**: Use `## N. Step Name`, not `@id` or `@name`.
3. **Colon for captures**: Use `var: path` not `var = path`.
4. **Prefer relative URLs**: Always use `/api/v1/...` instead of `https://...`.
5. **Meaningful variable names**: Use `authToken`, `projectID`, not `t`, `p`.
6. **Duration assertions**: Always include `ms` suffix: `duration < 500ms`.

## Prompt Template for Generating Flows
When asked to "generate a Kest flow", follow this pattern:
1. Define the business objective.
2. Step 1: Initial state (e.g., Login/Signup) + Capture credentials.
3. Steps 2-N: Domain operations + Chaining variables.
4. Final Step: Verification of state.
5. Add descriptive Markdown between code blocks.
