# System API

## Overview

The System module provides system-level endpoints for health checks, utilities, and administrative functions.

## Base Path

```
/v1
```

---

## 1. Health Check

### GET /health

Check the health status of the API service.

**Authentication**: Not required

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok",
    "version": "v1.0.0",
    "timestamp": "2024-02-05T03:00:00Z",
    "uptime": "72h30m15s",
    "environment": "production"
  }
}
```

---

## 2. Detailed Health Check

### GET /health/detailed

Get detailed health status including dependencies.

**Authentication**: Required (Admin access)

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "healthy",
    "version": "v1.0.0",
    "timestamp": "2024-02-05T03:00:00Z",
    "uptime": "72h30m15s",
    "environment": "production",
    "checks": {
      "database": {
        "status": "healthy",
        "response_time": "5ms",
        "details": {
          "connections": {
            "active": 15,
            "idle": 85,
            "max": 100
          }
        }
      },
      "redis": {
        "status": "healthy",
        "response_time": "2ms",
        "details": {
          "memory_usage": "45MB",
          "connected_clients": 10
        }
      },
      "queue": {
        "status": "healthy",
        "details": {
          "pending_jobs": 0,
          "failed_jobs": 0
        }
      },
      "storage": {
        "status": "healthy",
        "details": {
          "disk_usage": "45%",
          "available_space": "550GB"
        }
      }
    }
  }
}
```

---

## 3. System Information

### GET /system/info

Get system information and configuration.

**Authentication**: Required (Admin access)

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "application": {
      "name": "Kest API",
      "version": "1.0.0",
      "build": "20240205-1200",
      "go_version": "1.21.0"
    },
    "system": {
      "os": "linux",
      "architecture": "amd64",
      "cpu_cores": 8,
      "memory_total": "16GB",
      "memory_available": "8GB"
    },
    "configuration": {
      "environment": "production",
      "debug": false,
      "log_level": "info",
      "max_connections": 1000,
      "rate_limit": {
        "enabled": true,
        "requests_per_minute": 1000
      }
    },
    "features": {
      "authentication": true,
      "real_time_updates": true,
      "analytics": true,
      "advanced_permissions": true
    }
  }
}
```

---

## 4. System Metrics

### GET /system/metrics

Get system performance metrics.

**Authentication**: Required (Admin access)

#### Query Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `period` | string | ❌ No | 1h | Time period (5m, 15m, 1h, 6h, 24h) |
| `metric` | string | ❌ No | - | Specific metric name |

#### Example Request

```
GET /system/metrics?period=1h
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "1h",
    "metrics": {
      "requests": {
        "total": 15420,
        "success": 15230,
        "error": 190,
        "rate_per_second": 4.28
      },
      "response_time": {
        "avg": "125ms",
        "p50": "110ms",
        "p95": "250ms",
        "p99": "500ms"
      },
      "database": {
        "query_time_avg": "15ms",
        "connections_active": 25,
        "queries_per_second": 125
      },
      "memory": {
        "used": "2.1GB",
        "available": "13.9GB",
        "usage_percentage": 13.1
      },
      "cpu": {
        "usage_percentage": 35.2,
        "load_average": [1.2, 1.5, 1.8]
      }
    },
    "timeline": [
      {
        "timestamp": "2024-02-05T02:00:00Z",
        "requests_per_second": 4.1,
        "avg_response_time": 120,
        "error_rate": 0.01
      },
      {
        "timestamp": "2024-02-05T02:15:00Z",
        "requests_per_second": 4.5,
        "avg_response_time": 130,
        "error_rate": 0.02
      }
    ]
  }
}
```

---

## 5. System Logs

### GET /system/logs

Get system logs.

**Authentication**: Required (Admin access)

#### Query Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `level` | string | ❌ No | - | Filter by level (debug, info, warn, error) |
| `service` | string | ❌ No | - | Filter by service |
| `since` | string | ❌ No | 1h | Time since (5m, 15m, 1h, 6h, 24h) |
| `limit` | integer | ❌ No | 100 | Max number of entries |

#### Example Request

```
GET /system/logs?level=error&since=1h&limit=50
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "logs": [
      {
        "timestamp": "2024-02-05T02:45:00Z",
        "level": "error",
        "service": "api",
        "message": "Database connection failed",
        "context": {
          "error": "connection timeout",
          "query": "SELECT * FROM users",
          "duration": "30s"
        },
        "trace_id": "trace_123456789"
      },
      {
        "timestamp": "2024-02-05T02:30:00Z",
        "level": "error",
        "service": "auth",
        "message": "Invalid JWT token",
        "context": {
          "user_id": "unknown",
          "ip": "192.168.1.100"
        },
        "trace_id": "trace_987654321"
      }
    ],
    "total": 2,
    "filtered": true
  }
}
```

---

## 6. Cache Management

### POST /system/cache/clear

Clear system cache.

**Authentication**: Required (Admin access)

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | ❌ No | Cache type (all, api, auth, permissions) |

#### Example Request

```json
{
  "type": "api"
}
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Cache cleared successfully",
  "data": {
    "cleared": "api",
    "affected_keys": 150,
    "time_taken": "25ms"
  }
}
```

---

## 7. Queue Management

### GET /system/queue/status

Get queue status and statistics.

**Authentication**: Required (Admin access)

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "queues": {
      "default": {
        "pending": 5,
        "processing": 2,
        "failed": 0,
        "completed_today": 1250
      },
      "email": {
        "pending": 10,
        "processing": 1,
        "failed": 2,
        "completed_today": 500
      },
      "reports": {
        "pending": 0,
        "processing": 0,
        "failed": 0,
        "completed_today": 50
      }
    },
    "workers": {
      "active": 5,
      "total": 10,
      "busy": 2
    }
  }
}
```

