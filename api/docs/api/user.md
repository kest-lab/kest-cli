# User Module API

> 💡 This documentation is automatically synchronized with the source code.

## 📌 Overview

The `user` module provides the following API endpoints:

### {
  "description": "This endpoint allows a new user to register by providing essential details such as username, password, and email. Additional optional information like nickname and phone number can also be provided. The endpoint returns the newly created user object upon successful registration.",
  "request_example": {
    "username": "john_doe",
    "password": "SecureP@ss123",
    "email": "john.doe@example.com",
    "nickname": "JD",
    "phone": "123-456-7890"
  },
  "response_example": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "john_doe",
    "email": "john.doe@example.com",
    "nickname": "JD",
    "phone": "123-456-7890",
    "created_at": "2023-10-01T12:00:00Z"
  }
}

**Endpoint:**
<kbd>POST</kbd> `/register`

**Handler Implementation:**
`user.Register`

---

### {
  "description": "This endpoint allows a user to log in by providing their username and password. Upon successful authentication, it returns a response indicating the success of the login attempt.",
  "request_example": {
    "username": "johndoe",
    "password": "securepassword123"
  },
  "response_example": {
    "message": "Login successful",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    }
  }
}

**Endpoint:**
<kbd>POST</kbd> `/login`

**Handler Implementation:**
`user.Login`

---

### {
  "description": "This endpoint is used to initiate the password reset process for a user. It requires the user's email address, which must be valid. Upon successful submission, a password reset email will be sent to the provided email address.",
  "request_example": {
    "email": "user@example.com"
  },
  "response_example": {
    "message": "Password reset email sent"
  }
}

**Endpoint:**
<kbd>POST</kbd> `/password/reset`

**Handler Implementation:**
`user.ResetPassword`

---

### {
  "description": "The GET /users/profile endpoint retrieves the profile information of the authenticated user. It requires a valid user session or token to identify the user making the request. The endpoint returns the user's details such as nickname, avatar, phone, and bio.",
  "request_example": null,
  "response_example": {
    "status": "success",
    "data": {
      "id": "12345",
      "username": "johndoe",
      "nickname": "JohnD",
      "email": "johndoe@example.com",
      "avatar": "https://example.com/avatars/johndoe.png",
      "phone": "+1234567890",
      "bio": "Software developer with a passion for clean code and innovative solutions."
    }
  }
}

**Endpoint:**
<kbd>GET</kbd> `/users/profile`

**Handler Implementation:**
`user.GetProfile`

---

### {
  "description": "This endpoint allows an authenticated user to update their profile information, including nickname, avatar URL, phone number, and bio. The request must include the user's ID, which is typically obtained from the authentication token. Only the fields that need to be updated should be included in the request.",
  "request_example": {
    "nickname": "NewNickname123",
    "avatar": "https://example.com/avatar.jpg",
    "phone": "+1234567890",
    "bio": "Enthusiastic about technology and innovation."
  },
  "response_example": {
    "id": "user-12345",
    "username": "johndoe",
    "email": "johndoe@example.com",
    "nickname": "NewNickname123",
    "avatar": "https://example.com/avatar.jpg",
    "phone": "+1234567890",
    "bio": "Enthusiastic about technology and innovation.",
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-05T14:30:00Z"
  }
}

**Endpoint:**
<kbd>PUT</kbd> `/users/profile`

**Handler Implementation:**
`user.UpdateProfile`

---

### {
  "description": "This endpoint allows a user to change their password. The request must include the old password and the new password. The new password must be between 6 and 50 characters long. The endpoint requires the user to be authenticated, and the user ID is extracted from the context.",
  "request_example": {
    "old_password": "OldSecureP@ssw0rd",
    "new_password": "NewSecureP@ssw0rd123"
  },
  "response_example": {
    "message": "Password changed successfully"
  }
}

**Endpoint:**
<kbd>PUT</kbd> `/users/password`

**Handler Implementation:**
`user.ChangePassword`

---

### {
  "description": "This endpoint allows a user to delete their account. It requires the user to be authenticated, and the deletion process is irreversible. Upon successful deletion, the server will respond with a 204 No Content status, indicating that the operation was successful but there is no content to return. If the deletion fails, an error response will be returned.",
  "request": {
    "method": "DELETE",
    "path": "/users/account",
    "headers": {
      "Authorization": "Bearer <access_token>"
    },
    "body": "No request body is required for this endpoint."
  },
  "response": {
    "success": {
      "status_code": 204,
      "reason_phrase": "No Content",
      "body": "No content is returned upon successful account deletion."
    },
    "error": {
      "status_code": 401,
      "reason_phrase": "Unauthorized",
      "body": {
        "message": "User is not authorized to perform this action"
      }
    }
  }
}

**Endpoint:**
<kbd>DELETE</kbd> `/users/account`

**Handler Implementation:**
`user.DeleteAccount`

---

### {
  "description": "The GET /users endpoint retrieves a paginated list of users. It allows clients to fetch a specific page and number of items per page, returning the total count of available users along with the list of user details for the requested page.",
  "request_example": {
    "query_parameters": {
      "page": 1,
      "per_page": 10
    }
  },
  "response_example": {
    "data": [
      {
        "id": 1,
        "username": "john_doe",
        "nickname": "John",
        "avatar": "https://example.com/avatar.jpg",
        "phone": "+1234567890",
        "bio": "Software Developer and Tech Enthusiast"
      },
      {
        "id": 2,
        "username": "jane_smith",
        "nickname": "Jane",
        "avatar": "https://example.com/jane_avatar.jpg",
        "phone": "+0987654321",
        "bio": "Product Manager and Innovator"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 3,
      "total_items": 25,
      "items_per_page": 10,
      "next_page": "http://api.example.com/users?page=2&per_page=10",
      "prev_page": null
    }
  }
}

**Endpoint:**
<kbd>GET</kbd> `/users`

**Handler Implementation:**
`user.List`

---

### {
  "description": "This endpoint retrieves a user's information by their unique identifier. It requires the user ID to be provided as a path parameter. If the user with the given ID is found, it returns the user's details; otherwise, it returns an error indicating that the user was not found.",
  "request_example": {
    "method": "GET",
    "url": "/users/123"
  },
  "response_example_success": {
    "status_code": 200,
    "body": {
      "id": 123,
      "username": "john_doe",
      "email": "john.doe@example.com",
      "nickname": "JD",
      "phone": "+1-541-754-3010",
      "bio": "Enthusiastic developer and avid reader.",
      "avatar": "http://example.com/avatar.jpg"
    }
  },
  "response_example_error": {
    "status_code": 404,
    "body": {
      "error": "User not found"
    }
  }
}

**Endpoint:**
<kbd>GET</kbd> `/users/:id`

**Handler Implementation:**
`user.Get`

---

### {
  "description": "This endpoint retrieves detailed information about a specific user identified by the provided user ID. The response includes the user's username, email, nickname, phone number, avatar, and bio. No request body is needed for this GET request.",
  "request_example": null,
  "response_example": {
    "id": "12345",
    "username": "johndoe",
    "email": "johndoe@example.com",
    "nickname": "JD",
    "phone": "+1234567890",
    "avatar": "https://example.com/avatar.jpg",
    "bio": "Enthusiastic developer and avid reader."
  }
}

**Endpoint:**
<kbd>GET</kbd> `/users/:id/info`

**Handler Implementation:**
`user.GetUserInfo`

---

