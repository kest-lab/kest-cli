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

