package api

import (
	"encoding/json"
	// "fmt"
	"github.com/raymondvooo/doggy-date-app/server/postgres"
	"net/http"
)

type Email struct {
	Email string `json:"email"`
}

//CheckEmailExists checks against the database to see if email exists in the system
func CheckEmailExists(w http.ResponseWriter, req *http.Request, db *postgres.Db) {
	var e Email
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if e.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request: Please provide valid email!"))
	} else {
		exists, _ := db.CheckEmailExists(e.Email)
		if !exists {
			// exists, _ := json.Marshal(exists)
			// w.Write(exists)
			w.Write([]byte("false"))
		} else {
			// exists, _ := json.Marshal(exists)
			// w.Write(exists)
			w.Write([]byte("true"))
		}
	}
}
