package main

import (
	"fmt"
	"net/http"
	"os"
	"youtube-stream-notifier-linebot/api"
	"youtube-stream-notifier-linebot/controller"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// init
	controller.InitLineBot()
	controller.InitDatabase()
	defer controller.DB.Close()

	// automate subscription
	api.Subscribe()

	// create router and server
	router := mux.NewRouter()
	router.HandleFunc("/", api.Hello)
	router.HandleFunc("/callback/", api.LineCallbackHandler)
	router.HandleFunc("/subscribe/", api.PubsubCallbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, router)
}
