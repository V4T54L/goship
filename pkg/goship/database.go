package goship

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	_ "modernc.org/sqlite"
)

// ConnectToPostgresDb establishes a connection to a PostgreSQL database.
// It accepts a connection string and returns a *sql.DB object and an error if any.
//
// Example usage:
//
//	connStr := "postgres://user:password@localhost:5432/mydatabase?sslmode=disable"
func ConnectToPostgresDb(connStr string) (*sql.DB, error) {
	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err == nil {
		// Ping the database to ensure the connection is established
		err = db.Ping()
	}

	return db, err
}

// ConnectToRedisDb establishes a connection to a Redis database.
// It accepts a context and a Redis connection string, and returns a *redis.Client object and an error if any.
//
// Example usage:
//
//	connStr := "redis://localhost:6379/0"
//	ctx := context.Background()
func ConnectToRedisDb(ctx context.Context, connStr string) (*redis.Client, error) {
	// Parse the Redis connection string into options
	opt, err := redis.ParseURL(connStr)
	if err != nil {
		return nil, err
	}

	// Create a new Redis client with the parsed options
	client := redis.NewClient(opt)
	// Ping the Redis server to check the connection
	status := client.Ping(ctx)
	return client, status.Err()
}

// ConnectToSqliteDb establishes a connection to an SQLite database.
// It accepts the file path to the SQLite database and returns a *sql.DB object and an error if any.
//
// Example usage:
//
//	filePath := "/path/to/database.db"
func ConnectToSqliteDb(filePath string) (*sql.DB, error) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite", filePath)
	if err == nil {
		// Ping the database to ensure the connection is established
		err = db.Ping()
	}

	return db, err
}

// ConnectToMongoDB establishes a connection to a MongoDB database.
// It accepts a context and a MongoDB URI, and returns a *mongo.Client object and an error if any.
//
// Example usage:
//
//	dbURI := "mongodb://localhost:27017"
//	ctx := context.Background()
func ConnectToMongoDB(ctx context.Context, dbURI string) (*mongo.Client, error) {
	// Establish a connection to MongoDB using the provided URI
	client, err := mongo.Connect(options.Client().ApplyURI(dbURI))
	if err == nil {
		// Ping the MongoDB server to check the connection
		err = client.Ping(ctx, nil)
	}
	return client, err
}

func CloseMongoDBConn(ctx context.Context, client *mongo.Client) error {
	return client.Disconnect(ctx)
}

func CloseSqlDBConn(db *sql.DB) error {
	return db.Close()
}

func CloseRedisConn(client *redis.Client) error {
	return client.Close()
}
