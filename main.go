package main

import (
	"github.com/antaresvision/helloserver/api"
	"github.com/antaresvision/helloserver/db"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

func main() {

	sessions := map[string]string{}

	dataStore := db.NewConnection()
	defer dataStore.Close()

	srv := api.NewServer(dataStore)

	r := mux.NewRouter()

	r.Path("/login").Methods(http.MethodPost).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user := request.PostFormValue("user")
		password := request.PostFormValue("password")

		if user == "foo" && password == "bar" {
			encryptedUser := EncryptToBase64(user)
			http.SetCookie(writer, &http.Cookie{
				Name:       "myauth",
				Value:      encryptedUser,
			})
		} else {
			http.Error(writer, "invalid username/password", http.StatusBadRequest)
		}
	})

	r.Path("/login/token").Methods(http.MethodPost).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user, password, ok := request.BasicAuth()
		if !ok {
			http.Error(writer, "invalid username/password", http.StatusBadRequest)
			return
		}

		if user == "foo" && password == "bar" {
			key := uuid.New().String()
			sessions[key]=user
			writer.Write([]byte(key))
		} else {
			http.Error(writer, "invalid username/password", http.StatusBadRequest)
		}
	})

	r.Path("/greetings").Methods(http.MethodGet).HandlerFunc(api.GreetingsHandler)
	r.Path("/greetings/{name}").Methods(http.MethodGet).HandlerFunc(api.GreetingsHandler)

	itemsRouter := r.PathPrefix("/items").Subrouter()
	itemsRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("prima dell'handler")

			var user = ""

			authCookie, err := r.Cookie("myauth")
			if err == nil {
				user = DecryptFromBase64(authCookie.Value)
			}

			if user == "" {
				reqToken := r.Header.Get("Authorization")
				splitToken := strings.Split(reqToken, "Bearer ")
				reqToken = splitToken[1]

				user = sessions[reqToken]
			}

			log.Println("User:", user)
			var auth = (user == "foo")

			if auth {
				next.ServeHTTP(w, r)
			} else {
				log.Println("Not authorized")
				http.Error(w, "Auth required", http.StatusUnauthorized)
			}
			log.Println("dopo l'handler")
		})
	})


	itemsRouter.Path("/{id}").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Get item handler")
		api.GetItemById(writer, request, dataStore)
	})
	itemsRouter.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Remove item handler")
		api.RemoveItemById(writer, request, dataStore)
	})
	itemsRouter.Path("/{id}").Methods(http.MethodPost).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Save item handler")
		api.SaveItem(writer, request, dataStore)
	})
	itemsRouter.Path("/").Methods(http.MethodGet).HandlerFunc(srv.GetAllItems)


	http.ListenAndServe(":8000", r)
}
