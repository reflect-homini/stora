# API Contract

## Overview

- **API Purpose**: Authentication and User Management service.
- **Base URL**: `/api/v1`
- **Authentication mechanism**: JWT (JSON Web Token) using Bearer scheme.
- **Response conventions**: All responses are JSON.
- **Error conventions**: Standard error format provided by `ungerr` library.

## Authentication

Authentication is handled via JWT. Most endpoints require the `Authorization` header.

Header: `Authorization: Bearer <token>`

Example:

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Global Conventions

- **Content-Type**: `application/json`
- **Response Structure**: All responses follow a standardized wrapper.
- **Pagination**: Included in the response wrapper when applicable.

### Standard Response Wrapper

```json
{
  "data": { ... },
  "errors": [
    {
      "message": "Error message",
      "code": 400
    }
  ],
  "pagination": {
    "totalData": 120,
    "currentPage": 1,
    "totalPages": 6,
    "hasNextPage": true,
    "hasPrevPage": false
  }
}
```

Fields:

- `data`: The actual payload (omitted if null/empty).
- `errors`: An array of error objects (omitted if empty).
- `pagination`: Metadata for paginated results (omitted if zero).

## Endpoint Documentation

### Register User

#### Request

**Method**: `POST`  
**Path**: `/api/v1/auth/register`  
**Headers**: `Content-Type: application/json`

**Body**:

| Field                | Type   | Required | Description                       |
| -------------------- | ------ | -------- | --------------------------------- |
| email                | string | yes      | User email address (min length 3) |
| password             | string | yes      | User password                     |
| passwordConfirmation | string | yes      | Must match password               |

Example:

```json
{
  "email": "user@example.com",
  "password": "securepassword",
  "passwordConfirmation": "securepassword"
}
```

#### Response

**Status**: `201 Created`

```json
{
  "data": {
    "message": "registration successful"
  }
}
```

---

### Internal Login

#### Request

**Method**: `POST`  
**Path**: `/api/v1/auth/login`  
**Headers**: `Content-Type: application/json`

**Body**:

| Field    | Type   | Required | Description        |
| -------- | ------ | -------- | ------------------ |
| email    | string | yes      | User email address |
| password | string | yes      | User password      |

Example:

