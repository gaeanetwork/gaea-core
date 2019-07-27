package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gaeanetwork/gaea-core/did/address"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/register/{net}/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		netName := vars["net"]
		privKey, addr, err := address.Register(netName)
		if err != nil {
			log.Println(err)
		}
		doc := map[string]interface{}{
			"name":    vars["name"],
			"id":      addr,
			"privkey": privKey,
		}
		json.NewEncoder(w).Encode(doc)
	})

	r.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		value := map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&value)
		doc := map[string]interface{}{
			"result": true,
		}
		json.NewEncoder(w).Encode(doc)
	})
	http.Handle("/", r)
	http.ListenAndServe(":8082", r)
}
