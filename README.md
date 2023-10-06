# Notes Service

The Notes Service is a RESTful API that allows users to perform various operations related to note-taking. It provides endpoints for user authentication and for performing CRUD operations on notes. It also imposes a limit of 30 notes per user.

## Features
- [x] Authentication system
- [x] CRUD system on notes
- [ ] Share notes with other users
- [x] Toggle note visibility (public/private)
- [x] View list of public notes
- [ ] Add notes to personal favorites
- [ ] Save notes from other users
- [ ] Ban/unban users
- [ ] User permissions management system

## Default User
- **Login:** admin
- **Password:** qwerty12

## Limitations
- Each user can have a maximum of 30 notes.

## Documentation
Full API documentation is available via swagger
http://localhost:6000/v0/swagger

## How to Run
```
make up
```