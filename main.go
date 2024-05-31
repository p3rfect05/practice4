package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestDB struct {
	ID        int
	LastName  string
	FirstName string
}

type PostgreSQLContainer struct {
	testcontainers.Container
	MappedPort string
	Host       string
}

// NewPostgreSQLContainer создаёт контейнер postgresql
func NewPostgreSQLContainer(ctx context.Context) (*PostgreSQLContainer, error) {
	req := testcontainers.ContainerRequest{
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "1234",
			"POSTGRES_DB":       "postgres",
		},
		ExposedPorts: []string{"5432/tcp"},
		Image:        "postgres:latest",
		WaitingFor: wait.ForExec([]string{"pg_isready"}).
			WithPollInterval(3 * time.Second).
			WithExitCodeMatcher(func(exitCode int) bool {
				return exitCode == 0
			}),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	return &PostgreSQLContainer{
		Container:  container,
		MappedPort: mappedPort.Port(),
		Host:       host,
	}, nil
}

// urlExample := "postgres://username:password@localhost:5432/database_name"
var DATABASE_URL = "postgres://postgres:1234@localhost:5432/postgres"

func main() {

}

// pool_connection_sql выполняет 9 задание
func pool_connection_sql() {
	dbpool, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	var id int
	var firstName, lastName string
	err = dbpool.QueryRow(context.Background(), "select * from test_db;").Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Fatalf("QueryRow all fields failed: %v\n", err)
	}

	fmt.Println(id, firstName, lastName)
}

// gorm_sql выполняет 7 задание
func gorm_sql() error {
	dsn := "host=localhost user=postgres password=1234 dbname=postgres port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print(err)
		return err
	}

	res := db.Create(&TestDB{ID: 20, FirstName: "Super", LastName: "Man"})
	if res.Error != nil {
		log.Print(res.Error)
		return err
	}

	return nil

}

// pure_sql выполняет 1-6 задания
func pure_sql(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO test_db VALUES($1, $2, $3)", 1, "John", "Smith")
	if err != nil {
		log.Print(err)
		return err
	}
	var id int
	var firstName, lastName string
	err = db.QueryRow("select * from test_db;").Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Printf("QueryRow all fields failed: %v\n", err)
		return err
	}

	fmt.Println(id, firstName, lastName)

	err = db.QueryRow("SELECT id, first_name, last_name FROM test_db WHERE last_name = $1;", "Smith").Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Printf("QueryRow with condition failed: %v\n", err)
		return err

	}

	fmt.Println(id, firstName, lastName)

	result, err := db.Exec("INSERT INTO test_db VALUES($1, $2, $3)", 10, "Julia", "Roberts")
	if err != nil {
		log.Printf("Exec INSERT failed: %v\n", err)
		return err

	}
	inserted, err := result.RowsAffected()
	if err != nil {
		log.Printf("Exec RowsAffected INSERT failed: %v\n", err)
		return err

	}
	fmt.Printf("Inserted %d rows\n", inserted)

	result, err = db.Exec("UPDATE test_db SET first_name = $1 where first_name = $2", "Gulia", "Julia")
	if err != nil {
		log.Printf("Exec UPDATE failed: %v\n", err)
		return err

	}
	updated, err := result.RowsAffected()
	if err != nil {
		log.Printf("Exec RowsAffected UPDATE failed: %v\n", err)
		return err

	}
	fmt.Printf("Updated %d rows\n", updated)

	tx, err := db.Begin()
	if err != nil {
		log.Print(err)
		return err

	}
	_, execErr := tx.Exec("INSERT INTO test_db VALUES($1, $2, $3)", 17, "Tony", "Stark")
	if execErr != nil {
		_ = tx.Rollback()
		log.Print(execErr)
		return err

	}

	_, execErr = tx.Exec("UPDATE test_db SET first_name = $1, last_name = $2 WHERE id = $3", "Iron", "Man", 17)
	if execErr != nil {
		_ = tx.Rollback()
		log.Print(execErr)
		return err

	}
	if err := tx.Commit(); err != nil {
		log.Print(err)
		return err
	}

	return nil
}
