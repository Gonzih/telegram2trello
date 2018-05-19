package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

type sessionStore struct {
	redisClient *redis.Client
}

func newSessionStore() (*sessionStore, error) {
	store := sessionStore{}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()

	if err != nil {
		return &store, err
	}

	store.redisClient = client

	return &store, nil
}

func (session *sessionStore) key(id int) string {
	return fmt.Sprintf("session-%d", id)
}

func (session *sessionStore) Set(id int, vals map[string]interface{}) error {
	res := session.redisClient.HMSet(session.key(id), vals)
	return res.Err()
}

func (session *sessionStore) Get(id int, key string) (string, error) {
	res, err := session.redisClient.HMGet(session.key(id), key).Result()

	if err != nil {
		return "", err
	}

	if len(res) == 1 {
		v, ok := res[0].(string)

		if !ok {
			return "", fmt.Errorf("Cant convert %v to string", res[0])
		}

		return v, nil
	}

	return "", fmt.Errorf("Got weird result from redis %v", res)
}

func (session *sessionStore) Clear(id int) error {
	return session.redisClient.Del(session.key(id)).Err()
}
