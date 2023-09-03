package rest

import (
	"DobroBot/model"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	ch chan (model.Discont)
}

func NewHandler(ch chan (model.Discont)) *Handler {
	return &Handler{
		ch: ch,
	}
}

func (h *Handler) Init() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		if r.Method != http.MethodPost {
			http.Error(w, "Unknown path", http.StatusNotFound)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("cant read from body: %v\n", err)
			http.Error(w, "Some server error", http.StatusInternalServerError)
			return
		}

		discont := model.Discont{}
		err = json.Unmarshal(b, &discont)
		if err != nil {
			log.Printf("cant parse discont: %v\n", err)
			http.Error(w, "Some server error", http.StatusInternalServerError)
			return
		}

		if b != nil {
			h.ch <- discont
		} else {
			http.Error(w, "Empty body", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	return r
}
