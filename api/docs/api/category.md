# Category Module API

> 💡 This documentation is automatically synchronized with the source code.

## 📌 Overview

The `category` module provides the following API endpoints:

### {
  "summary": "Retrieves a list of all categories associated with a specific project, identified by the project ID. This endpoint is useful for getting an overview of how categories are organized within a given project.",
  "request_example": "",
  "response_example": [
    {
      "id": 1,
      "name": "Electronics",
      "parent_id": null,
      "description": "All electronic items",
      "sort_order": 1
    },
    {
      "id": 2,
      "name": "Laptops",
      "parent_id": 1,
      "description": "Portable computers",
      "sort_order": 2
    }
  ]
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/categories`

**Handler Implementation:**
`category.List`

---

### {
  "summary": "Creates a new category for the specified project. The category can have a name, an optional parent category, a description, and a sort order.",
  "request_example": {
    "name": "New Category",
    "parent_id": 1,
    "description": "This is a new category for the project.",
    "sort_order": 1
  },
  "response_example": {
    "id": 123,
    "name": "New Category",
    "parent_id": 1,
    "description": "This is a new category for the project.",
    "sort_order": 1,
    "created_at": "2023-10-05T14:25:30Z"
  }
}

**Endpoint:**
<kbd>POST</kbd> `/projects/:id/categories`

**Handler Implementation:**
`category.Create`

---

### {
  "summary": "This endpoint allows reordering of categories within a project by providing a list of category IDs in the desired order.",
  "request_example": {
    "category_ids": [1, 3, 2, 4]
  },
  "response_example": {
    "status": "success",
    "message": "Categories sorted successfully"
  }
}

**Endpoint:**
<kbd>PUT</kbd> `/projects/:id/categories/sort`

**Handler Implementation:**
`category.Sort`

---

### {
  "summary": "This endpoint retrieves a specific category associated with a given project by its ID and the category's ID.",
  "request_example": "",
  "response_example": {
    "id": 1,
    "name": "Example Category",
    "parent_id": null,
    "description": "This is an example category for demonstration purposes.",
    "sort_order": 1
  }
}

**Endpoint:**
<kbd>GET</kbd> `/projects/:id/categories/:cid`

**Handler Implementation:**
`category.Get`

---

### {
  "summary": "Updates the details of a specific category within a project, including its name, parent category, description, and sort order. The parent_id can be set to null to remove the association with a parent category.",
  "request_example": {
    "name": "Updated Category Name",
    "parent_id": null,
    "description": "This is an updated description for the category.",
    "sort_order": 2
  },
  "response_example": {
    "id": 1,
    "name": "Updated Category Name",
    "parent_id": null,
    "description": "This is an updated description for the category.",
    "sort_order": 2,
    "created_at": "2023-10-01T12:00:00Z",
    "updated_at": "2023-10-05T14:30:00Z"
  }
}

**Endpoint:**
<kbd>PATCH</kbd> `/projects/:id/categories/:cid`

**Handler Implementation:**
`category.Update`

---

### Deletes a specific category associated with a project. The endpoint requires the project ID and the category ID to identify which category to delete.

**Endpoint:**
<kbd>DELETE</kbd> `/projects/:id/categories/:cid`

**Response Example:**
```json
{}
```

**Handler Implementation:**
`category.Delete`

---

