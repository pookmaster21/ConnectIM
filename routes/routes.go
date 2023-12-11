package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pookmaster21/ConnectIM-Server/types"
)

type Message types.Message

func NewRoute() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	return mux
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		data, e := json.Marshal(Message{
			Status: "Error",
			Msg:    "",
		})
		if e != nil {
			io.WriteString(w, "Error in json")
		}

		io.WriteString(w, string(data))
	}

	w.WriteHeader(http.StatusOK)

	data, e := json.Marshal(Message{
		Status: "Ok",
		Msg:    "Root",
	})
	if e != nil {
		io.WriteString(w, "Error in json")
	}

	io.WriteString(w, string(data))
}

func getHello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)

		data, e := json.Marshal(Message{
			Status: "Error",
			Msg:    "",
		})
		if e != nil {
			io.WriteString(w, "Error in json")
		}

		io.WriteString(w, string(data))
	}

	w.WriteHeader(http.StatusOK)

	data, e := json.Marshal(Message{
		Status: "Ok",
		Msg:    "Hello",
	})
	if e != nil {
		io.WriteString(w, "Error in json")
	}

	io.WriteString(w, string(data))
}
