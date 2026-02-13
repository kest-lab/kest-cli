---
name: code-review-guide
description: Comprehensive code review process, checklists, and best practices
version: 1.0.0
category: development
tags: [code-review, quality, collaboration, pr]
author: ZGO Team
updated: 2026-01-24
---

# Code Review Guide

## 📋 Purpose

This skill provides a comprehensive guide to performing effective code reviews in the ZGO project, ensuring code quality, knowledge sharing, and team collaboration.

## 🎯 When to Use

- **As a Reviewer**: When reviewing a Pull Request
- **As an Author**: Before submitting a Pull Request
- **As a Team Lead**: Setting up review standards
- **Onboarding**: Teaching new team members review practices

## ⚙️ Prerequisites

- [ ] Understanding of ZGO coding standards
- [ ] Familiarity with Git and GitHub
- [ ] Knowledge of the module being reviewed

---

## 🔄 The 4-Phase Review Process

### Phase 1: Pre-Review (Author Self-Check) ⏱️ 5 minutes

**Before creating a PR, authors must:**

```bash
# 1. Run all automated checks
make test                    # All tests must pass
make lint                    # No linting errors
go fmt ./...                 # Code formatted

# 2. Run skill-specific validations
.agent/skills/coding-standards/scripts/verify-standards.sh <module>
.agent/skills/api-development/scripts/validate-api.sh <module>

# 3. Self-review changes
git diff main...HEAD         # Review your own changes

# 4. Update documentation
# - Update CHANGELOG.md
# - Update README.md if needed
# - Add/update code comments
```

**Pre-Review Checklist**:
- [ ] All tests passing locally
- [ ] No linter warnings
- [ ] Code formatted (gofmt)
- [ ] Self-reviewed the diff
- [ ] Documentation updated
- [ ] Commit messages are clear
- [ ] PR description filled out

---

### Phase 2: First Pass (Structure & Design) ⏱️ 10-15 minutes

**Focus**: High-level design, architecture, approach

**What to Check**:

#### 1. **PR Description Quality**
```markdown
✅ Good PR Description:

## What
Implemented user authentication with JWT tokens

## Why  
Users need secure login without storing passwords in session

## How
- Added JWT middleware
- Created auth service with token generation
- Added login/logout endpoints
- Migrated password hashing to bcrypt

## Testing
- Unit tests for auth service
- Integration tests for login flow
- Manual testing with Postman

## Screenshots
[Login flow screenshot]

## Checklist
- [x] Tests passing
- [x] Documentation updated
- [x] Breaking changes documented
```

#### 2. **Architecture & Design**

```go
// ✅ Good: Follows layered architecture
func (h *Handler) Create(c *gin.Context) {
    var req CreateRequest
    if !handler.BindJSON(c, &req) {
        return
    }
    user, err := h.service.Create(c.Request.Context(), &req)  // Uses service
    response.Success(c, ToResponse(user))
}

// ❌ Bad: Violates architecture
func (h *Handler) Create(c *gin.Context) {
    var po UserPO
    c.BindJSON(&po)
    h.db.Create(&po)  // Direct DB access in handler!
    c.JSON(200, po)
}
```

**Review Questions**:
- [ ] Does the solution fit the problem?
- [ ] Is the architecture sound?
- [ ] Are layers properly separated?
- [ ] Are there better alternatives?
- [ ] Will this scale?

**Comments to Leave**:
```markdown
🏗️ **Architecture**: This violates layer separation. Handlers should not 
access repositories directly. Please route through the service layer.

💡 **Suggestion**: Consider using the Circuit Breaker pattern for this 
external API call to prevent cascading failures.

⚠️ **Concern**: This approach will create N+1 query problem. Consider 
using eager loading or batch fetching.
```

---

### Phase 3: Deep Dive (Implementation) ⏱️ 20-30 minutes

**Focus**: Code quality, correctness, edge cases

#### 1. **Naming & Conventions**

```go
// ✅ Good naming
type UserPO struct { ... }              // PO suffix for models
type CreateUserRequest struct { ... }   // Request suffix
func GetByEmail(email string) { ... }   // Clear, specific

// ❌ Bad naming
type User struct { ... }                // Missing PO suffix in model.go
type CreateUserDTO struct { ... }       // Should be "Request"
func Get(param string) { ... }          // Vague, what does it get?
```

**Check**:
- [ ] Naming follows conventions ([coding-standards](../coding-standards/))
- [ ] Variable names are descriptive
- [ ] No single-letter variables (except i, j in loops)
- [ ] Constants are UPPER_CASE or const
- [ ] No Hungarian notation

#### 2. **Error Handling**

