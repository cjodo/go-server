package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/cjodo/go-server/pkg/lobby"
	"github.com/cjodo/go-server/pkg/socket"
	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "4040"
	}
	
	l := lobby.NewLobby()

	go l.Run()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/set-session", setSessionHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connect(w, r, l)
	})

	fmt.Println("server now running on port :", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:" + port, nil))
}

func connect(w http.ResponseWriter, r *http.Request, l *lobby.Lobby) {
	token, err := r.Cookie("token")

	if err != nil {
		// won't connect unless they have a valid session token
		fmt.Println("Error getting cookie or does not exist")
		return
	} else {
		fmt.Println("token: ", token)
	}

	conn, err := socket.Upgrade(w, r)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		return
	}

	connection := &socket.Connection{
		Id: token.String(),
		Conn: conn,
		Send: make(chan interface{}),
	}

	l.Register <- connection
}

func setSessionHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("session req received")

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	newToken := uuid.New()

	cookie := &http.Cookie{
		Name: "token",
		Value: newToken.String(),
		Path: "/",
		HttpOnly: false,
		MaxAge: 3600,
		Secure: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}
