package main

import (
	"backend/internal/models"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Contact)

func (app *application) wsConnectionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	app.ws = ws
	clients[ws] = true

	contacts, err := app.models.Contacts.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	ws.WriteJSON(map[string]any{"type": "allContacts", "contacts": contacts})

	for {
		var contact models.Contact
		err := ws.ReadJSON(&contact)
		if err != nil {
			log.Printf("error wsConnectionHandler: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- contact
	}

}

func (app *application) writeWsContact(t string, c *models.Contact) {
	for client := range clients {
		err := client.WriteJSON(map[string]any{"type": t, "contact": c})
		if err != nil {
			log.Printf("error handleContact: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func (app *application) pushContact() {
	for {
		c := <-broadcast
		for client := range clients {
			err := client.WriteJSON(c)
			if err != nil {
				log.Printf("error pushContact: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
