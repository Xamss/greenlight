package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// Add maxOpenConns, maxIdleConns and maxIdleTime fields to hold the configuration
// settings for the connection pool.
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	logger *log.Logger
}

func main() {

	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// Use the value of the GREENLIGHT_DB_DSN environment variable as the default value
	// for our db-dsn command-line flag.
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	// Read the connection pool settings from command-line flags into the config struct.
	// Notice the default values that we're using?
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	pool, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	// Defer a call to db.Close() so that the connection pool is closed before the
	// main() function exits.
	defer pool.Close()
	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.Printf("database connection pool established")
	app := &application{
		config: cfg,
		logger: logger,
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	// Because the err variable is now already declared in the code above, we need
	// to use the = operator here, instead of the := operator.
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*pgxpool.Pool, error) {

	pool, err := pgxpool.New(context.Background(), cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	pool.Config().MaxConnIdleTime, err = time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	pool.Config().MaxConns = int32(cfg.db.maxOpenConns)

	return pool, nil
}
