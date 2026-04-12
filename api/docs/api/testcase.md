# Testcase API

> Generated: 2026-04-12 23:33:37

## Base URL

See [API Documentation](./api.md) for environment-specific base URLs.

## Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/v1/projects/:id/test-cases` | Require Project Role testcase | 🔒 |
| `POST` | `/v1/projects/:id/test-cases` | Require Project Role testcase | 🔒 |
| `GET` | `/v1/projects/:id/test-cases/:tcid` | Require Project Role testcase | 🔒 |
| `PATCH` | `/v1/projects/:id/test-cases/:tcid` | Require Project Role testcase | 🔒 |
| `DELETE` | `/v1/projects/:id/test-cases/:tcid` | Require Project Role testcase | 🔒 |
| `POST` | `/v1/projects/:id/test-cases/:tcid/duplicate` | Require Project Role testcase | 🔒 |
| `POST` | `/v1/projects/:id/test-cases/from-spec` | Require Project Role testcase | 🔒 |
| `POST` | `/v1/projects/:id/test-cases/:tcid/run` | Require Project Role testcase | 🔒 |

---

## Details

### GET `/v1/projects/:id/test-cases`

**Require Project Role testcase**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Example

```bash
curl -X GET 'http://localhost:8025/api/v1/projects/1/test-cases' \
  -H 'Authorization: Bearer <token>'
```

---

### POST `/v1/projects/:id/test-cases`

**Require Project Role testcase**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Example

```bash
curl -X POST 'http://localhost:8025/api/v1/projects/1/test-cases' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/projects/:id/test-cases/:tcid`

**Require Project Role testcase**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |
| `tcid` | `integer` | Resource identifier |

#### Example

```bash
curl -X GET 'http://localhost:8025/api/v1/projects/1/test-cases/1' \
  -H 'Authorization: Bearer <token>'
```

---

### PATCH `/v1/projects/:id/test-cases/:tcid`

**Require Project Role testcase**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |
| `tcid` | `integer` | Resource identifier |

#### Example

```bash
curl -X PATCH 'http://localhost:8025/api/v1/projects/1/test-cases/1' \
  -H 'Authorization: Bearer <token>'
```

---

### DELETE `/v1/projects/:id/test-cases/:tcid`

**Require Project Role testcase**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |
| `tcid` | `integer` | Resource identifier |

#### Example

```bash
curl -X DELETE 'http://localhost:8025/api/v1/projects/1/test-cases/1' \
  -H 'Authorization: Bearer <token>'
```

---

### POST `/v1/projects/:id/test-cases/:tcid/duplicate`

**Require Project Role testcase**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |
| `tcid` | `integer` | Resource identifier |

#### Example

```bash
curl -X POST 'http://localhost:8025/api/v1/projects/1/test-cases/1/duplicate' \
  -H 'Authorization: Bearer <token>'
```

---

### POST `/v1/projects/:id/test-cases/from-spec`

**Require Project Role testcase**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Example

```bash
curl -X POST 'http://localhost:8025/api/v1/projects/1/test-cases/from-spec' \
  -H 'Authorization: Bearer <token>'
```

---

### POST `/v1/projects/:id/test-cases/:tcid/run`

**Require Project Role testcase**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |
| `tcid` | `integer` | Resource identifier |

#### Example

```bash
curl -X POST 'http://localhost:8025/api/v1/projects/1/test-cases/1/run' \
  -H 'Authorization: Bearer <token>'
```

---

