package api

import (
	"encoding/json"
	"github.com/antaresvision/helloserver/db"
	"github.com/antaresvision/helloserver/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (srv *server) GetAllItems(writer http.ResponseWriter, request *http.Request) {
	items, err := srv.ds.GetAllItems()
	if err != nil {
		http.Error(writer, "could not retrieve items", http.StatusInternalServerError)
		return
	}

	writer.Header().Set(`Content-Type`, `application/json`)
	enc := json.NewEncoder(writer)
	err = enc.Encode(items)
	if err != nil {
		http.Error(writer, "could not encode items in JSON", http.StatusInternalServerError)
		return
	}
}

func SaveItem(writer http.ResponseWriter, request *http.Request, ds *db.Store) {
	item := models.Item{}

	dec := json.NewDecoder(request.Body)
	err := dec.Decode(&item)
	if err != nil {
		http.Error(writer, "could not parse item id", http.StatusBadRequest)
		return
	}

	if item.Id == 0 {
		item, err = ds.CreateItem(item.NtinId, item.Serial, item.Status)
	} else {
		err = ds.UpdateItem(item)
	}
	if err != nil {
		http.Error(writer, "could not save item", http.StatusInternalServerError)
		return
	}

	writer.Header().Set(`Content-Type`, `application/json`)
	enc := json.NewEncoder(writer)
	err = enc.Encode(item)
	if err != nil {
		http.Error(writer, "could not encode items in JSON", http.StatusInternalServerError)
		return
	}
}

func RemoveItemById(writer http.ResponseWriter, request *http.Request, ds *db.Store) {
	itemIdVar := mux.Vars(request)["id"]
	itemId, err := strconv.Atoi(itemIdVar)
	if err != nil {
		http.Error(writer, "could not parse item id", http.StatusBadRequest)
		return
	}

	err = ds.RemoveItem(itemId)
	if err != nil {
		http.Error(writer, "could not remove item", http.StatusInternalServerError)
		return
	}
}

func GetItemById(writer http.ResponseWriter, request *http.Request, ds *db.Store) {
	itemIdVar := mux.Vars(request)["id"]
	itemId, err := strconv.Atoi(itemIdVar)
	if err != nil {
		http.Error(writer, "could not parse item id", http.StatusBadRequest)
		return
	}

	item, err := ds.GetItem(itemId)
	if err != nil {
		http.Error(writer, "could not retrieve item", http.StatusInternalServerError)
		return
	}

	writer.Header().Set(`Content-Type`, `application/json`)
	enc := json.NewEncoder(writer)
	err = enc.Encode(item)
	if err != nil {
		http.Error(writer, "could not encode items in JSON", http.StatusInternalServerError)
		return
	}
}
