# Kest CLI - New Features Summary

## 🚀 What's New

Kest CLI now includes **4 powerful features** inspired by Hurl, making it the most intelligent API testing tool for AI-driven development (Vibe Coding).

---

## 1️⃣ Duration Assertion (Performance Testing)

Assert that your API responds within a specific time limit.

### Usage
```bash
# Assert response time < 1000ms
kest get /api/users --max-duration 1000

# If slower, it will FAIL with clear error
# ❌ Request Failed: duration assertion failed: 1234ms > 1000ms
```

### Features
- Millisecond precision
- Automatic failure on timeout
- Clear error messages
- Perfect for CI/CD performance gates

---

## 2️⃣ Retry Mechanism (Reliability)

Handle flaky or rate-limited APIs with automatic retries.

### Usage
```bash
# Retry up to 3 times with 2-second intervals
kest post /api/order -d '{"item": "book"}' --retry 3 --retry-wait 2000
```

### Output Example
```
⏱️  Retry attempt 1/3 (waiting 2000ms)...
⏱️  Retry attempt 2/3 (waiting 2000ms)...
✅ Request succeeded on retry 2
```

### Features
- Configurable retry count (0 = no retry, -1 = unlimited)
- Configurable wait interval (in milliseconds)
- Clear retry progress indicators
- Works with duration assertion and all other flags

---

## 3️⃣ Parallel Execution (Speed)

Run multiple tests concurrently for blazing-fast test suites.

### Usage
```bash
# Run tests in parallel with 8 workers
kest run tests.kest --parallel --jobs 8

# Sequential mode (default)
kest run tests.kest
```

### Performance Comparison
| Tests | Sequential | Parallel (8 workers) |
|-------|-----------|---------------------|
| 10    | ~10s      | ~1.5s               |
| 50    | ~50s      | ~7s                 |
| 100   | ~100s     | ~13s                |

### Features
- Default: 4 workers
- Configurable with `--jobs N`
- Thread-safe execution
- Automatic output synchronization

---

## 4️⃣ Test Summary (Beautiful Reporting)

Get comprehensive test reports with pass/fail statistics.

### Output Example
```
🚀 Running 6 test(s) from demo.kest
⚡ Parallel mode: 6 workers

╭─────────────────────────────────────────────────────────────────────╮
│                        TEST SUMMARY                                 │
├─────────────────────────────────────────────────────────────────────┤
│ ✓ GET      https://httpbin.org/uuid                  178ms │
│ ✓ POST     https://httpbin.org/post                  234ms │
│ ✗ GET      https://httpbin.org/delay/10             10006ms │
│     Error: duration assertion failed: 10006ms > 3000ms      │
│ ✓ GET      https://httpbin.org/headers                12ms │
│ ✓ GET      https://httpbin.org/user-agent             45ms │
│ ✓ POST     https://httpbin.org/anything              123ms │
├─────────────────────────────────────────────────────────────────────┤
│ Total: 6  │  Passed: 5  │  Failed: 1  │  Time: 10.598s │
│ Elapsed: 1.892s                                                     │
╰─────────────────────────────────────────────────────────────────────╯

✗ 1 test(s) failed
```

### Features
- Automatic for `kest run` command
- Color-coded results (green ✓, red ✗)
- Individual test durations
- Total time and elapsed time
- Error details for failed tests
- Beautiful box-drawing UI

---

## 💡 Combined Usage

All features work together seamlessly:

```bash
# Create a robust test suite
cat > api-tests.kest << EOF
# Fast endpoint - must respond in 500ms
get /api/health --max-duration 500

# Flaky payment API - retry on failure
post /api/payment -d @payment.json --retry 3 --retry-wait 1000

# Multiple user tests
get /api/users/1
get /api/users/2
get /api/users/3
EOF

# Run with full power
kest run api-tests.kest --parallel --jobs 4

# Result: 
# - Fast parallel execution
# - Automatic retries on failures
# - Performance assertions
# - Beautiful summary report
```

---

## 🎯 Perfect for CI/CD

```yaml
# .github/workflows/api-test.yml
- name: API Performance Tests
  run: |
    kest run tests.kest --parallel --jobs 8
    # Fails if any test exceeds duration or fails assertions
```

---

## 📊 Comparison with Hurl

| Feature              | Hurl | Kest CLI |
|---------------------|------|----------|
| Duration Assertion   | ✅   | ✅       |
| Retry Mechanism      | ✅   | ✅       |
| Parallel Execution   | ✅   | ✅       |
| Test Summary         | ✅   | ✅       |
| gRPC Support         | ❌   | ✅       |
| Streaming Support    | ❌   | ✅       |
| AI Integration       | ❌   | ✅       |
| Variable Capture     | ✅   | ✅       |
| History & Replay     | ❌   | ✅       |

---

## 🚀 Quick Start

```bash
# Install
go install github.com/kest-lab/kest-cli/cmd/kest@latest

# Test with retry
kest get https://httpbin.org/uuid --retry 3 --max-duration 1000

# Run scenario
kest run my-tests.kest --parallel
```

---

## 📖 Documentation

- **Duration**: `--max-duration <milliseconds>`
- **Retry**: `--retry <count> --retry-wait <milliseconds>`
- **Parallel**: `kest run --parallel --jobs <workers>`
- **Summary**: Automatic for `kest run` command

---

## 🎉 Result

Kest CLI is now the **most powerful, AI-friendly API testing tool** with:
- ⚡ Performance testing built-in
- 🔄 Intelligent retry logic
- 🚀 Blazing-fast parallel execution
- 📊 Beautiful test reporting
- 🤖 Perfect for Vibe Coding
