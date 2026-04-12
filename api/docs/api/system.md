# System API

> Generated: 2026-04-12 23:33:37

## Base URL

See [API Documentation](./api.md) for environment-specific base URLs.

## Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/v1/system-features` | Get System Features system | 🔓 |
| `GET` | `/v1/setup-status` | Get Setup Status system | 🔓 |

---

## Details

### GET `/v1/system-features`

**Get System Features system**

| Property | Value |
|----------|-------|
| Auth | 🔓 Not required |

#### Example

```bash
curl -X GET 'http://localhost:8025/api/v1/system-features'
```

---

### GET `/v1/setup-status`

**Get Setup Status system**

| Property | Value |
|----------|-------|
| Auth | 🔓 Not required |

#### Example

```bash
curl -X GET 'http://localhost:8025/api/v1/setup-status'
```

---

