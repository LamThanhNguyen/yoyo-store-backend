package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WebSocketConnection struct {
	*websocket.Conn
}

type WsPayload struct {
	Action      string              `json:"action"`
	Message     string              `json:"message"`
	UserName    string              `json:"username"`
	MessageType string              `json:"message_type"`
	UserID      int                 `json:"user_id"`
	Conn        WebSocketConnection `json:"-"`
}

type WsJsonResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

var clients = make(map[WebSocketConnection]string)

var wsChan = make(chan WsPayload, 8)

func (server *Server) ListenToWsChannel(ctx context.Context) {
	var response WsJsonResponse
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("WsChannel stopped")
			return
		case e := <-wsChan:
			switch e.Action {
			case "deleteUser":
				response.Action = "logout"
				response.Message = "Your account has been deleted"
				response.UserID = e.UserID
				server.broadcastToAll(response)
			default:
			}
		}
	}
}

func (server *Server) broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		// broadcast to every connected client
		err := client.WriteJSON(response)
		if err != nil {
			log.Error().Err(err).Str("action", response.Action).Msg("Websocket error")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (server *Server) WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err)
		return
	}

	log.Info().Str("RemoteAddr", r.RemoteAddr).Msg("Client connected")
	var response WsJsonResponse
	response.Message = "Connected to server"

	err = ws.WriteJSON(response)
	if err != nil {
		log.Error().Err(err)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	go server.ListenForWS(&conn)
}

func (server *Server) ListenForWS(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msg(fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}