```go
// ✅ Good: Proper error handling
user, err := s.repo.GetByID(ctx, id)
if err != nil {
    if errors.Is(err, repository.ErrNotFound) {
        return nil, ErrUserNotFound
    }
    return nil, fmt.Errorf("failed to get user %d: %w", id, err)
}

// ❌ Bad: Ignoring errors
user, _ := s.repo.GetByID(ctx, id)  // Don't ignore errors!

// ❌ Bad: No context
user, err := s.repo.GetByID(ctx, id)
if err != nil {
    return nil, err  // No wrapping, no context
}

// ❌ Bad: Swallowing errors
user, err := s.repo.GetByID(ctx, id)
if err != nil {
    log.Println(err)  // Logged but not returned!
    return nil, nil
}
```

**Check**:
- [ ] All errors are handled
- [ ] Errors are wrapped with context
- [ ] No blank error returns (`_, _ :=`)
- [ ] Error messages are descriptive
- [ ] Uses custom errors where appropriate

#### 3. **Security**

```go
// ✅ Good: Secure password handling
hash, err := crypto.HashPassword(req.Password)
if err != nil {
    return nil, err
}
user.Password = hash  // Store hash, not plaintext

// JSON tag protects from exposure
type User struct {
    Password string `json:"-"`  // Never exposed in responses
}

// ❌ Bad: Security issues
user.Password = req.Password  // Storing plaintext!
log.Printf("User password: %s", password)  // Logging password!
db.Where("username = '" + username + "'")  // SQL injection!
```

**Check**:
- [ ] No plaintext passwords
- [ ] Sensitive data has `json:"-"` tag
- [ ] No SQL injection vulnerabilities
- [ ] Input validation present
- [ ] No secrets in code (use env vars)
- [ ] Authorization checks in place

#### 4. **Performance**

```go
// ✅ Good: Efficient query
users, err := db.Preload("Profile").Find(&users)

// ❌ Bad: N+1 query problem
users, _ := db.Find(&users)
for _, user := range users {
    profile, _ := db.Where("user_id = ?", user.ID).First(&profile)  // N queries!
}

// ✅ Good: Use pagination
paginator, err := pagination.PaginateFromContext[*domain.User](c, db)

// ❌ Bad: Load everything
var users []User
db.Find(&users)  // Could be millions of rows!
```

**Check**:
- [ ] No N+1 query problems
- [ ] Pagination used for lists
- [ ] Database indexes defined
- [ ] No unnecessary loops
- [ ] Efficient algorithms

#### 5. **Testing**

```go
// ✅ Good: Comprehensive test
func TestService_Create_Success(t *testing.T) {
    // Setup
    mockRepo := new(MockRepository)
    service := NewService(mockRepo)
    
    req := &CreateUserRequest{
        Email: "test@example.com",
        Username: "testuser",
    }
    
    // Expectations
    mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
        return user.Email == req.Email
    })).Return(nil)
    
    // Execute
    user, err := service.Create(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, req.Email, user.Email)
    mockRepo.AssertExpectations(t)
}

// ❌ Bad: Weak test
func TestCreate(t *testing.T) {
    service.Create(context.Background(), &req)  // No assertions!
}
```

**Check**:
- [ ] Tests cover happy path
- [ ] Tests cover error cases
- [ ] Tests cover edge cases
- [ ] Mocks are used appropriately
- [ ] Test names are descriptive
- [ ] No flaky tests

---

### Phase 4: Final Review (Polish) ⏱️ 5-10 minutes

**Focus**: Documentation, readability, maintainability

#### 1. **Code Comments**

```go
// ✅ Good: Helpful comments
// HashPassword generates a bcrypt hash from plaintext password.
// Returns error if password is empty or hashing fails.
func HashPassword(password string) (string, error) {
    if password == "" {
        return "", errors.New("password cannot be empty")
    }
    // Use cost=10 for balance between security and performance
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    return string(hash), err
}

// ❌ Bad: Useless comments
// Hash password
func HashPassword(password string) (string, error) {  // What does it do?
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)  // Missing validation
    return string(hash), err
}

// ❌ Bad: Commented-out code
func GetUser(id uint) {
    // user, _ := db.Find(id)
    // return user
    return db.First(id)  // Remove dead code!
}
```

**Check**:
- [ ] Complex logic has comments
- [ ] Public functions have godoc comments
- [ ] No commented-out code
- [ ] No TODO comments (create issues instead)
- [ ] Comments explain "why", not "what"

#### 2. **Swagger Documentation**

```go
// ✅ Good: Complete Swagger docs
// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user account with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "User creation request"
// @Success 201 {object} UserResponse
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 409 {object} response.ErrorResponse "Email already exists"
// @Router /api/users [post]
func (h *Handler) Create(c *gin.Context) {
    // ...
}

// ❌ Bad: Missing or incomplete Swagger
func (h *Handler) Create(c *gin.Context) {  // No docs!
    // ...
}
```

#### 3. **Code Readability**

