package database


import (
    "context"
    "github.com/go-redis/redis/v8"
)

func NewRedisClient(addr, password string) (*redis.Client, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       0,
    })

    // Test the connection
    _, err := client.Ping(context.Background()).Result()
    if err != nil {
        return nil, err
    }

    return client, nil
}