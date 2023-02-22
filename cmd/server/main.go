package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"login/internal/httptransport"
	"login/sqldb"
	"login/users"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func newExporter() (*jaeger.Exporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:16686")))

	return exp, err
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("login"),
		),
	)
	return r
}

func main() {

	db := sqldb.ConnectDB()
	var port int

	flag.IntVar(&port, "port", 0, "Address to bind the socket on.")

	flag.Parse()

	server := &http.Server{Handler: httptransport.NewHandler(users.NewInMemory(db))}

	l := log.New(os.Stdout, "", 0)

	exp, err := newExporter()
	if err != nil {
		l.Fatal(err)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(newResource()),
	)

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()

	otel.SetTracerProvider(tp)

	go func() {

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

		if err != nil {
			log.Panicf("cannot create tpc listener: %v", err)
		}

		log.Printf("      starting http server on %q", lis.Addr())
		if err := server.Serve(lis); err != nil {
			log.Panicf("cannot start http server: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	log.Printf("Got exit signal %q. Bye", <-sig)
}
