package models

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type Payment struct {
	Correlation_id string
	Amount         float64
	Processor      string
}

var Redis *redis.Client

func RedisConnect() error {

	addr := os.Getenv("REDIS_URL")

	Redis = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	_, err := Redis.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("failed to ping: %v", err)
	}

	return nil

}

func InsertPayment(p Payment) error {

	ctx := context.Background()
	pipe := Redis.TxPipeline()

	key := "payment:" + p.Correlation_id

	pipe.HSet(ctx, key, map[string]interface{}{
		"amount":    p.Amount,
		"processor": p.Processor,
	})

	pipe.Incr(ctx, "total_requests:"+p.Processor)
	pipe.IncrByFloat(ctx, "total_amount:"+p.Processor, p.Amount)

	_, err := pipe.Exec(ctx)

	return err

}

func GetSummary() (map[string]map[string]interface{}, error) {

	ctx := context.Background()
	result := make(map[string]map[string]interface{})

	processors := []string{"default", "fallback"}

	for _, proc := range processors {

		reqs, _ := Redis.Get(ctx, "total_requests:"+proc).Int()
		amount, _ := Redis.Get(ctx, "total_amount:"+proc).Float64()

		result[proc] = map[string]interface{}{
			"totalRequests": reqs,
			"totalAmount":   amount,
		}
	}

	return result, nil

}

func PurgeRedisData() error {
	ctx := context.Background()

	Redis.FlushAll(ctx)

	return nil
}
