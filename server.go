package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/harshitrajsinha/go-get-job/driver"
	"github.com/harshitrajsinha/go-get-job/graph"
	"github.com/harshitrajsinha/go-get-job/store"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8000"

//go:embed store/schema.sql
var schemaFS embed.FS
var db *sql.DB

type dbConfig struct {
	User string `envconfig:"DB_USER"`
	Host string `envconfig:"DB_HOST"`
	Port string `envconfig:"DB_PORT"`
	Pass string `envconfig:"DB_PASS"`
	Name string `envconfig:"DB_NAME"`
}

// Function to load data to database via schema file
func loadDataToDatabase(dbClient *sql.DB) error {

	// Read file content
	sqlFile, err := schemaFS.ReadFile("store/schema.sql")
	fmt.Println("...loading schema file")
	// sqlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Execute file content (queries)
	_, err = dbClient.Exec(string(sqlFile))
	if err != nil {
		return err
	}
	return nil
}

func init() {

	var cfg dbConfig
	var err error

	// load environment variables
	_ = godotenv.Load()
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=30", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
	dbDriver := "postgres"

	// Get db client
	db, err = driver.InitDB(dbDriver, connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Load data into database
	err = loadDataToDatabase(db)
	if err != nil {
		panic(err)
	} else {
		log.Println("SQL file executed successfully!")
	}

}

func main() {

	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbStore := store.NewJobStore(db)
	gqlQueryResolver := graph.NewGQLQueryResolver(dbStore)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: gqlQueryResolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
