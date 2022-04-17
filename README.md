# Instruction

1. cp .env_example .env
2. docker-compose up -d
3. Load `127.0.0.1:8889/posts`


# API:

* GET `127.0.0.1:8889/posts?limit=2&page=1`
* GET `127.0.0.1:8889/posts/:id`
* POST `/posts`
```json
    {"name":"required,gt=2,lt=10"}
```
* PUT `/posts/:id`
```json
    {"name":"required,gt=2,lt=10"}
```
* DELETE `/posts/:id`