```go
// ✅ Good: Readable
func (s *service) IsEligibleForDiscount(user *domain.User, order *domain.Order) bool {
    isPremiumMember := user.Tier == "premium"
    isLargeOrder := order.Total > 100
    isFirstOrder := user.OrderCount == 0
    
    return isPremiumMember || isLargeOrder || isFirstOrder
}

// ❌ Bad: Hard to read
func (s *service) IsEligibleForDiscount(u *domain.User, o *domain.Order) bool {
    return u.Tier == "premium" || o.Total > 100 || u.OrderCount == 0  // What does this mean?
}
```

**Check**:
- [ ] Functions are < 50 lines
- [ ] Files are < 500 lines
- [ ] No deeply nested logic (> 3 levels)
- [ ] Code is self-documenting
- [ ] No magic numbers (use constants)

---

## 📝 Review Feedback Guidelines

### ✅ Good Feedback

**1. Be Specific**
```markdown
❌ Bad: "This is wrong"
✅ Good: "The error is not being handled on line 45. This could cause a panic 
if the database connection fails. Consider wrapping with an error check."
```

**2. Be Constructive**
```markdown
❌ Bad: "This code is terrible"
✅ Good: "This approach works, but could be improved. Consider using the 
repository pattern to separate data access concerns. See module-creation 
skill for examples."
```

**3. Ask Questions**
```markdown
✅ "Could you explain why you chose to use a goroutine here? I'm concerned 
about potential race conditions."

✅ "Have you considered using the Circuit Breaker pattern for this external 
API call? See coding-standards skill section 9.2."
```

**4. Provide Context**
```markdown
✅ "According to our API standards (api-development skill), all list endpoints 
must use pagination. Can you add pagination.PaginateFromContext() here?"

✅ "This violates our naming convention. Database entities should have 'PO' 
suffix. See coding-standards Level 2 for details."
```

**5. Praise Good Work**
```markdown
✅ "Great use of the Circuit Breaker pattern here! This will prevent cascading 
failures if the payment gateway goes down."

✅ "Excellent test coverage! I appreciate the edge case tests."

✅ "This is a clean implementation of the repository pattern."
```

### ❌ Avoid

- Personal attacks or judgment
- Vague comments ("fix this", "bad code")
- Nitpicking without reasoning
- Blocking on style preferences
- Demanding changes without explanation

---

## 🎯 Priority Levels

### 🔴 MUST FIX (Blocking)

- Security vulnerabilities
- Breaking changes without migration
- Violations of architecture standards
- Failing tests
- Critical bugs
- Data loss risks

**Example**:
```markdown
🔴 **MUST FIX**: SQL injection vulnerability on line 67. User input is 
concatenated directly into query. Use parameterized queries instead:

db.Where("email = ?", email)  // ✅ Safe
not: db.Where("email = '" + email + "'")  // ❌ Vulnerable
```

### 🟡 SHOULD FIX (Important)

- Missing error handling
- Missing tests
- Performance issues
- Naming violations
- Missing documentation
- Code duplication

