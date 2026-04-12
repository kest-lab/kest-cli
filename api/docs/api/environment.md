# Environment API

> Generated: 2026-04-12 23:33:37

## Base URL

See [API Documentation](./api.md) for environment-specific base URLs.

## Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/v1/projects/:id/environments` | Require Project Role environment | 🔒 |
| `POST` | `/v1/projects/:id/environments` | Require Project Role environment | 🔒 |
| `GET` | `/v1/projects/:id/environments/:eid` | Require Project Role environment | 🔒 |
| `PATCH` | `/v1/projects/:id/environments/:eid` | Require Project Role environment | 🔒 |
| `DELETE` | `/v1/projects/:id/environments/:eid` | Require Project Role environment | 🔒 |
| `POST` | `/v1/projects/:id/environments/:eid/duplicate` | Require Project Role environment | 🔒 |

---

## Details

### GET `/v1/projects/:id/environments`

**Require Project Role environment**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Example

```bash
curl -X GET 'http://localhost:8025/api/v1/projects/1/environments' \
  -H 'Authorization: Bearer <token>'
```

---

### POST `/v1/projects/:id/environments`

**Require Project Role environment**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Example

```bash
curl -X POST 'http://localhost:8025/api/v1/projects/1/environments' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/projects/:id/environments/:eid`

**Require Project Role environment**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |
| `eid` | `integer` | Resource identifier |

#### Example

```bash
curl -X GET 'http://localhost:8025/api/v1/projects/1/environments/1' \
  -H 'Authorization: Bearer <token>'
```

---

### PATCH `/v1/projects/:id/environments/:eid`

**Require Project Role environment**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |
| `eid` | `integer` | Resource identifier |

#### Example

```bash
curl -X PATCH 'http://localhost:8025/api/v1/projects/1/environments/1' \
  -H 'Authorization: Bearer <token>'
```

---

### DELETE `/v1/projects/:id/environments/:eid`

**Require Project Role environment**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |
| `eid` | `integer` | Resource identifier |

#### Example

```bash
curl -X DELETE 'http://localhost:8025/api/v1/projects/1/environments/1' \
  -H 'Authorization: Bearer <token>'
```

---

### POST `/v1/projects/:id/environments/:eid/duplicate`

**Require Project Role environment**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |
| `eid` | `integer` | Resource identifier |

#### Example

```bash
curl -X POST 'http://localhost:8025/api/v1/projects/1/environments/1/duplicate' \
  -H 'Authorization: Bearer <token>'
```

---

