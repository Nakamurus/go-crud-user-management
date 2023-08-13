# Go CRUD User Management with Gin and Docker

Go言語、Ginを使ったユーザー管理CRUDアプリケーションです。Dockerでコンテナ化して、PostgreSQLとRedisを使ったシンプルなシステムです。
[For English, click here 英語はこちら](./README_ja.md)

## 特徴

- **Golang**
- **Gin**
- **Docker**
- **Security**: JWTベースの認証システムで、ブラックリストシステムのためRedisを使っています。
- **CI/CD Test**: Github ActionsのCI/CDテストでプッシュのたびにテストを走らせています。

## API エンドポイントとサンプルリクエスト

> `uuid` は実際のUUIDで置き換えてください。

## Hello World

- **ルート:** GET /
- **レスポンス:**

  - **200 OK:**
  - **Hello World**

## Authentication ルートs

### 1. Login

- **ルート:** POST /login
- **リクエストボディ:**

```json
{
    "email": "example@email.com",
    "password": "password123"
}
```

- **レスポンス:**

  - **200 OK:**

  ```json
  {
      "token": "your_jwt_token_here"
  }
  ```

  - 400 Bad Request: { "error": "error_message_here" }
  - 401 Unauthorized: { "error": "Invalid email or password" }
  - 500 Internal Server Error: { "error": "Error retrieving user" }

### 2. Change Password

- **ルート:** POST /me/uuid/password
- **リクエストボディ:**

```json
{
    "old_password": "old_password_here",
    "new_password": "new_password_here"
}
```

- **リクエストヘッダ**

```bash
Authorization: Bearer your_jwt_token_here
```

- **レスポンス:**

  - **200 OK:**

  ```json
  {
      "message": "Password updated successfully"
  }
  ```

  - Other Status Codes: { "error": "error_message_here" }

### 3. Refresh Token

- **ルート:** GET /me/refresh-token
- **リクエストヘッダ**

```bash
Authorization: Bearer your_jwt_token_here
```

- **レスポンス:**

  - **200 OK:**

  ```json
  {
      "token": "new_jwt_token_here"
  }
  ```

  - Other Status Codes: { "error": "error_message_here" }

### 4. Logout

- **ルート:** GET /me/logout
- **リクエストヘッダ**

```bash
Authorization: Bearer your_jwt_token_here
```

- **レスポンス:**

  - 200 OK: { "message": "Successfully logged out" }
  - Other Status Codes: { "error": "error_message_here" }

## User ルートs

### 1. List Users

- **ルート:** GET /users
- **レスポンス:**

  - 200 OK: Array of user objects.
  - 500 Internal Server Error: { "error": "error_message_here" }

### 2. Get a User

- **ルート:** GET /user/:uuid
- **レスポンス:**

  - 200 OK: User object.
  - 404 Not Found: { "error": "User not found" }
  - 500 Internal Server Error: { "error": "error_message_here" }

### 3. Create a User

- **ルート:** POST /user
- **リクエストボディ:**

```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
}
```

- **レスポンス:**

  - 200 OK: { "message": "user created John Doe" }
  - 400 Bad Request: { "error": "Missing required fields" }
  - 500 Internal Server Error: { "error": "error_message_here" }

### 4. Update a User

- **ルート:** PUT /me/:uuid
- **リクエストボディ (Partial updates allowed):**

```json
{
    "name": "New Name",
    "email": "newemail@example.com"
}
```

- **リクエストヘッダ**

```bash
Authorization: Bearer your_jwt_token_here
```

- **レスポンス:**

  - **200 OK:**

  ```json
  {
      "token": "new_jwt_token_here",
      "user": updated_user_object,
      "message": "User updated successfully"
  }
  ```

  - Other Status Codes: { "error": "error_message_here" }

### 5. Delete a User

- **ルート:** DELETE /me/:uuid
- **リクエストヘッダ**

```bash
Authorization: Bearer your_jwt_token_here
```

- **レスポンス:**

  - 200 OK: { "message": "user deleted" }
  - 500 Internal Server Error: { "error": "error_message_here" }
