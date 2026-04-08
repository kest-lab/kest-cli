---
title: zgi-next
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# zgi-next

Base URLs:

# Authentication

- HTTP Authentication, scheme: bearer

# Setup

## GET status [Legacy]

GET /console/api/setup

> 返回示例

```json
{
  "step": "finished",
  "setup_at": "2025-02-26T23:05:31"
}
```

```json
{
  "step": "not_started"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» step|string|true|none||none|
|» setup_at|string|true|none||none|

## POST setup [Compatible]

POST /console/api/setup

> Body 请求参数

```json
{
  "email": "i@zgi.ai",
  "name": "zgi",
  "password": "zgi@zgi123"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» email|body|string| 是 |none|
|» name|body|string| 是 |none|
|» password|body|string| 是 |none|

> 返回示例

```json
{
  "result": "success"
}
```

```json
{
  "code": "already_setup",
  "message": "Dify has been successfully installed. Please refresh the page or return to the dashboard homepage.",
  "status": 403
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|none|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» result|string|true|none||none|

## GET system-features [Legacy]

GET /console/api/system-features

> 返回示例

```json
{
  "sso_enforced_for_signin": false,
  "sso_enforced_for_signin_protocol": "",
  "sso_enforced_for_web": false,
  "sso_enforced_for_web_protocol": "",
  "enable_web_sso_switch_component": false,
  "enable_marketplace": true,
  "max_plugin_package_size": 52428800,
  "enable_email_code_login": false,
  "enable_email_password_login": true,
  "enable_social_oauth_login": false,
  "is_allow_register": false,
  "is_allow_create_workspace": false,
  "is_email_setup": true,
  "license": {
    "status": "none",
    "expired_at": ""
  },
  "is_public_deployment": true
}
```

```json
{
  "code": "already_setup",
  "message": "Dify has been successfully installed. Please refresh the page or return to the dashboard homepage.",
  "status": 403
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|none|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» result|string|true|none||none|

# EnterpriseGroup

## POST enterprise-groups create

POST /console/api/enterprise-groups

> Body 请求参数

```json
{
  "name": "test group1",
  "short_name": ""
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» name|body|string| 是 |none|
|» short_name|body|string| 否 |none|

> 返回示例

> 201 Response

```json
{
  "message": "success"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|none|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|

## GET enterprise-groups list

GET /console/api/enterprise-groups

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|page|query|string| 否 |none|
|limit|query|string| 否 |none|
|status|query|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 20,
  "total": 1,
  "has_more": false,
  "data": [
    {
      "id": "c7efcf23-b60f-46c7-9c90-5b77422e3159",
      "name": "Default Group",
      "status": "active",
      "created_at": 1741294643
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|false|none||none|
|»» name|string|false|none||none|
|»» status|string|false|none||none|
|»» created_at|integer|false|none||none|

## PUT enterprise-groups update

PUT /console/api/enterprise-groups

> Body 请求参数

```json
{
  "id": "2d604491-df3e-491a-89b2-007ce56d1a9e",
  "name": "update test group",
  "short_name": "update group short_name"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|string| 是 |none|
|» name|body|string| 是 |none|
|» short_name|body|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 0,
  "total": 1,
  "has_more": false,
  "data": [
    {
      "id": "4d6982b5-97d0-44ee-a99c-d904b57445d8",
      "name": "Default Group",
      "status": "active",
      "created_at": 1741265662
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|false|none||none|
|»» name|string|false|none||none|
|»» status|string|false|none||none|
|»» created_at|integer|false|none||none|

## DELETE enterprise-groups delete

DELETE /console/api/enterprise-groups

> Body 请求参数

```json
{
  "id": "3e09f270-3435-46cb-8d38-5e4c4de6c679"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» id|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 0,
  "total": 1,
  "has_more": false,
  "data": [
    {
      "id": "4d6982b5-97d0-44ee-a99c-d904b57445d8",
      "name": "Default Group",
      "status": "active",
      "created_at": 1741265662
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|false|none||none|
|»» name|string|false|none||none|
|»» status|string|false|none||none|
|»» created_at|integer|false|none||none|

## GET enterprise-groups info

GET /console/api/enterprise-groups/info/{string:group_id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|string:group_id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "id": "2d604491-df3e-491a-89b2-007ce56d1a9e",
  "name": "test group",
  "short_name": "tests",
  "status": "active",
  "created_at": 1743972307
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|string|true|none||none|
|» name|string|true|none||none|
|» short_name|string|true|none||none|
|» status|string|true|none||none|
|» created_at|integer|true|none||none|

## GET tenants list by group_id

GET /console/api/enterprise-groups/{group_id}/tenants

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|group_id|path|string| 是 |none|
|page|query|string| 否 |none|
|limit|query|string| 否 |none|
|status|query|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 0,
  "total": 1,
  "has_more": false,
  "data": [
    {
      "id": "4d6982b5-97d0-44ee-a99c-d904b57445d8",
      "name": "Default Group",
      "status": "active",
      "created_at": 1741265662
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|false|none||none|
|»» name|string|false|none||none|
|»» status|string|false|none||none|
|»» created_at|integer|false|none||none|

## POST tenant create by group id

POST /console/api/enterprise-groups/{group_id}/tenants

> Body 请求参数

```json
{
  "name": "my workspace"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|group_id|path|string| 是 |none|
|body|body|object| 否 |none|
|» name|body|string| 是 |none|

> 返回示例

> 201 Response

```json
{
  "message": "workspace created."
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|none|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|

## GET unjoined-tenants by group account_id

GET /console/api/enterprise-groups/{group_id}/unjoined-tenants/{account_id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|group_id|path|string| 是 |none|
|account_id|path|string| 是 |none|
|page|query|string| 否 |none|
|limit|query|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 0,
  "total": 1,
  "has_more": false,
  "data": [
    {
      "id": "4d6982b5-97d0-44ee-a99c-d904b57445d8",
      "name": "Default Group",
      "status": "active",
      "created_at": 1741265662
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|false|none||none|
|»» name|string|false|none||none|
|»» status|string|false|none||none|
|»» created_at|integer|false|none||none|

## GET joined-tenants by group account_id

GET /console/api/enterprise-groups/{group_id}/joined-tenants/{account_id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|group_id|path|string| 是 |none|
|account_id|path|string| 是 |none|
|page|query|string| 否 |none|
|limit|query|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 0,
  "total": 1,
  "has_more": false,
  "data": [
    {
      "id": "4d6982b5-97d0-44ee-a99c-d904b57445d8",
      "name": "Default Group",
      "status": "active",
      "created_at": 1741265662
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|false|none||none|
|»» name|string|false|none||none|
|»» status|string|false|none||none|
|»» created_at|integer|false|none||none|

## GET joined-tenants roles by group account_id

GET /console/api/enterprise-groups/{group_id}/joined-tenants-roles/{account_id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|group_id|path|string| 是 |none|
|account_id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
[
  {
    "tenant_id": "9e9f5925-6ca5-455e-bf8b-e7a413a2280a",
    "tenant_name": "my workspace",
    "role": "normal",
    "position": "dev",
    "permissions": [
      "create_agent",
      "create_knowledge"
    ]
  },
  {
    "tenant_id": "be5d2163-76cd-40f9-b2ee-c8d0a0cc2f6d",
    "tenant_name": "zgi's Workspace",
    "role": "admin",
    "position": "高级主管1",
    "permissions": [
      "create_agent",
      "create_knowledge"
    ]
  }
]
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» tenant_id|string|true|none||none|
|» tenant_name|string|true|none||none|
|» role|string|true|none||none|
|» position|string|true|none||none|
|» permissions|[string]|true|none||none|

## GET joined groups by account_id

GET /console/api/enterprise-groups/joined-groups/{account_id}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|account_id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 20,
  "total": 1,
  "has_more": false,
  "data": [
    {
      "id": "c7efcf23-b60f-46c7-9c90-5b77422e3159",
      "name": "Default Group",
      "status": "active",
      "created_at": 1741294643
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|false|none||none|
|»» name|string|false|none||none|
|»» status|string|false|none||none|
|»» created_at|integer|false|none||none|

## GET joined groups current user

GET /console/api/enterprise-groups/joined-groups

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 20,
  "total": 1,
  "has_more": false,
  "data": [
    {
      "id": "c7efcf23-b60f-46c7-9c90-5b77422e3159",
      "name": "Default Group",
      "status": "active",
      "created_at": 1741294643
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|false|none||none|
|»» name|string|false|none||none|
|»» status|string|false|none||none|
|»» created_at|integer|false|none||none|

## GET managed-tenants check-manage-permission

GET /console/api/enterprise-groups/{string:group_id}/check-manage-permission

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|string:group_id|path|string| 是 |none|
|switch_tenant|query|boolean| 否 |none|

> 返回示例

> 200 Response

```json
{
  "has_permission": true
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» has_permission|boolean|true|none||none|

## GET enterprise-groups current

GET /console/api/enterprise-groups/current

> 返回示例

> 200 Response

```json
{
  "id": "dd074bc2-45b7-4207-ada8-76d74a297b15",
  "name": "Default Group",
  "short_name": null,
  "status": "active",
  "created_at": 1741294643,
  "group_role": "super_admin"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|string|true|none||none|
|» name|string|true|none||none|
|» short_name|null|true|none||none|
|» status|string|true|none||none|
|» created_at|integer|true|none||none|
|» group_role|string|true|none||none|

## GET enterprise-groups current detail

GET /console/api/enterprise-groups/current-detail

> 返回示例

> 200 Response

```json
{
  "enterprise_group": {
    "id": "aeb0aa5f-9010-45bc-a775-41312ca845c9",
    "name": "阿里巴巴",
    "short_name": "阿里巴巴",
    "status": "active",
    "created_at": 1744870218,
    "group_role": "super_admin"
  },
  "shadow_tenant": {
    "id": "aeb0aa5f-9010-45bc-a775-41312ca845c9",
    "name": "阿里巴巴",
    "status": "normal",
    "created_at": 1744870218
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» enterprise_group|object|true|none||none|
|»» id|string|true|none||none|
|»» name|string|true|none||none|
|»» short_name|string|true|none||none|
|»» status|string|true|none||none|
|»» created_at|integer|true|none||none|
|»» group_role|string|true|none||none|
|» shadow_tenant|object|true|none||none|
|»» id|string|true|none||none|
|»» name|string|true|none||none|
|»» status|string|true|none||none|
|»» created_at|integer|true|none||none|

## GET managed-tenants current user

GET /console/api/enterprise-groups/{string:group_id}/managed-tenants

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|string:group_id|path|string| 是 |none|
|page|query|string| 是 |none|
|limit|query|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 20,
  "total": 4,
  "has_more": false,
  "data": [
    {
      "id": "9e9f5925-6ca5-455e-bf8b-e7a413a2280a",
      "name": "my workspace",
      "status": "normal",
      "created_at": 1741317173
    },
    {
      "id": "be5d2163-76cd-40f9-b2ee-c8d0a0cc2f6d",
      "name": "zgi's Workspace",
      "status": "normal",
      "created_at": 1741294639
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|true|none||none|
|»» name|string|true|none||none|
|»» status|string|true|none||none|
|»» created_at|integer|true|none||none|

## GET managed-app-tenants current user

GET /console/api/enterprise-groups/{string:group_id}/managed-app-tenants

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|string:group_id|path|string| 是 |none|
|page|query|string| 是 |none|
|limit|query|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 20,
  "total": 4,
  "has_more": false,
  "data": [
    {
      "id": "9e9f5925-6ca5-455e-bf8b-e7a413a2280a",
      "name": "my workspace",
      "status": "normal",
      "created_at": 1741317173
    },
    {
      "id": "be5d2163-76cd-40f9-b2ee-c8d0a0cc2f6d",
      "name": "zgi's Workspace",
      "status": "normal",
      "created_at": 1741294639
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|true|none||none|
|»» name|string|true|none||none|
|»» status|string|true|none||none|
|»» created_at|integer|true|none||none|

## GET managed-dataset-tenants current user

GET /console/api/enterprise-groups/{string:group_id}/managed-dataset-tenants

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|string:group_id|path|string| 是 |none|
|page|query|string| 是 |none|
|limit|query|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 20,
  "total": 4,
  "has_more": false,
  "data": [
    {
      "id": "9e9f5925-6ca5-455e-bf8b-e7a413a2280a",
      "name": "my workspace",
      "status": "normal",
      "created_at": 1741317173
    },
    {
      "id": "be5d2163-76cd-40f9-b2ee-c8d0a0cc2f6d",
      "name": "zgi's Workspace",
      "status": "normal",
      "created_at": 1741294639
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|true|none||none|
|»» name|string|true|none||none|
|»» status|string|true|none||none|
|»» created_at|integer|true|none||none|

## GET members of enterprise-groups

GET /console/api/enterprise-groups/{string:group_id}/members

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|string:group_id|path|string| 是 |none|
|page|query|string| 是 |none|
|limit|query|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 5,
  "total": 11,
  "has_more": true,
  "data": [
    {
      "id": "a90e5435-0e04-454a-b4d3-fcd5dd61102f",
      "name": "owner_invite_name",
      "avatar": null,
      "avatar_url": null,
      "email": "owner_invite@zgi.ai",
      "is_password_set": false,
      "interface_language": "zh-Hans",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1743465505,
      "status": "pending",
      "account_role": null,
      "extension": {
        "mobile": "13000000020",
        "wechat": null,
        "address": null
      }
    },
    {
      "id": "d9c6d0fd-99fd-4bd7-8b3c-87f29ffb3dd4",
      "name": "super_invite3",
      "avatar": null,
      "avatar_url": null,
      "email": "super_invite3@zgi.ai",
      "is_password_set": false,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1743135164,
      "status": "pending",
      "account_role": {
        "role_type": "system_admin"
      },
      "extension": {
        "mobile": "13000000555",
        "wechat": null,
        "address": null
      }
    },
    {
      "id": "ca95938c-2902-4636-9f30-387a84acef02",
      "name": "super_invite",
      "avatar": null,
      "avatar_url": null,
      "email": "super_invite2@zgi.ai",
      "is_password_set": false,
      "interface_language": "zh-Hans",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1743134685,
      "status": "pending",
      "account_role": {
        "role_type": "system_admin"
      },
      "extension": {
        "mobile": "13000000000",
        "wechat": null,
        "address": null
      }
    },
    {
      "id": "6fef7a4d-b5f5-4350-9db8-31dad9480b72",
      "name": "test_invite4_name",
      "avatar": null,
      "avatar_url": null,
      "email": "test_invite4@zgi.ai",
      "is_password_set": false,
      "interface_language": "zh-Hans",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1742960214,
      "status": "pending",
      "account_role": {
        "role_type": "system_admin"
      },
      "extension": {
        "mobile": "13000000002",
        "wechat": null,
        "address": null
      }
    },
    {
      "id": "04d35c42-96ac-4f28-956d-e6c1929a8953",
      "name": "test_invite2_name",
      "avatar": null,
      "avatar_url": null,
      "email": "test_invite3@zgi.ai",
      "is_password_set": false,
      "interface_language": "zh-Hans",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1742793725,
      "status": "pending",
      "account_role": {
        "role_type": "normal"
      },
      "extension": {
        "mobile": "13000000001",
        "wechat": null,
        "address": null
      }
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|true|none||none|
|»» name|string|true|none||none|
|»» avatar|null|true|none||none|
|»» avatar_url|null|true|none||none|
|»» email|string|true|none||none|
|»» is_password_set|boolean|true|none||none|
|»» interface_language|string|true|none||none|
|»» interface_theme|string|true|none||none|
|»» timezone|string|true|none||none|
|»» last_login_at|null|true|none||none|
|»» last_login_ip|null|true|none||none|
|»» created_at|integer|true|none||none|
|»» status|string|true|none||none|
|»» account_role|object|true|none||none|
|»»» role_type|string|true|none||none|
|»» extension|object|true|none||none|
|»»» mobile|string|true|none||none|
|»»» wechat|null|true|none||none|
|»»» address|null|true|none||none|

## GET datasets of enterprise-groups

GET /console/api/enterprise-groups/{string:group_id}/datasets

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|string:group_id|path|string| 是 |none|
|page|query|string| 是 |none|
|limit|query|string| 是 |none|
|search|query|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 5,
  "total": 1,
  "has_more": false,
  "data": [
    {
      "id": "564ea104-92de-4838-838b-e817d07f42c7",
      "name": "test..",
      "description": "useful for when you want to answer queries about",
      "provider": "vendor",
      "permission": "only_me",
      "data_source_type": "upload_file",
      "indexing_technique": "economy",
      "app_count": 0,
      "document_count": 1,
      "word_count": 0,
      "created_by": "5cfa8de0-f562-409f-bb12-49cbefb7edbf",
      "created_at": 1743730867,
      "updated_by": "5cfa8de0-f562-409f-bb12-49cbefb7edbf",
      "updated_at": 1744278309,
      "embedding_model": "",
      "embedding_model_provider": "",
      "embedding_available": null,
      "retrieval_model_dict": {
        "search_method": "semantic_search",
        "reranking_enable": false,
        "reranking_mode": null,
        "reranking_model": {
          "reranking_provider_name": "",
          "reranking_model_name": ""
        },
        "weights": null,
        "top_k": 2,
        "score_threshold_enabled": false,
        "score_threshold": 1
      },
      "tags": [],
      "doc_form": "text_model",
      "external_knowledge_info": {
        "external_knowledge_id": null,
        "external_knowledge_api_id": null,
        "external_knowledge_api_name": null,
        "external_knowledge_api_endpoint": null
      },
      "external_retrieval_model": {
        "top_k": 2,
        "score_threshold": 1,
        "score_threshold_enabled": false
      },
      "icon": "xxxx",
      "icon_background": "bbbbb",
      "owner": null,
      "owner_account": {
        "name": "zgi"
      },
      "tenant": {
        "name": "zgi's Workspace"
      }
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|false|none||none|
|»» name|string|false|none||none|
|»» description|string|false|none||none|
|»» provider|string|false|none||none|
|»» permission|string|false|none||none|
|»» data_source_type|string|false|none||none|
|»» indexing_technique|string|false|none||none|
|»» app_count|integer|false|none||none|
|»» document_count|integer|false|none||none|
|»» word_count|integer|false|none||none|
|»» created_by|string|false|none||none|
|»» created_at|integer|false|none||none|
|»» updated_by|string|false|none||none|
|»» updated_at|integer|false|none||none|
|»» embedding_model|string|false|none||none|
|»» embedding_model_provider|string|false|none||none|
|»» embedding_available|null|false|none||none|
|»» retrieval_model_dict|object|false|none||none|
|»»» search_method|string|true|none||none|
|»»» reranking_enable|boolean|true|none||none|
|»»» reranking_mode|null|true|none||none|
|»»» reranking_model|object|true|none||none|
|»»»» reranking_provider_name|string|true|none||none|
|»»»» reranking_model_name|string|true|none||none|
|»»» weights|null|true|none||none|
|»»» top_k|integer|true|none||none|
|»»» score_threshold_enabled|boolean|true|none||none|
|»»» score_threshold|integer|true|none||none|
|»» tags|[string]|false|none||none|
|»» doc_form|string|false|none||none|
|»» external_knowledge_info|object|false|none||none|
|»»» external_knowledge_id|null|true|none||none|
|»»» external_knowledge_api_id|null|true|none||none|
|»»» external_knowledge_api_name|null|true|none||none|
|»»» external_knowledge_api_endpoint|null|true|none||none|
|»» external_retrieval_model|object|false|none||none|
|»»» top_k|integer|true|none||none|
|»»» score_threshold|integer|true|none||none|
|»»» score_threshold_enabled|boolean|true|none||none|
|»» icon|string|false|none||none|
|»» icon_background|string|false|none||none|
|»» owner|null|false|none||none|
|»» owner_account|object|false|none||none|
|»»» name|string|true|none||none|
|»» tenant|object|false|none||none|
|»»» name|string|true|none||none|

## POST repair model providers

POST /console/api/enterprise-groups/workspaces/repair-model-providers

> Body 请求参数

```json
{}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET check app permission

GET /console/api/enterprise-groups/{string:group_id}/check-app-permission

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|string:group_id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "has_permission": true
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» has_permission|boolean|true|none||none|

## GET check dataset permission

GET /console/api/enterprise-groups/{string:group_id}/check-dataset-permission

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|string:group_id|path|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "has_permission": true
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» has_permission|boolean|true|none||none|

## POST repair system default

POST /console/api/enterprise-groups/workspaces/repair-system-default

> Body 请求参数

```json
{}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# Account

## POST register

POST /console/api/account/register

> Body 请求参数

```json
{
  "email": "test1@zgi.ai",
  "name": "test1",
  "password": "zgi@zgi123"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» email|body|string| 是 |none|
|» name|body|string| 是 |none|
|» password|body|string| 是 |none|

> 返回示例

```json
{
  "result": "success"
}
```

```json
{
  "code": "user_register_failed",
  "message": "Email already exists."
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|string|true|none||none|
|» name|string|true|none||none|
|» avatar|null|true|none||none|
|» avatar_url|null|true|none||none|
|» email|string|true|none||none|
|» is_password_set|boolean|true|none||none|
|» interface_language|string|true|none||none|
|» interface_theme|string|true|none||none|
|» timezone|string|true|none||none|
|» last_login_at|null|true|none||none|
|» last_login_ip|null|true|none||none|
|» created_at|integer|true|none||none|

## GET profile [Legacy]

GET /console/api/account/profile

> 返回示例

> 200 Response

```json
{
  "id": "4d436d1a-8910-41cb-b208-574726726f26",
  "name": "zgi",
  "avatar": null,
  "avatar_url": null,
  "email": "i@zgi.ai",
  "is_password_set": true,
  "interface_language": "en-US",
  "interface_theme": "light",
  "timezone": "America/New_York",
  "last_login_at": 1741107948,
  "last_login_ip": "127.0.0.1",
  "created_at": 1741083308
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|string|true|none||none|
|» name|string|true|none||none|
|» avatar|null|true|none||none|
|» avatar_url|null|true|none||none|
|» email|string|true|none||none|
|» is_password_set|boolean|true|none||none|
|» interface_language|string|true|none||none|
|» interface_theme|string|true|none||none|
|» timezone|string|true|none||none|
|» last_login_at|integer|true|none||none|
|» last_login_ip|string|true|none||none|
|» created_at|integer|true|none||none|

## GET profile-ex

GET /console/api/account-ex/profile

> 返回示例

> 200 Response

```json
{
  "id": "5cfa8de0-f562-409f-bb12-49cbefb7edbf",
  "name": "zgi",
  "avatar": "a70eaaf8-90b5-4cfc-b848-06d6b3b5792f",
  "avatar_url": "http://127.0.0.1:5001/files/a70eaaf8-90b5-4cfc-b848-06d6b3b5792f/file-preview?timestamp=1743496141&nonce=dec59fce111dfc1b797d3b2f728e6090&sign=sZQqWLLgZi6h6Sm0ptiucRx1rjyoOhMKAH5oa_mADKU=",
  "email": "i@zgi.ai",
  "is_password_set": true,
  "interface_language": "en-US",
  "interface_theme": "light",
  "timezone": "America/New_York",
  "last_login_at": 1743464629,
  "last_login_ip": "127.0.0.1",
  "created_at": 1741294639,
  "status": "active",
  "account_role": {
    "role_type": "super_admin"
  },
  "extension": null
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|string|true|none||none|
|» name|string|true|none||none|
|» avatar|null|true|none||none|
|» avatar_url|null|true|none||none|
|» email|string|true|none||none|
|» is_password_set|boolean|true|none||none|
|» interface_language|string|true|none||none|
|» interface_theme|string|true|none||none|
|» timezone|string|true|none||none|
|» last_login_at|integer|true|none||none|
|» last_login_ip|string|true|none||none|
|» created_at|integer|true|none||none|

## GET account-ex list

GET /console/api/account-ex/list

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|page|query|string| 是 |none|
|limit|query|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "page": 1,
  "limit": 30,
  "total": 21,
  "has_more": false,
  "data": [
    {
      "id": "82932ea0-03f6-4772-b1e5-66cb654553eb",
      "name": "部门管理员测试",
      "avatar": "5b14afcb-6d20-437e-877d-523359a56c2c",
      "avatar_url": "http://127.0.0.1:5001/files/5b14afcb-6d20-437e-877d-523359a56c2c/file-preview?timestamp=1745293621&nonce=b7d3c36d8a22e249f284b34fd54cb3ce&sign=8-nKAQvdrfOWnpaGkke6z1-kODmP5sietibMyq4RDBo=",
      "email": "cuixuezheng@topmh.com",
      "is_password_set": true,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": 1745177120,
      "last_login_ip": "118.247.113.150",
      "created_at": 1744139996,
      "status": "active",
      "group_role": "normal",
      "account_role": null,
      "extension": {
        "mobile": "13401151771",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "19186dbe-8cb0-4306-94ae-26ef00c08195",
      "name": "mj",
      "avatar": null,
      "avatar_url": null,
      "email": "511933643@qq.com",
      "is_password_set": true,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": 1745208257,
      "last_login_ip": "127.0.0.1",
      "created_at": 1741385215,
      "status": "active",
      "group_role": "admin",
      "account_role": null,
      "extension": {
        "mobile": "13000000000",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "39e37e31-d901-498e-a70e-98a8715fe27e",
      "name": "ru",
      "avatar": null,
      "avatar_url": null,
      "email": "15724725214@qq.com",
      "is_password_set": false,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1742944642,
      "status": "pending",
      "group_role": "normal",
      "account_role": null,
      "extension": null
    },
    {
      "id": "c3016769-9c24-4f94-965e-73d18192f327",
      "name": "测试账户",
      "avatar": null,
      "avatar_url": null,
      "email": "544327587@qq.com",
      "is_password_set": false,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1744259977,
      "status": "pending",
      "group_role": "normal",
      "account_role": null,
      "extension": {
        "mobile": "13232211234",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "e8d3d9fa-0081-4bd1-95ed-c3a1b64ac693",
      "name": "jaxson",
      "avatar": null,
      "avatar_url": null,
      "email": "jaxsonnwy123@gmail.com",
      "is_password_set": false,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1743123730,
      "status": "pending",
      "group_role": "normal",
      "account_role": null,
      "extension": {
        "mobile": "18822644410",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "5cfa777e-1be2-4aa1-9f73-31657dad7f74",
      "name": "ruer",
      "avatar": null,
      "avatar_url": null,
      "email": "15724725214@163.com",
      "is_password_set": false,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1742945636,
      "status": "pending",
      "group_role": "normal",
      "account_role": {
        "role_type": "normal"
      },
      "extension": null
    },
    {
      "id": "3ae89aeb-5426-4cc9-9189-889000e8cf42",
      "name": "admin",
      "avatar": null,
      "avatar_url": null,
      "email": "i@zgi.ai2",
      "is_password_set": true,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": 1741157689,
      "last_login_ip": "193.1.245.250",
      "created_at": 1741076754,
      "status": "active",
      "group_role": "normal",
      "account_role": null,
      "extension": null
    },
    {
      "id": "80c40820-0315-4a3b-a1cf-8f4c769e6cd1",
      "name": "ru2",
      "avatar": null,
      "avatar_url": null,
      "email": "690030428@qq.com",
      "is_password_set": false,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1743619403,
      "status": "pending",
      "group_role": "normal",
      "account_role": null,
      "extension": {
        "mobile": "15724725214",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "d10ac372-f817-46a6-98e6-45b63c712a09",
      "name": "飘",
      "avatar": null,
      "avatar_url": null,
      "email": "2945414535@qq.com",
      "is_password_set": false,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1743478445,
      "status": "pending",
      "group_role": "normal",
      "account_role": {
        "role_type": "system_admin"
      },
      "extension": {
        "mobile": "1111",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "bb5af7f9-2ead-402a-8bfe-dec20645a500",
      "name": "Jaxsonwang",
      "avatar": null,
      "avatar_url": null,
      "email": "jaxsonwy123@gmail.com",
      "is_password_set": true,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": 1743465791,
      "last_login_ip": "38.150.67.63",
      "created_at": 1743123819,
      "status": "active",
      "group_role": "normal",
      "account_role": null,
      "extension": {
        "mobile": "188222644410",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "e607ed8f-3c0d-4f14-9e22-1ad3bdcfd0ec",
      "name": "普通用户",
      "avatar": "dd16e716-169e-4d17-a4e1-4dc8d5fde9b3",
      "avatar_url": "http://127.0.0.1:5001/files/dd16e716-169e-4d17-a4e1-4dc8d5fde9b3/file-preview?timestamp=1745293622&nonce=c52e9dc2cf372ebe017f913966b32b89&sign=yYH-lhWwXB14I_EpuGtVPamHnMqBOGTjswsCNxdth18=",
      "email": "491450838@qq.com",
      "is_password_set": true,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": 1745088251,
      "last_login_ip": "219.142.184.58",
      "created_at": 1743557445,
      "status": "active",
      "group_role": "normal",
      "account_role": {
        "role_type": "normal"
      },
      "extension": {
        "mobile": "13444444444",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "f5bba1bc-645d-4e81-b0cb-a2cec4d940b4",
      "name": "jack",
      "avatar": null,
      "avatar_url": null,
      "email": "894532265@qq.com",
      "is_password_set": false,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1743920162,
      "status": "pending",
      "group_role": "normal",
      "account_role": null,
      "extension": {
        "mobile": "18735357273",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "a7bb33a8-a24d-40fb-b647-7ea52bcf298f",
      "name": "chen22",
      "avatar": "fe3217d6-fef1-4ef0-85af-0c5d7839fdf7",
      "avatar_url": "http://127.0.0.1:5001/files/fe3217d6-fef1-4ef0-85af-0c5d7839fdf7/file-preview?timestamp=1745293622&nonce=86ac0d7799973fa53d4c2f6a4db2ab04&sign=tD9QttBGkIqS0xGxI0demedYf07qGivxMjlC2abDsgs=",
      "email": "690030427@qq.com",
      "is_password_set": true,
      "interface_language": "zh-Hans",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": 1744527701,
      "last_login_ip": "115.183.159.105",
      "created_at": 1742767951,
      "status": "active",
      "group_role": "normal",
      "account_role": {
        "role_type": "normal"
      },
      "extension": {
        "mobile": "15724725218",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "710026aa-fe43-457c-82d6-3ce29f7f9eb3",
      "name": "jack",
      "avatar": null,
      "avatar_url": null,
      "email": "8945322265@qq.com",
      "is_password_set": false,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": null,
      "last_login_ip": null,
      "created_at": 1744314534,
      "status": "pending",
      "group_role": "normal",
      "account_role": null,
      "extension": {
        "mobile": "18735357273",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "14b20176-2811-42b2-a9dd-d345577c0de6",
      "name": "飘丨雪",
      "avatar": "38d2986d-9148-4ab5-bd0b-ff47a3842044",
      "avatar_url": "http://127.0.0.1:5001/files/38d2986d-9148-4ab5-bd0b-ff47a3842044/file-preview?timestamp=1745293622&nonce=e58d786b9b3144588ba85eb8dbe70097&sign=sQpKuFJuzxuxhwRxJsK6mkl6y2IJ0g_z1NT2CJFg59s=",
      "email": "zaixiapiaoxue@foxmail.com",
      "is_password_set": true,
      "interface_language": "zh-Hans",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": 1743657073,
      "last_login_ip": "117.154.107.239",
      "created_at": 1741730266,
      "status": "active",
      "group_role": "normal",
      "account_role": {
        "role_type": "normal"
      },
      "extension": {
        "mobile": "18271466597",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "c48069cf-c116-4412-a04a-75279ee52eb8",
      "name": "491450838",
      "avatar": null,
      "avatar_url": null,
      "email": " 491450838@qq.com",
      "is_password_set": true,
      "interface_language": "zh-Hans",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": 1743401163,
      "last_login_ip": "127.0.0.1",
      "created_at": 1741412562,
      "status": "banned",
      "group_role": "normal",
      "account_role": {
        "role_type": "normal"
      },
      "extension": null
    },
    {
      "id": "d22b6ea2-fbf9-44fe-88dd-8b7a6c1b1be2",
      "name": "mj",
      "avatar": null,
      "avatar_url": null,
      "email": "mj@test.com",
      "is_password_set": true,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": 1744155337,
      "last_login_ip": "113.224.158.124",
      "created_at": 1744155064,
      "status": "active",
      "group_role": "normal",
      "account_role": null,
      "extension": {
        "mobile": "13000000000",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "13539725-99cf-4ad7-b8bd-3cc0569b3d52",
      "name": "chen",
      "avatar": null,
      "avatar_url": null,
      "email": "1572725214@163.com",
      "is_password_set": true,
      "interface_language": "zh-Hans",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": 1744323241,
      "last_login_ip": "218.30.116.113",
      "created_at": 1744314697,
      "status": "active",
      "group_role": "normal",
      "account_role": null,
      "extension": {
        "mobile": "15724678732",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "ac8e00cf-fa55-4721-b664-95b547bd865f",
      "name": "test",
      "avatar": null,
      "avatar_url": null,
      "email": "imageidol@126.com",
      "is_password_set": true,
      "interface_language": "zh-Hans",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": 1745037882,
      "last_login_ip": "188.253.115.252",
      "created_at": 1743478342,
      "status": "active",
      "group_role": "normal",
      "account_role": {
        "role_type": "system_admin"
      },
      "extension": {
        "mobile": "13274567890",
        "wechat": null,
        "address": null,
        "gender": null
      }
    },
    {
      "id": "637c1fd0-3c78-473f-b3eb-2080bccaf7c7",
      "name": "ruru",
      "avatar": null,
      "avatar_url": null,
      "email": "879661217@qq.com",
      "is_password_set": true,
      "interface_language": "en-US",
      "interface_theme": "light",
      "timezone": "America/New_York",
      "last_login_at": 1744955637,
      "last_login_ip": "123.114.97.152",
      "created_at": 1744955635,
      "status": "active",
      "group_role": "normal",
      "account_role": null,
      "extension": null
    },
    {
      "id": "2f729148-28b3-47fb-8398-5fcef1d3f4a8",
      "name": "宁静夏天",
      "avatar": "81de87ec-dfa4-4d14-ba02-aba2644683ae",
      "avatar_url": "http://127.0.0.1:5001/files/81de87ec-dfa4-4d14-ba02-aba2644683ae/file-preview?timestamp=1745293623&nonce=1c64ad1f6779009f34f2407e5bc90535&sign=8MFzSIcAeVX9VZMeMCdh5yAIaBcP-UDxczEO2V89C_8=",
      "email": "i@zgi.ai",
      "is_password_set": true,
      "interface_language": "zh-Hans",
      "interface_theme": "light",
      "timezone": "Asia/Shanghai",
      "last_login_at": 1745264048,
      "last_login_ip": "127.0.0.1",
      "created_at": 1741076754,
      "status": "active",
      "group_role": "admin",
      "account_role": {
        "role_type": "super_admin"
      },
      "extension": {
        "mobile": "13828372384",
        "wechat": null,
        "address": null,
        "gender": null
      }
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» page|integer|true|none||none|
|» limit|integer|true|none||none|
|» total|integer|true|none||none|
|» has_more|boolean|true|none||none|
|» data|[object]|true|none||none|
|»» id|string|true|none||none|
|»» name|string|true|none||none|
|»» avatar|string¦null|true|none||none|
|»» avatar_url|string¦null|true|none||none|
|»» email|string|true|none||none|
|»» is_password_set|boolean|true|none||none|
|»» interface_language|string|true|none||none|
|»» interface_theme|string|true|none||none|
|»» timezone|string|true|none||none|
|»» last_login_at|integer¦null|true|none||none|
|»» last_login_ip|string¦null|true|none||none|
|»» created_at|integer|true|none||none|
|»» status|string|true|none||none|
|»» group_role|string|true|none||none|
|»» account_role|object|true|none||none|
|»»» role_type|string|true|none||none|
|»» extension|object|true|none||none|
|»»» mobile|string|true|none||none|
|»»» wechat|null|true|none||none|
|»»» address|null|true|none||none|
|»»» gender|null|true|none||none|

## GET account-ex by email

GET /console/api/account-ex/email

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|email|query|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "id": "f65d9e29-06bf-4fbd-8e6a-c779c9f3ad17",
  "name": "test_invite2_name",
  "avatar": null,
  "avatar_url": null,
  "email": "test_invite2@zgi.ai",
  "is_password_set": false,
  "interface_language": "zh-Hans",
  "interface_theme": "light",
  "timezone": "Asia/Shanghai",
  "last_login_at": null,
  "last_login_ip": null,
  "created_at": 1741189335,
  "extension": {
    "mobile": "13000000000",
    "wechat": null,
    "address": null
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|string|true|none||none|
|» name|string|true|none||none|
|» avatar|null|true|none||none|
|» avatar_url|null|true|none||none|
|» email|string|true|none||none|
|» is_password_set|boolean|true|none||none|
|» interface_language|string|true|none||none|
|» interface_theme|string|true|none||none|
|» timezone|string|true|none||none|
|» last_login_at|null|true|none||none|
|» last_login_ip|null|true|none||none|
|» created_at|integer|true|none||none|
|» extension|object|true|none||none|
|»» mobile|string|true|none||none|
|»» wechat|null|true|none||none|
|»» address|null|true|none||none|

## GET account-ex by id

GET /console/api/account-ex/id

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|query|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "id": "f65d9e29-06bf-4fbd-8e6a-c779c9f3ad17",
  "name": "test_invite2_name",
  "avatar": null,
  "avatar_url": null,
  "email": "test_invite2@zgi.ai",
  "is_password_set": false,
  "interface_language": "zh-Hans",
  "interface_theme": "light",
  "timezone": "Asia/Shanghai",
  "last_login_at": null,
  "last_login_ip": null,
  "created_at": 1741189335,
  "extension": {
    "mobile": "13000000000",
    "wechat": null,
    "address": null
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|string|true|none||none|
|» name|string|true|none||none|
|» avatar|null|true|none||none|
|» avatar_url|null|true|none||none|
|» email|string|true|none||none|
|» is_password_set|boolean|true|none||none|
|» interface_language|string|true|none||none|
|» interface_theme|string|true|none||none|
|» timezone|string|true|none||none|
|» last_login_at|null|true|none||none|
|» last_login_ip|null|true|none||none|
|» created_at|integer|true|none||none|
|» extension|object|true|none||none|
|»» mobile|string|true|none||none|
|»» wechat|null|true|none||none|
|»» address|null|true|none||none|

## PUT account-ex update

PUT /console/api/account-ex/id

> Body 请求参数

```json
{
  "name": "test User",
  "avatar": "https://example.com/avatar.jpg",
  "status": "active",
  "mobile": "13800138000",
  "role_type": "system_admin",
  "group_role": "admin"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|query|string| 否 |none|
|body|body|object| 否 |none|
|» name|body|string| 否 |none|
|» avatar|body|string| 否 |none|
|» status|body|string| 否 |none|
|» mobile|body|string| 否 |none|
|» role_type|body|string| 否 |none|
|» group_role|body|[enterprise group role enum](#schemaenterprise group role enum)| 否 |none|

#### 枚举值

|属性|值|
|---|---|
|» group_role|admin|
|» group_role|normal|

> 返回示例

> 200 Response

```json
{
  "id": "f65d9e29-06bf-4fbd-8e6a-c779c9f3ad17",
  "name": "test_invite2_name",
  "avatar": null,
  "avatar_url": null,
  "email": "test_invite2@zgi.ai",
  "is_password_set": false,
  "interface_language": "zh-Hans",
  "interface_theme": "light",
  "timezone": "Asia/Shanghai",
  "last_login_at": null,
  "last_login_ip": null,
  "created_at": 1741189335,
  "extension": {
    "mobile": "13000000000",
    "wechat": null,
    "address": null
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|string|true|none||none|
|» name|string|true|none||none|
|» avatar|null|true|none||none|
|» avatar_url|null|true|none||none|
|» email|string|true|none||none|
|» is_password_set|boolean|true|none||none|
|» interface_language|string|true|none||none|
|» interface_theme|string|true|none||none|
|» timezone|string|true|none||none|
|» last_login_at|null|true|none||none|
|» last_login_ip|null|true|none||none|
|» created_at|integer|true|none||none|
|» extension|object|true|none||none|
|»» mobile|string|true|none||none|
|»» wechat|null|true|none||none|
|»» address|null|true|none||none|

## PUT account-ex update role

PUT /console/api/account-ex/role

> Body 请求参数

```json
{
  "role_type": "system_admin"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|query|string| 否 |none|
|body|body|object| 否 |none|
|» role_type|body|string| 是 |none|

#### 枚举值

|属性|值|
|---|---|
|» role_type|system_admin|
|» role_type|normal|

> 返回示例

> 200 Response

```json
{
  "id": "87b23068-0a1c-493e-bf6c-e910ab34d74a",
  "name": "testme",
  "avatar": null,
  "avatar_url": null,
  "email": "511933643@qq.com",
  "is_password_set": true,
  "interface_language": "zh-Hans",
  "interface_theme": "light",
  "timezone": "Asia/Shanghai",
  "last_login_at": 1741396672,
  "last_login_ip": "127.0.0.1",
  "created_at": 1741300607,
  "status": "active",
  "account_role": {
    "role_type": "system_admin"
  },
  "extension": {
    "mobile": "13000000000",
    "wechat": null,
    "address": null
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» id|string|true|none||none|
|» name|string|true|none||none|
|» avatar|null|true|none||none|
|» avatar_url|null|true|none||none|
|» email|string|true|none||none|
|» is_password_set|boolean|true|none||none|
|» interface_language|string|true|none||none|
|» interface_theme|string|true|none||none|
|» timezone|string|true|none||none|
|» last_login_at|integer|true|none||none|
|» last_login_ip|string|true|none||none|
|» created_at|integer|true|none||none|
|» status|string|true|none||none|
|» account_role|object|true|none||none|
|»» role_type|string|true|none||none|
|» extension|object|true|none||none|
|»» mobile|string|true|none||none|
|»» wechat|null|true|none||none|
|»» address|null|true|none||none|

# Auth

## POST login [Legacy]

POST /console/api/login

> Body 请求参数

```json
{
  "email": "i@zgi.ai",
  "name": "zgi",
  "password": "zgi@zgi123"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» email|body|string| 是 |none|
|» name|body|string| 是 |none|
|» password|body|string| 是 |none|
|» invite_token|body|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "result": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGQ0MzZkMWEtODkxMC00MWNiLWIyMDgtNTc0NzI2NzI2ZjI2IiwiZXhwIjoxNzQxMDg3MjM5LCJpc3MiOiJTRUxGX0hPU1RFRCIsInN1YiI6IkNvbnNvbGUgQVBJIFBhc3Nwb3J0In0.M7kRERwJogXqpr0EnEbKdc0uMv-LiJ7sqdnfMaCzyAw",
    "refresh_token": "5c3affea7d2e2e28c8712b7f7317f6c7f240a44ec31811e16452b64bb3df595d7e9e3003f8d613d836573fd477751cfa21f477833eadc8a611921133bb2da52b"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» result|string|true|none||none|
|» data|object|true|none||none|
|»» access_token|string|true|none||none|
|»» refresh_token|string|true|none||none|

## GET activate check [Legacy]

GET /console/api/activate/check

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|workspace_id|query|string| 否 |none|
|email|query|string| 是 |none|
|token|query|string| 是 |none|

> 返回示例

```json
{
  "is_valid": false
}
```

```json
{
  "is_valid": true,
  "data": {
    "workspace_name": "zgi's Workspace",
    "workspace_id": "865addd7-f137-4c04-8d54-f038ba372d90",
    "email": "test2@zgi.ai"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST activate [Legacy]

POST /console/api/activate

> Body 请求参数

```json
{
  "token": "d3733ec3-3fa0-42e8-81ec-61706674c153",
  "name": "test3",
  "interface_language": "en-US",
  "timezone": "Asia/Shanghai"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» token|body|string| 是 |none|
|» name|body|string| 是 |none|
|» interface_language|body|string| 是 |none|
|» timezone|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## POST forgot-password [Legacy]

POST /console/api/forgot-password

> Body 请求参数

```json
{
  "email": "test5@zgi.ai",
  "language": "en-US"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» email|body|string| 是 |none|
|» language|body|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "result": "success",
  "data": "af3148c6-632e-4e8d-adeb-ec28e9290758"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» result|string|true|none||none|
|» data|string|true|none||none|

## POST forgot-password validity [Legacy]

POST /console/api/forgot-password/validity

> Body 请求参数

```json
{
  "email": "test5@zgi.ai",
  "code": "793775",
  "token": "f5a68dd1-8df6-4602-ba45-68fe9459277f"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» email|body|string| 是 |none|
|» code|body|string| 是 |none|
|» token|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "is_valid": true,
  "email": "test5@zgi.ai"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» is_valid|boolean|true|none||none|
|» email|string|true|none||none|

## POST forgot-password resets [Legacy]

POST /console/api/forgot-password/resets

> Body 请求参数

```json
{
  "email": "i@zgi.ai",
  "name": "zgi",
  "password": "zgi@zgi123"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» email|body|string| 是 |none|
|» password|body|string| 是 |none|
|» password_confirm|body|string| 是 |none|
|» token|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "result": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGQ0MzZkMWEtODkxMC00MWNiLWIyMDgtNTc0NzI2NzI2ZjI2IiwiZXhwIjoxNzQxMDg3MjM5LCJpc3MiOiJTRUxGX0hPU1RFRCIsInN1YiI6IkNvbnNvbGUgQVBJIFBhc3Nwb3J0In0.M7kRERwJogXqpr0EnEbKdc0uMv-LiJ7sqdnfMaCzyAw",
    "refresh_token": "5c3affea7d2e2e28c8712b7f7317f6c7f240a44ec31811e16452b64bb3df595d7e9e3003f8d613d836573fd477751cfa21f477833eadc8a611921133bb2da52b"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» result|string|true|none||none|
|» data|object|true|none||none|
|»» access_token|string|true|none||none|
|»» refresh_token|string|true|none||none|

## POST register

POST /console/api/register

> Body 请求参数

```json
{
  "email": "test5@zgi.ai",
  "language": "en-US"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» email|body|string| 是 |none|
|» language|body|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "result": "success",
  "data": "af3148c6-632e-4e8d-adeb-ec28e9290758"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» result|string|true|none||none|
|» data|string|true|none||none|

## POST register validity

POST /console/api/register/validity

> Body 请求参数

```json
{
  "email": "test5@zgi.ai",
  "code": "793775",
  "token": "f5a68dd1-8df6-4602-ba45-68fe9459277f"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» email|body|string| 是 |none|
|» code|body|string| 是 |none|
|» token|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "is_valid": true,
  "email": "test5@zgi.ai"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» is_valid|boolean|true|none||none|
|» email|string|true|none||none|

## POST register finish

POST /console/api/register/finish

> Body 请求参数

```json
{
  "email": "register5@zgi.ai",
  "name": "register5",
  "password": "zgi@zgi123",
  "password_confirm": "zgi@zgi123",
  "token": "bd8a8d9e-c6e1-4d2d-a062-85c74e5330c8"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» email|body|string| 是 |none|
|» name|body|string| 否 |none|
|» password|body|string| 是 |none|
|» password_confirm|body|string| 是 |none|
|» token|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "result": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMWVhNTk3NDUtZjYwOC00YWExLTgyZTMtZGZlOTE5NjY4ZTAzIiwiZXhwIjoxNzQ0OTYwNzI4LCJpc3MiOiJTRUxGX0hPU1RFRCIsInN1YiI6IkNvbnNvbGUgQVBJIFBhc3Nwb3J0In0.R-ao5QTRxZJMAeE4b_CJzLRzjnROuNl-u2ZUj3r4BPw",
    "refresh_token": "1731b2a37fdbc64df69924bbfe9c69f16f04605110544a02d4465cd13ba71114cb980d3e052971c786c84098a52524b8867c37f0761f24059d6eea4caf1c5794"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» result|string|true|none||none|
|» data|object|true|none||none|
|»» access_token|string|true|none||none|
|»» refresh_token|string|true|none||none|

## GET oauth login [Legacy]

GET /console/api/oauth/login/{provider}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|provider|path|string| 是 |none|

> 返回示例

> 200 Response

```json
"string"
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|string|

## GET oauth authorize [Legacy]

GET /console/api/oauth/authorize/{provider}

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|provider|path|string| 是 |none|
|code|query|string| 是 |none|
|state|query|string| 是 |none|

> 返回示例

> 200 Response

```json
"string"
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|string|

# 数据模型

<h2 id="tocS_enterprise group role enum">enterprise group role enum</h2>

<a id="schemaenterprise group role enum"></a>
<a id="schema_enterprise group role enum"></a>
<a id="tocSenterprise group role enum"></a>
<a id="tocsenterprise group role enum"></a>

```json
"admin"

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|*anonymous*|string|false|none||none|

#### 枚举值

|属性|值|
|---|---|
|*anonymous*|admin|
|*anonymous*|normal|

