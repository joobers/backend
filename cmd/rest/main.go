package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/joobers/backend/pkg/models/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	clients  map[primitive.ObjectID][]*websocket.Conn
	messages mongodb.MessageModel
}

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := server.ListenAndServe()
	errorLog.Fatal(err)
}
