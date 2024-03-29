# Backend test task

## Table of contents

- [Description](#description)
- [Technologies](#technologies)
- [Installation](#installation)


## Description
This test task involves creating a user with augmented data from three different open APIs. There are also methods to retrieve all created users with sorting, to retrieve user by id, also update user fields and delete

### Endpoints
- GET: http://localhost:1234/users (all users)
- GET: http://localhost:1234/users?sort_by=name&sort_order=asc (you can sort by all fields: name, age... asc/desc)
> http://localhost:1234/users?sort_by=surname&age&sort_order=asc (sort by surname and age)
- POST: http://localhost:1234/users
```json
{
    "name": "Vitalik",
    "surname": "Buterin",
    "patronymic": "can be null"
}
```
- GET: http://localhost:1234/users/:uuid
- PUT: http://localhost:1234/users/update/:uuid (you can change all fields)
- DELETE: http://localhost:1234/users/delete/:uuid (you can delete users by ids)

## Technologies
- Go
- PostgreSQL
- Docker

## Installation
#### Installation by docker

(NOTE: make sure that you have created postgres image)

```bash
$ docker login

$ docker compose up
```
