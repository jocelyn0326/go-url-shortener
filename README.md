# ShortenUrl- 郭宜萱(jocelyn0326@gmail.com)


## Summary
Implement a one-to-one url shorten service.


## Built With
* [GO](https://go.dev/)
* [Gin](https://github.com/gin-gonic/gin)
* [Redis](https://pkg.go.dev/github.com/go-redis/redis/v8#section-readme)
* [Docker](https://www.docker.com/)

## How to start the server

`docker compose up`

## How to run test

1. `docker compose -f docker-compose-test.yaml build`
2. `docker compose -f docker-compose-test.yaml run test-funnow-api-server go test -v`

## API Introduction

### `POST |shorten`
url: http://localhost:8080/shorten
| Property | Type | Description |
| -------- | -------- | -------- |
| longUrl     | string     | The original url which will be shorten and be kept in cache. |


```
{
  "longurl": "https://www.google.com.tw/"
}

```

## How to avoid race condition

Use Redis watch and transaction to wrap get and set actions in `controllers/url.go`.
    
## Future TODO list
- [ ] Redirect short url.
- [ ]  LRU cache implementation
- [ ] Analyze user behavior with using RDB.




