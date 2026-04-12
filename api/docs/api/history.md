# History API

> Generated: 2026-04-12 23:33:37

## Base URL

See [API Documentation](./api.md) for environment-specific base URLs.

## Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/v1/projects/:id/history` | List historys | 🔒 |
| `GET` | `/v1/projects/:id/history/:hid` | Get history details | 🔒 |

---

## Details

### GET `/v1/projects/:id/history`

**List historys**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `history.list` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |

#### Response

```json
{
  "action": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "data": "object",
  "diff": "object",
  "entity_id": 1,
  "entity_type": "string",
  "id": 1,
  "message": "string",
  "project_id": 1,
  "user_id": 1
}
```

#### Example

```bash
curl -X GET 'http://localhost:8025/api/v1/projects/1/history' \
  -H 'Authorization: Bearer <token>'
```

---

### GET `/v1/projects/:id/history/:hid`

**Get history details**

| Property | Value |
|----------|-------|
| Auth | 🔒 JWT Required |
| Route Name | `history.show` |

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | `integer` | Resource identifier |
| `hid` | `integer` | Resource identifier |

#### Response

```json
{
  "action": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "data": "object",
  "diff": "object",
  "entity_id": 1,
  "entity_type": "string",
  "id": 1,
  "message": "string",
  "project_id": 1,
  "user_id": 1
}
```

#### Example

```bash
curl -X GET 'http://localhost:8025/api/v1/projects/1/history/1' \
  -H 'Authorization: Bearer <token>'
```

---

