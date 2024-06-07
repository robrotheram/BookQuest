package auth

import (
	"BookQuest/internal/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/allegro/bigcache/v3"
)

type Cache struct {
	*bigcache.BigCache
}

func initCache() *Cache {
	config := bigcache.DefaultConfig(10 * time.Minute) // set cache expiration time
	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	return &Cache{
		BigCache: cache,
	}
}

func (cache *Cache) SetUser(user models.User) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return cache.Set(user.Email, userBytes)
}

func (cache *Cache) GetUser(email string) (models.User, error) {
	userBytes, err := cache.Get(email)
	if err != nil {
		return models.User{}, err
	}
	var user models.User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
