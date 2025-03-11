package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	Name string `json:"name"`
}

var userCashe = make(map[int]User)

var casheMutex sync.RWMutex

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users/{id}", getUser)
	// mux.HandleFunc("De /users/{id}", getUser)


	fmt.Println("Server listening to :8080")

	http.ListenAndServe(":8080", mux)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	casheMutex.Lock()

	userCashe[len(userCashe ) + 1] = user

	casheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func getUser(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.PathValue("id")) 

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	casheMutex.RLock()

	user, ok := userCashe[id]

	casheMutex.RUnlock()

	if !ok {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	w.WriteHeader("Content-Type":"application/json")
	j, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write(j)

}