### POST /system/queue/retry

Retry failed jobs in queue.

**Authentication**: Required (Admin access)

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `queue` | string | ❌ No | Queue name (all queues if not specified) |
| `limit` | integer | ❌ No | Max jobs to retry (default: 100) |

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Jobs retried successfully",
  "data": {
    "retried": 15,
    "queue": "email"
  }
}
```

---

## 8. Maintenance Mode

### POST /system/maintenance

Enable or disable maintenance mode.

**Authentication**: Required (Admin access)

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `enabled` | boolean | ✅ Yes | Enable/disable maintenance |
| `message` | string | ❌ No | Maintenance message |
| `bypass_key` | string | ❌ No | Key to bypass maintenance |

#### Example Request

```json
{
  "enabled": true,
  "message": "Scheduled maintenance in progress. Expected completion: 30 minutes.",
  "bypass_key": "maintenance_bypass_123"
}
```

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Maintenance mode updated",
  "data": {
    "enabled": true,
    "message": "Scheduled maintenance in progress. Expected completion: 30 minutes.",
    "enabled_at": "2024-02-05T03:00:00Z",
    "enabled_by": "admin_user"
  }
}
```

---

## 9. API Version Information

### GET /version

Get API version information.

**Authentication**: Not required

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "version": "1.0.0",
    "build": "20240205-1200",
    "commit": "a1b2c3d4e5f6",
    "built_at": "2024-02-05T12:00:00Z",
    "api_versions": {
      "v1": {
        "status": "current",
        "deprecated_at": null,
        "sunset_at": null
      }
    }
  }
}
```

---

## 10. Feature Flags

### GET /system/features

Get feature flag status.

**Authentication**: Required (Admin access)

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "features": {
      "real_time_updates": {
        "enabled": true,
        "description": "Enable real-time WebSocket updates"
      },
      "advanced_analytics": {
        "enabled": true,
        "description": "Advanced analytics and reporting"
      },
      "beta_api": {
        "enabled": false,
        "description": "Beta API endpoints"
      },
      "new_ui": {
        "enabled": false,
        "description": "New user interface"
      }
    }
  }
}
```

### PATCH /system/features/:feature

Toggle a feature flag.

**Authentication**: Required (Admin access)

#### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `feature` | string | ✅ Yes | Feature name |

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `enabled` | boolean | ✅ Yes | Feature status |

#### Response (200 OK)

```json
{
  "code": 0,
  "message": "Feature updated successfully",
  "data": {
    "feature": "beta_api",
    "enabled": true,
    "updated_at": "2024-02-05T03:00:00Z"
  }
}
```

---

## Usage Examples

### JavaScript (Fetch API)

```javascript
const adminToken = 'your-admin-token';

// Health check
const healthCheck = async () => {
  const response = await fetch('http://localhost:8025/health');
  return await response.json();
};

// Get system metrics
const getMetrics = async (period = '1h') => {
  const response = await fetch(`http://localhost:8025/system/metrics?period=${period}`, {
    headers: {
      'Authorization': `Bearer ${adminToken}`
    }
  });
  return await response.json();
};

// Clear cache
const clearCache = async (type = 'all') => {
  const response = await fetch('http://localhost:8025/system/cache/clear', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${adminToken}`
    },
    body: JSON.stringify({ type })
  });
  return await response.json();
};

// Enable maintenance mode
const enableMaintenance = async () => {
  const response = await fetch('http://localhost:8025/system/maintenance', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${adminToken}`
    },
    body: JSON.stringify({
      enabled: true,
      message: 'System maintenance in progress'
    })
  });
  return await response.json();
};
```

### cURL

```bash
# Health check
curl -X GET http://localhost:8025/health

# Detailed health check
curl -X GET http://localhost:8025/health/detailed \
  -H "Authorization: Bearer ADMIN_TOKEN"

# System metrics
curl -X GET "http://localhost:8025/system/metrics?period=6h" \
  -H "Authorization: Bearer ADMIN_TOKEN"

# System logs
curl -X GET "http://localhost:8025/system/logs?level=error&since=1h" \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Clear cache
curl -X POST http://localhost:8025/system/cache/clear \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{
    "type": "all"
  }'

# Enable maintenance mode
curl -X POST http://localhost:8025/system/maintenance \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{
    "enabled": true,
    "message": "System maintenance for 30 minutes"
  }'
```

---

## Security Considerations

1. **Admin Access**: Most endpoints require admin privileges
2. **Sensitive Data**: Logs may contain sensitive information
3. **Rate Limiting**: System endpoints have stricter rate limits
4. **Audit Trail**: All admin actions are logged
5. **Maintenance Bypass**: Secure bypass key should be kept secret

---

## Best Practices

1. **Monitoring**: Regularly check health endpoints
2. **Log Analysis**: Monitor error logs for issues
3. **Performance**: Track metrics for optimization
4. **Maintenance**: Schedule maintenance during low traffic
5. **Cache Management**: Clear cache after deployments
