package main

import (
	"database/sql"
	"fmt"
	"govideoconverter/internal/converter"
	"govideoconverter/internal/rabbitmq"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

func connectPostgres() (*sql.DB, error) {
	user := getEnvOrDefault("POSTGRES_USER", "postgres")
	password := getEnvOrDefault("POSTGRES_PASSWORD", "postgres")
	host := getEnvOrDefault("POSTGRES_HOST", "postgres")
	port := getEnvOrDefault("POSTGRES_PORT", "5432")
	dbName := getEnvOrDefault("POSTGRES_DB", "converter")
	sslmode := getEnvOrDefault("POSTGRES_SSLMODE", "disable")

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s", user, password, host, port, dbName, sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Error connecting to database", slog.String("conn_str", connStr))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		slog.Error("Error pinging database", slog.String("conn_str", connStr))
		return nil, err
	}
	
	slog.Info("Connected to database successfully")

	return db, nil
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	db, err := connectPostgres()
	if err != nil {
		panic(err)
	}

	rabbitmqUrl := getEnvOrDefault("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	rabbitClient, err := rabbitmq.NewRabbitClient(rabbitmqUrl)
	if err != nil {
		panic(err)
	}

	defer rabbitClient.Close()

	convertionExchange := getEnvOrDefault("CONVERSION_EXCHANGE", "conversion_exchange")
	queueName := getEnvOrDefault("CONVERSION_QUEUE", "video_conversion_queue")
	conversionKey := getEnvOrDefault("CONVERSION_KEY", "conversion")
	confirmationKey := getEnvOrDefault("CONFIRMATION_KEY", "finish-conversion")
	confirmationQueue := getEnvOrDefault("CONFIRMATION_QUEUE", "video_confirmation_queue")

	vc := converter.NewVideoConverter(rabbitClient, db)
	// vc.Handle([]byte(`{"video_id": 1, "path": "/media/uploads/1"}`))

	msgs, err := rabbitClient.ConsumeMessages(convertionExchange, conversionKey, queueName)
	if err != nil {
		slog.Error("Failed to consume messages", slog.String("error", err.Error()))
	}

	for d := range msgs {
		go func(delivery amqp.Delivery) {
			vc.Handle(delivery, convertionExchange, confirmationKey, confirmationQueue)
		}(d)
	}

}

