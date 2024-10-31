package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

var databaseInstance *sqlstore.Container

func Initialize() {
	fmt.Println("===== Initialize Database =====")

	// Replace the DSN with your own PostgreSQL DSN
	dsn := "postgres://<username>:<password>@<host>:<port>/<db_name>?sslmode=disable"
	container, err := sqlstore.New("pgx", dsn, nil)
	if err != nil {
		panic(fmt.Errorf("failed to initialize sqlstore: %w", err))
	}

	databaseInstance = container
	fmt.Println("connected to: waclient.db")
}

func GetDBInstance() *sqlstore.Container {
	return databaseInstance
}
