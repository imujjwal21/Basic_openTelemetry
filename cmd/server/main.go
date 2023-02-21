package main

import (
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
)

func main() {

	db := sqldb.ConnectDB()
	var port int

	flag.IntVar(&port, "port", 0, "Address to bind the socket on.")

	flag.Parse()

	server := &http.Server{Handler: httptransport.NewHandler(users.NewInMemory(db))}

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
