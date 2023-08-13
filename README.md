# Go CRUD User Management: Simple and Modern

Welcome to the `go-crud-user-management` project. Built with **Golang**, **Gin**, and containerized with **Docker**, this is a toy application designed to demonstrate a basic user management system.

## Features

- **Golang**: Using Go provides simplicity and clarity, making the codebase easy to understand.
- **Web Framework**: We're utilizing Gin, a framework that offers a lean, efficient web server with easy JSON rendering.
- **Containerization**: Thanks to Docker, you can ensure that the app runs the same regardless of where it's deployed.
- **Security**: Simple JWT-based authentication complemented with a Redis-backed blacklisting system.
- **Structured Code**: The code is modular, making it easier for newcomers to follow and learn from.
- **CI/CD Test**: CI/CD Test ensures stability and reliability.

## API Endpoints & Sample Requests

> Note: For these sample requests, replace `uuid` with your actual UUID.

### Public Routes

# User Management API (Built with Golang, Gin, and Docker)


## Hello World

- **Route:** GET /
- **Responses:**

  - **200 OK:**
  - **Hello World**

## Authentication Routes

### 1. Login

- **Route:** POST /login
- **Request Body:**

```json
{
    "email": "example@email.com",
    "password": "password123"
}
```

- **Responses:**

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

- **Route:** POST /me/uuid/password
- **Request Body:**

```json
{
    "old_password": "old_password_here",
    "new_password": "new_password_here"
}
```

- **Request Header**

```bash
Authorization: Bearer your_jwt_token_here
```

- **Responses:**

  - **200 OK:**

  ```json
  {
      "message": "Password updated successfully"
  }
  ```

  - Other Status Codes: { "error": "error_message_here" }

### 3. Refresh Token

- **Route:** GET /me/refresh-token
- **Request Header**

```bash
Authorization: Bearer your_jwt_token_here
```

- **Responses:**

  - **200 OK:**

  ```json
  {
      "token": "new_jwt_token_here"
  }
  ```

  - Other Status Codes: { "error": "error_message_here" }

### 4. Logout

- **Route:** GET /me/logout
- **Request Header**

```bash
Authorization: Bearer your_jwt_token_here
```

- **Responses:**

  - 200 OK: { "message": "Successfully logged out" }
  - Other Status Codes: { "error": "error_message_here" }

## User Routes

### 1. List Users

- **Route:** GET /users
- **Responses:**

  - 200 OK: Array of user objects.
  - 500 Internal Server Error: { "error": "error_message_here" }

### 2. Get a User

- **Route:** GET /user/:uuid
- **Responses:**

  - 200 OK: User object.
  - 404 Not Found: { "error": "User not found" }
  - 500 Internal Server Error: { "error": "error_message_here" }

### 3. Create a User

- **Route:** POST /user
- **Request Body:**

```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
}
```

- **Responses:**

  - 200 OK: { "message": "user created John Doe" }
  - 400 Bad Request: { "error": "Missing required fields" }
  - 500 Internal Server Error: { "error": "error_message_here" }

### 4. Update a User

- **Route:** PUT /me/:uuid
- **Request Body (Partial updates allowed):**

```json
{
    "name": "New Name",
    "email": "newemail@example.com"
}
```

- **Request Header**

```bash
Authorization: Bearer your_jwt_token_here
```

- **Responses:**

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

- **Route:** DELETE /me/:uuid
- **Request Header**

```bash
Authorization: Bearer your_jwt_token_here
```

- **Responses:**

  - 200 OK: { "message": "user deleted" }
  - 500 Internal Server Error: { "error": "error_message_here" }
