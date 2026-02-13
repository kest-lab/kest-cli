# Pull Request Template

## 📋 Description

### What
<!-- Briefly describe what this PR does -->

### Why
<!-- Explain the problem this PR solves or the feature it adds -->

### How
<!-- Describe the approach and key implementation details -->

---

## 🔗 Related Issues

Closes #<!-- issue number -->

---

## 🧪 Testing

### Test Coverage
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing completed

### Test Results
```bash
# Paste test output here
make test
```

### Coverage Report
```bash
# Paste coverage report here
go test -cover ./...
```

---

## 📸 Screenshots/Recordings

<!-- If applicable, add screenshots or recordings to demonstrate the changes -->

---

## ✅ Checklist

### Code Quality
- [ ] Code follows [coding standards](../.agent/skills/coding-standards/)
- [ ] Naming conventions followed
- [ ] Error handling implemented
- [ ] No security vulnerabilities
- [ ] Performance considered

### API Standards (if applicable)
- [ ] List endpoints use pagination
- [ ] Uses `response.*` for all responses
- [ ] Proper HTTP methods (GET/POST/PATCH/DELETE)
- [ ] RESTful URL naming
- [ ] Request validation with `binding` tags
- [ ] Swagger documentation added

### Testing
- [ ] All tests passing (`make test`)
- [ ] Coverage > 80%
- [ ] Edge cases tested
- [ ] Error cases tested

### Documentation
- [ ] Code comments added for complex logic
- [ ] Swagger/godoc comments added
- [ ] README updated (if needed)
- [ ] CHANGELOG.md updated

### Automation
- [ ] Linter passing (`make lint`)
- [ ] Code formatted (`go fmt ./...`)
- [ ] Wire DI regenerated (if providers changed)
- [ ] Migrations tested (if DB changes)

---

## 🔍 Self-Review

### Standards Validation
```bash
# Run these checks before submitting:
.agent/skills/coding-standards/scripts/verify-standards.sh <module>
.agent/skills/api-development/scripts/validate-api.sh <module>
.agent/skills/logging-standards/scripts/validate-logging.sh <module>
```

### What I've checked:
- [ ] Self-reviewed the code changes
- [ ] Checked for code duplication
- [ ] Verified no commented-out code
- [ ] Ensured no TODOs left (created issues instead)
- [ ] Confirmed no debugging code left

---

## 🚀 Deployment Notes

<!-- Any special deployment considerations? Database migrations? Environment variables? -->

### Breaking Changes
<!-- List any breaking changes and migration steps -->

### Configuration Changes
<!-- List any new environment variables or config updates -->

### Database Migrations
<!-- List migration files and what they do -->

---

## 📝 Additional Context

<!-- Any other information that reviewers should know -->

---

## 🙏 Reviewer Guidelines

**Estimated Review Time**: <!-- Small (10 min) / Medium (30 min) / Large (60 min) -->

**Focus Areas**: 
<!-- What should reviewers pay special attention to? -->

**Questions for Reviewers**:
<!-- Any specific questions or concerns? -->
