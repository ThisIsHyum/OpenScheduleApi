package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/config"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

const dbName string = "osadb"

var testDB *gorm.DB
var container *mysql.MySQLContainer

func truncateAll() {
	var tables []string

	testDB.Raw(`
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = DATABASE();
    `).Scan(&tables)

	testDB.Exec("SET FOREIGN_KEY_CHECKS = 0;")

	for _, t := range tables {
		testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", t))
	}

	testDB.Exec("SET FOREIGN_KEY_CHECKS = 1;")
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	container, err = mysql.Run(ctx,
		"mysql:8.0",
		mysql.WithDatabase("osadb"),
		mysql.WithUsername("root"),
		mysql.WithPassword("root"),
		testcontainers.WithReuseByName("test_osadb"),
	)
	if err != nil {
		panic(err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "3306/tcp")

	testDB, err = database.NewDb(&config.DatabaseConfig{
		Host:     host,
		Port:     int(port.Num()),
		User:     "root",
		Password: "root",
		Dbname:   dbName,
	})
	if err != nil {
		panic(err)
	}

	code := m.Run()

	container.Terminate(ctx)
	os.Exit(code)
}
