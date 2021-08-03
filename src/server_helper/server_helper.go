package server_helper

import (
	"db_helper"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"logger_helper"
	"net/http"
	"telegram_helper"

	"github.com/gorilla/mux"
)

type api struct {
	router http.Handler
}

type Server interface {
	Router() http.Handler
}

type User struct {
	Username string `json:"username"`
	UserId int `json:"user_id"`
}

var (
	db string = ""
	table string = ""
	token string = ""
)

const (
	GET_USERS string = "/getUsers"
	ADD_USER string = "/addUser"
	DEL_USER string = "/deleteUser"
)


func (a *api) getUsers(w http.ResponseWriter, r *http.Request) {
	logger_helper.LogInfo("Getting users...")
	usernames, user_ids := db_helper.GetUsers(db, table)
	var users []User
	for i, _ := range usernames {
		users = append(users, User{usernames[i], user_ids[i]})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}


func (a *api) addUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(body, &user)
	logger_helper.LogInfo(fmt.Sprintf("Adding user (%s, %d)",
									  user.Username, user.UserId))
	db_helper.AddUser(db, table, user.Username, user.UserId)
	telegram_helper.InitBot(token, user.UserId)
	json.NewEncoder(w).Encode(user)
}


func (a *api) deleteUser(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var user User
	logger_helper.LogInfo(fmt.Sprintf("Deleting user (%s, %d)",
									  user.Username, user.UserId))
	json.Unmarshal(body, &user)
	db_helper.DeleteUser(db, table, user.UserId)
	telegram_helper.DeleteBot(user.UserId)
	json.NewEncoder(w).Encode(user)
}


func New() Server {
	a := &api{}
	r:= mux.NewRouter()
	r.HandleFunc(GET_USERS, a.getUsers).Methods(http.MethodGet)
	r.HandleFunc(ADD_USER, a.addUser).Methods(http.MethodPost)
	r.HandleFunc(DEL_USER, a.deleteUser).Methods(http.MethodPost)
	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}


func InitServer(port string, db_name string, table_name string,
				token_id string) {
	db = db_name
	table = table_name
	token = token_id
	s := New()
	logger_helper.LogInfo(fmt.Sprintf("Starting server on port %s...", port))
	logger_helper.LogError(http.ListenAndServe(port, s.Router()).Error())
}


func SetDb(db_name string, table_name string) {
	logger_helper.LogInfo(fmt.Sprintf("Changing DB to %s:%s",
									  db_name, table_name))
	db = db_name
	table = table_name
}
