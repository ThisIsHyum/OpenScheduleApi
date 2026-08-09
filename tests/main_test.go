package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ThisIsHyum/OpenScheduleApi/internal/config"
	"github.com/ThisIsHyum/OpenScheduleApi/internal/database"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

const dbName string = "osadb_test"

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
		mysql.WithDatabase(dbName),
		mysql.WithUsername("root"),
		mysql.WithPassword("root"),
		testcontainers.WithReuseByName("test_osadb"),
	)
	if err != nil {
		panic(err)
	}
	defer container.Terminate(ctx)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "3306/tcp")

	if err := migrateDB(host, int(port.Num())); err != nil {
		panic(err)
	}

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
	os.Exit(code)
}

func migrateDB(host string, port int) error {

	m, err := migrate.New(
		"file://../migrations",
		fmt.Sprintf("mysql://root:root@tcp(%s:%d)/%s?multiStatements=true",
			host, port, dbName),
	)
	if err != nil {
		return err
	}
	defer m.Close()
	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}
	return err

}
