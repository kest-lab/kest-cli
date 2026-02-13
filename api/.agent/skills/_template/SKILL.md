---
name: module-creation-template
description: Template for creating new skills in the ZGO project
version: 1.0.0
category: meta
tags: [template, scaffolding]
author: ZGO Team
updated: 2026-01-24
---

# Skill Template

## 📋 Purpose

[Describe what this skill does and what problem it solves]

Example: This skill guides you through creating a standardized DDD module in the ZGO framework.

## 🎯 When to Use

[Describe the scenarios when this skill should be invoked]

Example:
- Creating a new business module (e.g., User, Blog, Product)
- Need to follow ZGO's DDD layered architecture
- Want to ensure all necessary files and patterns are included

## ⚙️ Prerequisites

List all required tools, permissions, and knowledge:

- [ ] Prerequisite 1 (e.g., Go 1.21+ installed)
- [ ] Prerequisite 2 (e.g., Wire tool available)
- [ ] Prerequisite 3 (e.g., Understanding of DDD concepts)

## 🚀 Workflow Steps

Break down the task into clear, actionable steps:

### Step 1: [Step Title]

Brief description of what this step accomplishes.

**Actions**:
```bash
# Command examples
mkdir -p path/to/directory
```

**Expected outcome**:
- ✓ Outcome 1
- ✓ Outcome 2

### Step 2: [Step Title]

Continue with next step...

**Code example**:
```go
// Provide code snippets
package example

type ExampleStruct struct {
    Field string
}
```

### Step 3: Validation

Always include a validation step:

**Automated checks**:
```bash
# Run validation script
./scripts/validate-something.sh
```

**Manual verification**:
- [ ] Check 1
- [ ] Check 2
- [ ] Check 3

## 🔍 Troubleshooting

### Common Error 1: [Error Name]

**Symptom**: Description of the error message or behavior

**Cause**: Why this error happens

**Solution**:
1. Step 1 to resolve
2. Step 2 to resolve

**Prevention**: How to avoid this in the future

### Common Error 2: [Error Name]

[Same structure as above]

## 📚 Examples

### Example 1: [Scenario Name]

[Provide a complete, real-world example]

```go
// Full code example
```

### Example 2: [Scenario Name]

[Another example showing variation]

## 🔗 Related Skills

Link to other relevant skills:

- [`skill-name-1`](../skill-name-1/): Description of how it relates
- [`skill-name-2`](../skill-name-2/): Description of how it relates

## 📖 References

- [Project Documentation Link](../../docs/path)
- [External Resource Link](https://example.com)
- [Related Tool Documentation](https://example.com)

---

## 📝 Notes for Skill Authors

When creating a new skill, replace all bracketed sections with actual content.

### Checklist for Complete Skill

- [ ] YAML frontmatter complete and accurate
- [ ] Clear purpose and when-to-use sections
- [ ] Prerequisites listed with checkboxes
- [ ] Step-by-step workflow with code examples
- [ ] Troubleshooting section with common errors
- [ ] At least one complete example
- [ ] Related skills and references linked
- [ ] Scripts tested (if any)
- [ ] Examples validated (if any)

### Tips for Writing Great Skills

1. **Be specific**: Don't say "do X", say "run `command Y` to do X"
2. **Show, don't tell**: Include code examples for every concept
3. **Think errors**: What can go wrong? How to fix it?
4. **Link resources**: Don't duplicate docs, link to them
5. **Keep updated**: Review and update when patterns change
