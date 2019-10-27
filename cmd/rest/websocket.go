package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joobers/backend/pkg/models"
	"github.com/joobers/backend/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *application) upgradeToWebsocket(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	claims, err := utils.ParseJwt(token)
	if err != nil {
		app.serverError(w, err)
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Printf("User with ID '%s' connected", claims.ID.String())

	app.clients[claims.ID] = append(app.clients[claims.ID], ws)

	app.handleMessages(w, ws, claims.ID)
}

func (app *application) handleMessages(w http.ResponseWriter, conn *websocket.Conn, id primitive.ObjectID) {
	collection, ctx, cancel := app.messages.GetCollection(id)

	for {
		messageType, messageBytes, err := conn.ReadMessage()
		if err != nil {
			app.serverError(w, err)
			continue
		}

		if messageType == websocket.TextMessage {
			var message models.Message
			err := json.Unmarshal(messageBytes, &message)
			if err != nil {
				app.clientError(w, http.StatusBadRequest)
				continue
			}

			result, err := app.messages.Insert(id, message, collection, ctx, cancel)
			if err != nil {
				app.serverError(w, err)
				continue
			}

			message.ID = result.InsertedID.(primitive.ObjectID)
			messageBytes, err := json.Marshal(message)
			if err != nil {
				app.serverError(w, err)
				continue
			}

			deadline := time.Now().Add(time.Minute)
			for _, client := range app.clients[id] {
				client.WriteControl(messageType, messageBytes, deadline)
			}
		}

	}
}
