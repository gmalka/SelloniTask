package main

import (
	"DobroBot/model"
	customestore "DobroBot/store/customeStore"
	"DobroBot/transport/rest"
	"DobroBot/transport/telegram"
	"net/http"
)

func main() {
	store := customestore.NewStore()
	ch := make(chan (model.Discont), 10)
	tg := telegram.NewTelegram(store, ch)

	go tg.Run("6538437322:AAGjNbhTJedoQO9I2Em0J3KqbQEFKK_2g78")

	handler := rest.NewHandler(ch)

	http.ListenAndServe("localhost:8080", handler.Init())
}
