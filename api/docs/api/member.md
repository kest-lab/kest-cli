# Member Module API

> 💡 This documentation is automatically synchronized with the source code.

## 📌 Overview

The `member` module provides the following API endpoints:

### This endpoint is not provided in the given code. The provided DTOs (AddMemberRequest and UpdateMemberRequest) are for adding or updating member roles, but no handler code is given for a GET / request.

**Endpoint:**
<kbd>GET</kbd> `/`

**Handler Implementation:**
`member.unknown`

---

### {
  "summary": "This endpoint allows the creation of a new member with a specified role for a given user ID.",
  "request_example": {
    "user_id": 1,
    "role": "admin"
  },
  "response_example": ""
}

**Endpoint:**
<kbd>POST</kbd> `/`

**Handler Implementation:**
`member.Create`

---

### {
  "summary": "Updates the role of a member identified by their user ID.",
  "request_example": {
    "role": "admin"
  },
  "response_example": {
    "message": "Member role updated successfully",
    "data": {
      "user_id": 1,
      "role": "admin"
    }
  }
}

**Endpoint:**
<kbd>PATCH</kbd> `/:uid`

**Handler Implementation:**
`member.Update`

---

### This endpoint allows for the removal of a member from a group or organization, identified by their unique user ID (uid).

**Endpoint:**
<kbd>DELETE</kbd> `/:uid`

**Handler Implementation:**
`member.Delete`

---

