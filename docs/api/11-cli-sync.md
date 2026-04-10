# CLI Sync API

## Overview

CLI Sync lets the Kest CLI upload API usage-derived specifications into a project with a project-scoped CLI token.

This flow has two parts:

1. A project member generates a CLI token in the Web Console or through the authenticated API.
2. The CLI stores `platform_url`, `platform_token`, and `platform_project_id` in `.kest/config.yaml`, then uploads specs with `kest sync push`.

---

## Authentication Modes

### 1. Generate Token

Use a normal user JWT and project write permission.

```http
Authorization: Bearer <jwt-token>
```

### 2. Upload Specs

Use the generated project-scoped CLI token.

```http
Authorization: Bearer <kest_pat_...>
```

The upload token is scoped to a single `project_id`. If the URL project does not match the token scope, the request is rejected.

---

## 1. Generate CLI Token

### POST /projects/:id/cli-tokens

Generate a project-scoped CLI token for sync uploads.

**Authentication**: JWT + project write access

#### Request Body

```json
{
  "name": "Payments API CLI sync",
  "scopes": ["spec:write"]
}
```

#### Response

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "kest_pat_3f3b7c...",
    "token_type": "bearer",
    "project_id": 12,
    "token_info": {
      "id": 5,
      "project_id": 12,
      "name": "Payments API CLI sync",
      "token_prefix": "kest_pat_3f3b7c12",
      "scopes": ["spec:write"],
      "created_at": "2026-04-09T10:00:00Z"
    }
  }
}
```

#### Notes

- The full token is only returned once.
- Current supported scopes are `spec:write` and `run:write`.
- `spec:write` is required for CLI spec upload.

---

## 2. CLI Spec Sync Upload

### POST /projects/:id/cli/spec-sync

Upload API specs derived from CLI history into a project.

**Authentication**: CLI token with `spec:write`

#### Request Body

```json
{
  "project_id": 12,
  "source": "cli",
  "metadata": {
    "cli_version": "0.1.0",
    "sync_time": "2026-04-09T10:00:00Z"
  },
  "specs": [
    {
      "method": "POST",
      "path": "/v1/orders",
      "title": "Create order",
      "summary": "Create order",
      "description": "Create a pending order from cart items.",
      "version": "v1",
      "request_body": {
        "description": "Order payload",
        "required": true,
        "content_type": "application/json",
        "schema": {
          "type": "object"
        }
      },
      "responses": {
        "201": {
          "description": "Order created",
          "content_type": "application/json",
          "schema": {
            "type": "object"
          }
        }
      },
      "examples": [
        {
          "name": "Create order",
          "request_headers": {
            "Authorization": "Bearer secret"
          },
          "request_body": "{\"cart_id\":\"cart_1\",\"password\":\"secret\"}",
          "response_status": 201,
          "response_body": "{\"data\":{\"id\":\"ord_1\"}}",
          "duration_ms": 132
        }
      ]
    }
  ]
}
```

#### Response

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "created": 1,
    "updated": 0,
    "skipped": 0,
    "errors": []
  }
}
```

#### Sanitization

During sync, Kest redacts common secrets before storing uploaded examples:

- Sensitive headers such as `Authorization`, `Cookie`, and `X-API-Key`
- Common secret-shaped JSON fields such as `password`, `token`, `secret`, `api_key`, and `client_secret`

---

## CLI Configuration

Run inside a Kest project so the CLI writes to `.kest/config.yaml`:

```bash
kest sync config \
  --platform-url "https://api.kest.dev/v1" \
  --platform-token "kest_pat_..." \
  --project-id "12"
```

Then push local history:

```bash
kest sync push
```
