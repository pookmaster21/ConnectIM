package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pookmaster21/ConnectIM/types"
)

type Message types.Message

func NewRoute() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", getRoot)

	return mux
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		data, e := json.Marshal(map[string]string{
			"msg": "root",
		})
		if e != nil {
			data = []byte("Error")
		}
		io.WriteString(w, string(data))
	}
}
