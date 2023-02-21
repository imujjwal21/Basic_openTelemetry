package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
)

type inMemory struct {
	db *sql.DB
}

func NewInMemory(db *sql.DB) Storage {
	return &inMemory{db}
}

func (i *inMemory) Create(ctx context.Context, name, password string) error {

	log.Print("Register -> ", name, " : ", password)
	log.Print("\nDatabase -> ", i.db)

	_, err := i.db.Exec(`INSERT INTO users (name, password) VALUES (?,?)`, name, password)
	if err != nil {
		return fmt.Errorf("can't insert user : %v", err)
	}
	return nil

}

func (i *inMemory) Check(ctx context.Context, name, password string) error {

	_, span := otel.Tracer(name).Start(ctx, "Check")
	defer span.End()

	row, err := i.db.Query(`SELECT id FROM users WHERE name=? and password=?`, name, password)

	log.Print("Login ->  ", name, " : ", password)

	if err != nil {
		return fmt.Errorf("can't retrieved data from database : %v", err)
	}

	var flag bool

	for row.Next() {
		flag = true
	}

	if !flag {
		return errors.New("can't connect to database with id")
	}

	return nil
}
