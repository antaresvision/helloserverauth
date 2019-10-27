package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func GreetingsHandler(writer http.ResponseWriter, request *http.Request) {
	variables := mux.Vars(request)
	name, found := variables["name"]
	if !found {
		name = "Generic User"
	}

	resp := fmt.Sprintf("Hello, Dear %s", name)

	fmt.Println(resp)

	jsonResp := Response{
		Message:   resp,
		TimeStamp: time.Now(),
	}

	jsonBuffer, err := json.Marshal(jsonResp)
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(jsonBuffer)
	if err != nil {
		fmt.Println(err)
	}
}