**Example**:
```markdown
🟡 **Should Fix**: Error is ignored on line 34. This could hide failures. 
Please add error handling:

if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

### 🟢 CONSIDER (Suggestions)

- Code style improvements
- Refactoring opportunities
- Alternative approaches
- Minor optimizations
- Nice-to-have features

**Example**:
```markdown
🟢 **Consider**: This function could be simplified using the MultiError 
pattern from coding-standards skill section 9.5. Not blocking, but would 
improve error visibility in batch operations.
```

---

## 📋 Review Checklists

### Quick Review Checklist (10 min)

For small changes (< 100 lines):

- [ ] **Functionality**: Does the code do what it claims?
- [ ] **Tests**: Are there tests? Do they pass?
- [ ] **Errors**: Are errors handled properly?
- [ ] **Security**: No obvious security issues?
- [ ] **Style**: Follows coding standards?

### Full Review Checklist (30 min)

For significant changes:

**Architecture** (5 min):
- [ ] Follows 8-file module structure
- [ ] Layer separation respected (Handler→Service→Repository)
- [ ] Uses domain entities, not POs in service layer
- [ ] Uses DTOs for API requests/responses

**Code Quality** (10 min):
- [ ] Naming follows conventions
- [ ] Error handling comprehensive
- [ ] No security vulnerabilities
- [ ] Performance considerations addressed
- [ ] No code duplication

**Testing** (5 min):
- [ ] Unit tests present and passing
- [ ] Tests cover happy path
- [ ] Tests cover error cases
- [ ] Mocks used appropriately
- [ ] Test coverage > 80%

**Documentation** (5 min):
- [ ] Swagger comments on handlers
- [ ] Complex logic has comments
- [ ] README updated if needed
- [ ] CHANGELOG.md updated

**API Standards** (5 min):
- [ ] List endpoints use pagination
- [ ] Uses `response.*` for all responses
- [ ] Proper HTTP methods (GET/POST/PATCH/DELETE)
- [ ] RESTful URL naming
- [ ] Request validation with `binding` tags

### Expert Review Checklist (60 min)

For critical or complex changes:

**Deep Architecture Review**:
- [ ] Design patterns appropriate
- [ ] Scalability considered
- [ ] Database schema optimal
- [ ] Indexes defined
- [ ] Caching strategy if needed
- [ ] Error handling patterns (Circuit Breaker, Retry)

**Security Audit**:
- [ ] Authentication/Authorization correct
- [ ] Input validation comprehensive
- [ ] SQL injection prevented
- [ ] XSS prevented
- [ ] CSRF protection if needed
- [ ] Sensitive data protected

**Performance Analysis**:
- [ ] No N+1 query problems
- [ ] Pagination implemented
- [ ] Indexes used effectively
- [ ] No memory leaks
- [ ] Goroutine safety
- [ ] Resource cleanup (defer)

---

## 🤝 PR Author Responsibilities

### Before Creating PR

1. **Self-Review**: Review your own changes first
2. **Run Checks**: All tests and linters passing
3. **Write Description**: Clear "what, why, how"
4. **Add Tests**: Ensure coverage > 80%
5. **Update Docs**: README, CHANGELOG, comments

### During Review

1. **Respond Promptly**: Within 24 hours
2. **Be Open**: Accept feedback gracefully
3. **Ask Questions**: If feedback unclear
4. **Make Changes**: Address all blocking issues
5. **Explain Decisions**: When disagreeing with feedback

### After Approval

1. **Merge Promptly**: Don't leave approved PRs open
2. **Monitor**: Watch for issues after merge
3. **Follow Up**: Fix any post-merge bugs quickly

---

## 👥 Reviewer Responsibilities

### Before Review

1. **Understand Context**: Read PR description and linked issues
2. **Check Out Code**: Review running code, not just diff
3. **Run Tests**: Verify tests pass locally
4. **Allocate Time**: Block 30-60 min for thorough review

### During Review

1. **Be Timely**: Review within 24-48 hours
2. **Be Thorough**: Check all checklist items
3. **Be Kind**: Constructive, not destructive
4. **Be Clear**: Specific, actionable feedback
5. **Be Consistent**: Follow review standards

### After Review

1. **Follow Up**: Check if author has questions
2. **Re-Review**: When changes are made
3. **Approve**: When all issues addressed
4. **Unblock**: Don't leave PRs waiting

---

## 🛠️ Tools & Automation

### GitHub Review Tools

```bash
# View PR diff locally
gh pr checkout <pr-number>
gh pr diff <pr-number>

# Leave review comments
gh pr review <pr-number> --comment -b "Review comments here"

# Approve PR
gh pr review <pr-number> --approve

# Request changes
gh pr review <pr-number> --request-changes -b "Please fix X"
```

### Automated Checks (CI/CD)

```yaml
# .github/workflows/pr.yml
name: PR Checks
on: [pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: make test
      - run: make lint
      - run: go test -coverprofile=coverage.out ./...
      - run: go tool cover -func=coverage.out | grep total
```

### Review Scripts

```bash
# Quick standards check
.agent/skills/coding-standards/scripts/verify-standards.sh <module>

# API standards check
.agent/skills/api-development/scripts/validate-api.sh <module>

# Logging standards check
.agent/skills/logging-standards/scripts/validate-logging.sh <module>
```

---

## 📚 Examples

- [**PR Template**](./templates/pull-request-template.md) - Standard PR description format
- [**Review Example**](./examples/good-review-example.md) - Example of a good code review
- [**Common Issues**](./examples/common-review-issues.md) - Frequently found problems

---

## 🔗 Related Skills

- [`coding-standards`](../coding-standards/) - Code quality standards
- [`api-development`](../api-development/) - API best practices
- [`module-creation`](../module-creation/) - Module structure
- [`logging-standards`](../logging-standards/) - Logging practices

---

## ✅ Quick Reference

**Before creating PR**:
```bash
make test && make lint
.agent/skills/coding-standards/scripts/verify-standards.sh <module>
git diff main...HEAD  # Self-review
```

**When reviewing**:
1. Phase 1: Check PR description & architecture (10 min)
2. Phase 2: Review implementation details (20 min)
3. Phase 3: Check tests & documentation (10 min)
4. Phase 4: Leave clear, constructive feedback (5 min)

**Feedback levels**:
- 🔴 **MUST FIX**: Security, breaking changes, critical bugs
- 🟡 **SHOULD FIX**: Missing tests, poor error handling
- 🟢 **CONSIDER**: Suggestions, improvements, style

---

**Version**: 1.0.0  
**Last Updated**: 2026-01-24  
**Maintainer**: ZGO Team
