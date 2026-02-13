# System Module API

> 💡 This documentation is automatically synchronized with the source code.

## 📌 Overview

The `system` module provides the following API endpoints:

### {
  "description": "The GET /system-features endpoint retrieves a set of system features and their current enablement status. This includes options such as enabling email/password login, social OAuth login, user registration, API documentation, test runner, and CLI synchronization.",
  "request_example": {},
  "response_example": {
    "EnableEmailPasswordLogin": true,
    "EnableSocialOAuthLogin": false,
    "IsAllowRegister": true,
    "EnableAPIDocumentation": true,
    "EnableTestRunner": true,
    "EnableCLISync": true
  }
}

**Endpoint:**
<kbd>GET</kbd> `/system-features`

**Handler Implementation:**
`system.GetSystemFeatures`

---

### {
  "description": "This endpoint retrieves the current setup status of the system. It provides information about whether the setup process is completed, when it was set up, if an admin user exists, and the version of the system.",
  "request": {
    "method": "GET",
    "url": "/setup-status",
    "headers": {
      "Content-Type": "application/json"
    }
  },
  "response": {
    "status": 200,
    "body": {
      "step": "finished",
      "setupAt": "2024-01-01T00:00:00Z",
      "isSetup": true,
      "hasAdmin": true,
      "version": "1.0.0"
    }
  }
}

**Endpoint:**
<kbd>GET</kbd> `/setup-status`

**Handler Implementation:**
`system.GetSetupStatus`

---

