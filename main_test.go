package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// Test_pure_sql выполняет 10 задание (интеграционный тест создаёт контейнер с БД, выполняет операции)
func Test_pure_sql(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	container, err := NewPostgreSQLContainer(ctx)

	if err != nil {
		t.Fatal(err)
	}
	// миграция (создание тестовой базы данных)
	_, _, err = container.Exec(ctx, []string{"psql", "-U", "postgres", "-d", "postgres", "-c", "CREATE TABLE test_db (id int, first_name varchar, last_name varchar)"})
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(context.Background())
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", "postgres", "1234", container.Host, container.MappedPort, "postgres")
	if err := pure_sql(dsn); err != nil {
		t.Fatal(err)
	}
}
