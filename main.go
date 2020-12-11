package main

import (
	"database/sql"
	"fmt"
	"github.com/Alex-ttt/pokerbank/handlers"
	"github.com/Alex-ttt/pokerbank/middlewares"
	"github.com/Alex-ttt/pokerbank/models"
	"github.com/Alex-ttt/pokerbank/services"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"os"
)

// initTCPConnectionPool initializes a TCP connection pool for a Cloud SQL
// instance of SQL Server.
func initTCPConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_postgres_databasesql_create_tcp]
	var (
		dbUser    = mustGetenv("DB_USER")
		dbPwd     = mustGetenv("DB_PASS")
		dbTcpHost = mustGetenv("DB_TCP_HOST")
		dbPort    = mustGetenv("DB_PORT")
		dbName    = mustGetenv("DB_NAME")
	)

	var dbURI string
	dbURI = fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s", dbTcpHost, dbUser, dbPwd, dbPort, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_postgres_databasesql_create_tcp]
}

// initSocketConnectionPool initializes a Unix socket connection pool for
// a Cloud SQL instance of SQL Server.
func initSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_postgres_databasesql_create_socket]s
	var (
		dbUser                 = mustGetenv("DB_USER")
		dbPwd                  = mustGetenv("DB_PASS")
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME")
		dbName                 = mustGetenv("DB_NAME")
	)

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dbURI string
	dbURI = fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", dbUser, dbPwd, dbName, socketDir, instanceConnectionName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_postgres_databasesql_create_socket]
}

func main() {
	var (
		db  *sql.DB
		err error
	)

	databaseUrl := os.Getenv("DATABASE_URL")
	if len(databaseUrl) > 0 {
		db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalf("Error opening database: %q", err)
		}
	} else {
		if os.Getenv("DB_TCP_HOST") != "" {
			db, err = initTCPConnectionPool()
			if err != nil {
				log.Fatalf("initTCPConnectionPool: unable to connect: %s", err)
			} else {
				log.Println("initTCPConnectionPool: success")
			}
		} else {
			db, err = initSocketConnectionPool()
			if err != nil {
				log.Fatalf("initSocketConnectionPool: unable to connect: %s", err)
			} else {
				log.Println("initSocketConnectionPool: success")
			}
		}
	}

	services.Db = db
	models.CreateDatabaseStructure(db)

	services.InitRedis()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	var router = gin.Default()
	router.Static("/static", "./static")
	templateFuncs := template.FuncMap{
		"add": func(x, y int) int { return x + y },
		"seq": func(n int) []int { return make([]int, n, n) },
	}

	router.SetFuncMap(templateFuncs)
	router.LoadHTMLFiles("templates/index.html")

	router.GET("/", middlewares.TokenAuthMiddleware, handlers.IndexPage)
	router.POST("/addGameResult", middlewares.TokenAuthMiddleware, handlers.AddGameResult)
	router.POST("/payDebts", middlewares.TokenAuthMiddleware, handlers.AddDebtPayment)
	router.POST("/login", handlers.Login)
	router.POST("/logout", middlewares.TokenAuthMiddleware, handlers.Logout)
	//http.HandleFunc("/", handlers.IndexPage)
	//http.HandleFunc("/payDebts", handlers.AddDebtPayment)

	if err = router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func configureConnectionPool(dbPool *sql.DB) {
	// [START cloud_sql_postgres_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	dbPool.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	dbPool.SetMaxOpenConns(7)

	// [END cloud_sql_postgres_databasesql_limit]

	// [START cloud_sql_postgres_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	dbPool.SetConnMaxLifetime(1800)

	// [END cloud_sql_postgres_databasesql_lifetime]
}
