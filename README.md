## Simple REST API in GO for todo lists app.

## Table of Contents

- [API Endpoints](#api-endpoints)
- [Build](#build)
- [Run](#run)
- [Database Migrations](#database-migrations)

## API Endpoints

| Port  | URI	Description       | Description                                                  | Method |
|-------|-----------------------|--------------------------------------------------------------|--------|
| :8080 | /auth/sign-up         | sign-up endpoint                                             | POST   |
| :8080 | /auth/sign-in         | sign-in endpoint                                             | POST   |
| :8080 | /auth/refresh-tokens  | refresh access and refresh JWT tokens                        | POST   |
| :8080 | /api/lists            | create a new todo list for the user                          | POST   |
| :8080 | /api/lists            | get all todo lists for the user                              | GET    |
| :8080 | /api/lists/{id}       | get users todo list by id                                    | GET    |
| :8080 | /api/lists/{id}       | update users todo list by id                                 | PUT    |
| :8080 | /api/lists/{id}       | delete users todo list by id                                 | DELETE |
| :8080 | /api/lists/{id}/items | create a new todo item for a todo list belonging to the user | POST   |
| :8080 | /api/lists/{id}/items | get all todo list items for the user by list id              | GET    |
| :8080 | /api/items/{id}       | get todo item for the user by item id                        | GET    |
| :8080 | /api/items/{id}       | update todo item for the user by item id                     | PUT    |
| :8080 | /api/items/{id}       | delete users todo item by id                                 | DELETE |
| :8080 | /swagger/doc.json     | Swagger Docs                                                 | GET    |
| :8080 | /swagger/index.html   | Swagger Docs                                                 | GET    |

## Build

To build the Docker image, run the following command:

```bash
make build
```

## Run

To run the Docker container, use the following command:

```bash
make run
```

## Database Migrations
Before running the application, ensure you have the migrate tool installed.  
Follow the [official instructions](https://github.com/golang-migrate/migrate) to install the migrate tool.

After installing the migrate tool, you can apply database migrations with:

```bash
make migrate
```
