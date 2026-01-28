package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// Order represents a trading order
type Order struct {
	ID        string    `json:"id"`
	Symbol    string    `json:"symbol"`
	Side      string    `json:"side"`
	Price     float64   `json:"price"`
	Quantity  float64   `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	redisClient *redis.Client
	ctx         = context.Background()
	ordersCount = 0
	mu          sync.Mutex
)

func init() {
	// Load environment variables
	_ = godotenv.Load()

	// Redis connection
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:         redisHost + ":" + redisPort,
		PoolSize:     10,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("✓ Connected to Redis")
}

// PublishOrder publishes an order to Redis Stream
func PublishOrder(order Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	// Fast path: publish to stream with minimal serialization
	result := redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "orders.fast",
		ID:     "*", // Auto-generate ID
		Values: map[string]interface{}{
			"data": string(data),
		},
	})

	if err := result.Err(); err != nil {
		return fmt.Errorf("failed to publish order: %w", err)
	}

	mu.Lock()
	ordersCount++
	mu.Unlock()

	return nil
}

// PreRiskCheck performs O(1) pre-risk validation
func PreRiskCheck(order Order) bool {
	// Stupid-fast validations
	if order.Price <= 0 || order.Quantity <= 0 {
		return false
	}
	if order.Symbol == "" {
		return false
	}
	if order.Side != "BUY" && order.Side != "SELL" {
		return false
	}
	return true
}

// ProcessOrder processes a single order
func ProcessOrder(order Order) error {
	// Pre-risk check (O(1))
	if !PreRiskCheck(order) {
		log.Printf("❌ Pre-risk check failed for order %s", order.ID)
		return fmt.Errorf("pre-risk check failed")
	}

	// Publish to Redis Stream
	if err := PublishOrder(order); err != nil {
		log.Printf("❌ Failed to publish order %s: %v", order.ID, err)
		return err
	}

	log.Printf("✓ Order published: %s | %s %s @ %.2f x %.0f",
		order.ID, order.Side, order.Symbol, order.Price, order.Quantity)

	return nil
}

// HealthCheck returns service health status
func HealthCheck() map[string]interface{} {
	mu.Lock()
	count := ordersCount
	mu.Unlock()

	return map[string]interface{}{
		"status":        "healthy",
		"orders_sent":   count,
		"timestamp":     time.Now().Unix(),
		"redis_status":  "connected",
	}
}

// ProcessBatch processes multiple orders in parallel (fast!)
func ProcessBatch(orders []Order) {
	// Use goroutines for parallel processing
	semaphore := make(chan struct{}, 100) // Max 100 concurrent
	var wg sync.WaitGroup

	for _, order := range orders {
		wg.Add(1)
		go func(o Order) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			if err := ProcessOrder(o); err != nil {
				log.Printf("Error processing order: %v", err)
			}
		}(order)
	}

	wg.Wait()
}

// SimulateLiveOrders simulates incoming orders (for testing)
func SimulateLiveOrders(count int, interval time.Duration) {
	symbols := []string{"BTC", "ETH", "SOL", "USDT"}
	sides := []string{"BUY", "SELL"}

	for i := 0; i < count; i++ {
		order := Order{
			ID:        fmt.Sprintf("ORD-%d", i+1),
			Symbol:    symbols[i%len(symbols)],
			Side:      sides[i%len(sides)],
			Price:     45000 + float64(i),
			Quantity:  float64((i%10) + 1),
			Timestamp: time.Now(),
		}

		if err := ProcessOrder(order); err != nil {
			log.Printf("Error: %v", err)
		}

		time.Sleep(interval)
	}
}

func main() {
	log.Println("🚀 Order Core Service (Go) - Stupid-Fast Edition")
	log.Println("================================================")

	// Health check
	health := HealthCheck()
	log.Printf("Health: %+v\n", health)

	// Simulate orders or wait for API calls
	// For demo: publish 10 orders
	go func() {
		log.Println("\n📤 Publishing 10 sample orders...")
		SimulateLiveOrders(10, 100*time.Millisecond)
	}()

	// Keep service running
	select {}
}
