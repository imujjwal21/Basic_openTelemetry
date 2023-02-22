package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
)

type inMemory struct {
	db *sql.DB
}

func NewInMemory(db *sql.DB) Storage {
	return &inMemory{db}
}

const name = "login"

type Telemet interface {
	Done(ctx context.Context) error
}

type checkedLogin struct {
}

func (l *checkedLogin) Done(ctx context.Context) error {
	_, span := otel.Tracer(name).Start(ctx, "checkedLogin")
	defer span.End()
	return nil
}

type checkedQuery struct {
	first Telemet
}

func (c *checkedQuery) Done(ctx context.Context) error {
	_, span := otel.Tracer(name).Start(ctx, "checkedQuery")
	defer span.End()
	return c.first.Done(ctx)
}

type insideCheckFun struct {
	second Telemet
}

func (i *insideCheckFun) Done(ctx context.Context) error {
	_, span := otel.Tracer(name).Start(ctx, "insideCheckFun")
	defer span.End()
	return i.second.Done(ctx)
}

func (i *inMemory) Create(ctx context.Context, name, password string) error {

	_, err := i.db.Exec(`INSERT INTO users (name, password) VALUES (?,?)`, name, password)
	if err != nil {
		return fmt.Errorf("can't insert user : %v", err)
	}
	return nil

}

func (i *inMemory) Check(ctx context.Context, name, password string) error {

	var temp Telemet
	{
		temp = &checkedLogin{}
	}

	row, err := i.db.Query(`SELECT id FROM users WHERE name=? and password=?`, name, password)

	if err != nil {
		return fmt.Errorf("can't retrieved data from database : %v", err)
	}

	temp = &checkedQuery{temp}

	var flag bool

	for row.Next() {
		flag = true
	}

	if !flag {
		return errors.New("can't connect to database with id")
	}

	temp = &insideCheckFun{temp}

	temp.Done(ctx)

	return nil
}