```json
{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Internal Login Response

**Status**: `200 OK`

```json
{
  "data": {
    "type": "Bearer",
    "token": "access_token_here",
    "refreshToken": "refresh_token_here"
  }
}
```

---

### Refresh Token

#### Request

**Method**: `PUT`  
**Path**: `/api/v1/auth/refresh`  
**Headers**: `Content-Type: application/json`

**Body**:

| Field        | Type   | Required | Description         |
| ------------ | ------ | -------- | ------------------- |
| refreshToken | string | yes      | Valid refresh token |

#### Response

**Status**: `200 OK`

```json
{
  "data": {
    "type": "Bearer",
    "token": "new_access_token_here",
    "refreshToken": "new_refresh_token_here"
  }
}
```

---

### OAuth2 Login Redirect

#### Request

**Method**: `GET`  
**Path**: `/api/v1/auth/:provider`  
**Description**: Redirects to OAuth2 provider's login page.

---

### OAuth2 Callback

#### Request

**Method**: `GET`  
**Path**: `/api/v1/auth/:provider/callback`  
**Query Params**:

- `code`: Authorization code from provider
- `state`: State parameter for CSRF protection

#### OAuth2 Callback Response

**Status**: `200 OK`

```json
{
  "data": {
    "type": "Bearer",
    "token": "access_token_here",
    "refreshToken": "refresh_token_here"
  }
}
```

---

### Verify Registration

#### Request

**Method**: `GET`  
**Path**: `/api/v1/auth/verify-registration`  
**Query Params**:

- `token`: Verification token sent via email

#### Response

**Status**: `200 OK`

---

### Send Password Reset

#### Request

**Method**: `POST`  
**Path**: `/api/v1/auth/password-reset`  
**Headers**: `Content-Type: application/json`

**Body**:

| Field | Type   | Required | Description        |
| ----- | ------ | -------- | ------------------ |
| email | string | yes      | User email address |

#### Response

**Status**: `201 Created`

---

### Reset Password

#### Request

**Method**: `PATCH`  
**Path**: `/api/v1/auth/reset-password`  
**Headers**: `Content-Type: application/json`

**Body**:

| Field                | Type   | Required | Description            |
| -------------------- | ------ | -------- | ---------------------- |
| token                | string | yes      | Reset token from email |
| password             | string | yes      | New password           |
| passwordConfirmation | string | yes      | Must match password    |

#### Response

**Status**: `200 OK`

---

### Logout

#### Request

**Method**: `DELETE`  
**Path**: `/api/v1/auth/logout`  
**Headers**: `Authorization: Bearer <token>`

#### Response

**Status**: `204 No Content`

---

### Get Current User Profile (Me)

#### Request

**Method**: `GET`  
**Path**: `/api/v1/me`  
**Headers**: `Authorization: Bearer <token>`

#### Response

**Status**: `200 OK`

```json
{
  "data": {
    "id": "uuid-here",
    "createdAt": "2026-03-10T00:00:00Z",
    "updatedAt": "2026-03-10T00:00:00Z",
    "email": "user@example.com",
    "profile": {
      "id": "uuid-here",
      "createdAt": "2026-03-10T00:00:00Z",
      "updatedAt": "2026-03-10T00:00:00Z",
      "userId": "uuid-here",
      "name": "User Name",
      "avatar": "url_to_avatar"
    }
  }
}
```

---

### Create Project

#### Request

**Method**: `POST`
**Path**: `/api/v1/projects`
**Headers**:

- `Authorization: Bearer <token>`
- `Content-Type: application/json`

**Body**:

| Field       | Type   | Required | Description                 |
| ----------- | ------ | -------- | --------------------------- |
| name        | string | yes      | Project name (min length 3) |
| description | string | no       | Project description         |

Example:

```json
{
  "name": "My New Project",
  "description": "A description of the project"
}
```

#### Response

**Status**: `201 Created`

```json
{
  "data": {
    "id": "uuid-here",
    "createdAt": "2026-03-11T00:00:00Z",
    "updatedAt": "2026-03-11T00:00:00Z",
    "userId": "uuid-here",
    "name": "My New Project",
    "description": "A description of the project",
    "lastInteractedAt": "2026-03-11T00:00:00Z",
    "entries": []
  }
}
```

---

### Get All Projects

#### Request

**Method**: `GET`  
**Path**: `/api/v1/projects`  
**Headers**: `Authorization: Bearer <token>`

#### Response

**Status**: `200 OK`

```json
{
  "data": [
    {
      "id": "uuid-here",
      "createdAt": "2026-03-11T00:00:00Z",
      "updatedAt": "2026-03-11T00:00:00Z",
      "userId": "uuid-here",
      "name": "My New Project",
      "description": "A description of the project",
      "lastInteractedAt": "2026-03-11T00:00:00Z",
      "entries": []
    }
  ]
}
```

---

### Get Project By ID

#### Request

**Method**: `GET`  
**Path**: `/api/v1/projects/:projectID`  
**Headers**: `Authorization: Bearer <token>`

#### Response

**Status**: `200 OK`

```json
{
  "data": {
    "id": "uuid-here",
    "createdAt": "2026-03-11T00:00:00Z",
    "updatedAt": "2026-03-11T00:00:00Z",
    "userId": "uuid-here",
    "name": "My New Project",
    "description": "A description of the project",
    "lastInteractedAt": "2026-03-11T00:00:00Z",
    "entries": [
      {
        "id": "uuid-here",
        "createdAt": "2026-03-11T00:00:00Z",
        "updatedAt": "2026-03-11T00:00:00Z",
        "projectId": "uuid-here",
        "content": "Entry content here"
      }
    ]
  }
}
```

---

### Create Entry

#### Request

**Method**: `POST`
**Path**: `/api/v1/projects/:projectID/entries`
**Headers**:

- `Authorization: Bearer <token>`
- `Content-Type: application/json`

**Body**:

| Field   | Type   | Required | Description                  |
| ------- | ------ | -------- | ---------------------------- |
| content | string | yes      | Entry content (min length 3) |

Example:

```json
{
  "content": "A new entry for the project"
}
```

#### Response

**Status**: `201 Created`

```json
{
  "data": {
    "id": "uuid-here",
    "createdAt": "2026-03-11T00:00:00Z",
    "updatedAt": "2026-03-11T00:00:00Z",
    "projectId": "uuid-here",
    "content": "A new entry for the project"
  }
}
```

---

### Update Entry

#### Request

**Method**: `PUT`
**Path**: `/api/v1/projects/:projectID/entries/:entryID`
**Headers**:

- `Authorization: Bearer <token>`
- `Content-Type: application/json`

**Body**:

| Field   | Type   | Required | Description                  |
| ------- | ------ | -------- | ---------------------------- |
| content | string | yes      | Entry content (min length 3) |

Example:

```json
{
  "content": "Updated entry content"
}
```

#### Response

**Status**: `200 OK`

```json
{
  "data": {
    "id": "uuid-here",
    "createdAt": "2026-03-11T00:00:00Z",
    "updatedAt": "2026-03-11T00:00:00Z",
    "projectId": "uuid-here",
    "content": "Updated entry content"
  }
}
```

---

### Delete Entry

#### Request

**Method**: `DELETE`
**Path**: `/api/v1/projects/:projectID/entries/:entryID`
**Headers**:

- `Authorization: Bearer <token>`

#### Response

**Status**: `204 No Content`

## Data Models

### BaseDTO

| Field     | Type             | Description           |
| --------- | ---------------- | --------------------- |
| id        | string (UUID)    | Unique identifier     |
| createdAt | string (ISO8601) | Creation timestamp    |
| updatedAt | string (ISO8601) | Last update timestamp |

### UserResponse

| Field   | Type            | Description          |
| ------- | --------------- | -------------------- |
| id      | string (UUID)   | Unique user ID       |
| email   | string          | User email address   |
| profile | ProfileResponse | User profile details |

### ProfileResponse

| Field  | Type          | Description        |
| ------ | ------------- | ------------------ |
| id     | string (UUID) | Unique profile ID  |
| userId | string (UUID) | Associated User ID |
| name   | string        | User display name  |
| avatar | string        | Avatar URL         |

### ProjectResponse

| Field            | Type             | Description                    |
| ---------------- | ---------------- | ------------------------------ |
| id               | string (UUID)    | Unique project ID              |
| createdAt        | string (ISO8601) | Creation timestamp             |
| updatedAt        | string (ISO8601) | Last update timestamp          |
| userId           | string (UUID)    | Owner's User ID                |
| name             | string           | Project name                   |
| description      | string           | Project description            |
| lastInteractedAt | string (ISO8601) | Last interaction timestamp     |
| entries          | EntryResponse[]  | List of entries in the project |

### EntryResponse

| Field     | Type             | Description           |
| --------- | ---------------- | --------------------- |
| id        | string (UUID)    | Unique entry ID       |
| createdAt | string (ISO8601) | Creation timestamp    |
| updatedAt | string (ISO8601) | Last update timestamp |
| projectId | string (UUID)    | Associated Project ID |
| content   | string           | Entry content         |

### Pagination

| Field       | Type    | Description                      |
| ----------- | ------- | -------------------------------- |
| totalData   | integer | Total number of items            |
| currentPage | integer | Current page number              |
| totalPages  | integer | Total number of pages            |
| hasNextPage | boolean | True if there is a next page     |
| hasPrevPage | boolean | True if there is a previous page |
