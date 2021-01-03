package webServer

import (
	"log"
	"net/http"

	userManager "../userManager"
	wrappers "../wrappers"
)

type application struct {
	db *wrappers.StringStringDb
	um *userManager.UserManager
	pm userManager.PermissionManager
}

func Start(port string, database *wrappers.StringStringDb, usermanager *userManager.UserManager, permissionmanager userManager.PermissionManager) {
	app := &application{db: database, um: usermanager, pm: permissionmanager}
	http.HandleFunc("/", app.databaseHandler)
	http.HandleFunc("/user", app.userHandler)
	http.HandleFunc("/node", app.nodeHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (app *application) databaseHandler(w http.ResponseWriter, r *http.Request) {
	uname, passwd, _ := r.BasicAuth()
	var status int
	var response string
	switch r.Method {
	case "GET":
		status, response = app.getData(uname, passwd, r.Body)
	case "POST":
		status, response = app.addData(uname, passwd, r.Body)
	case "DELETE":
		status, response = app.delData(uname, passwd, r.Body)
	default:
		status, response = 400, ""
	}
	w.WriteHeader(status)
	if response != "" {
		w.Write([]byte(response))
	}
	return
}

func (app *application) userHandler(w http.ResponseWriter, r *http.Request) {
	uname, passwd, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(401)
		return
	}
	var status int
	var response string
	switch r.Method {
	case "GET":
		status, response = app.getUser(uname, passwd, r.Body)
	case "POST":
		status, response = app.addUser(uname, passwd, r.Body)
	case "PUT":
		status, response = app.modUser(uname, passwd, r.Body)
	case "DELETE":
		status, response = app.delUser(uname, passwd, r.Body)
	default:
		status, response = 400, ""
	}
	w.WriteHeader(status)
	if response != "" {
		w.Write([]byte(response))
	}
	return
}

func (app *application) nodeHandler(w http.ResponseWriter, r *http.Request) {
	//to be discussed
}
