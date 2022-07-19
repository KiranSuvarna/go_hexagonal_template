package driven

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/hex_microservice_template/internal/core/domain"
	"github.com/hex_microservice_template/internal/core/ports/outbound"
	"github.com/hex_microservice_template/internal/core/usecase"
	"github.com/pkg/errors"
	"strconv"

	"os"
	"time"
)

type redisRepository struct {
	client        *redis.Client
	clusterClient *redis.ClusterClient
	expiration    time.Duration
}

func newRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}

func NewRedisRepository(host string) (outbound.RedirectRepository, error) {
	redisRepo := &redisRepository{}
	client, err := newRedisClient()
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	redisRepo.client = client
	return redisRepo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r redisRepository) Find(code string) (*domain.Redirect, error) {
	redirect := &domain.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(usecase.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r redisRepository) Store(redirect *domain.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
