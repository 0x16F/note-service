# Notes Service

The Notes Service is a RESTful API that allows users to perform various operations related to note-taking. It provides endpoints for user authentication and for performing CRUD operations on notes. It also imposes a limit of 30 notes per user.

## Features
- [x] Unit tests
- [x] Integrations tests
- [x] Swagger

## Default User
- **Login:** admin
- **Password:** qwerty12

## Limitations
- Each user can have a maximum of 30 notes.

## Documentation
Full API documentation is available via swagger
http://localhost:6000/v0/swagger

## How to Use
1. **User Registration and Authentication:**
   - Register a new user using the `/v0/auth/register` endpoint.
   - Authenticate using the `/v0/auth/login` endpoint to get the token.
   - Use this token in the `Authorization` header for subsequent requests.

2. **Creating and Managing Notes:**
   - Use the `/v0/notes` endpoint to create, update, fetch, or delete notes.
   - The authenticated user can perform actions on their notes using the above endpoint.

## How to Run
```
make up
```