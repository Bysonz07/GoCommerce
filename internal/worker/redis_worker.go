package worker

import (
	"context"
	"log"

	"github.com/redis/go-redis/v8"
)

func StartRedisWorker(redisClient *redis.Client) {
	ctx := context.Background()

	for {
		job, err := redisClient.RPop(ctx, "order_jobs").Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			log.Printf("Failed to fetch job from Redis: %v", err)
			continue
		}

		log.Printf("Processing job: %s", job)
		// Process the job (e.g., send email, update status, etc.)
	}
}
