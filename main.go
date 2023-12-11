package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pookmaster21/ConnectIM-Server/routes"
)

func main() {
	server := routes.NewRoute()

	e := http.ListenAndServe(":8080", server)

	if errors.Is(e, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else {
		fmt.Println(e)
	}
}
