package main

import "net/http"

func (app *application) upgradeConnection(w http.ResponseWriter, r *http.Request) {
	app.upgradeToWebsocket(w, r)
}
