package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	opentracing "github.com/opentracing/opentracing-go"
	
	"github.com/opentracing/opentracing-go/ext"

	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/ubser/jaeger-lib/matrics"
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

// initialize a tracer

tracer := opentracing.GlobalTracer()

cfg := &config/Configuration{
	serviceName: "client",
	sampler: &config.SamplerConfig{
		Type: "const",
		Param: 1,
	},
	Reporter: &config.ReporterConfig{
		LogSpan: true,
	},
}

// 1
tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
defer closer.close()

if err!=nil{
	panic(fmt.Sprintf("Can't init jaeger : %v\n", err))
}

// 2

clientSpan := tracer.StartSpan("clientspan")
defer clientSpan.Finish()
time.Sleep(time.Second)


func (i *inMemory) Check(ctx context.Context, name, password string) error {

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

	// ext.SpanKindRPCClient.set(clientSpan)
	// ext.HTTPUrl.Set(clientSpan,url)
    // ext.HTTPMethod.Set(clientSpan, "GET")

	tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, _ := http.DefaultClient.Do(req)
	fmt.Println(resp)

	return nil
}
