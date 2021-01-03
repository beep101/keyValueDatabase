package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"../storage"
	"../userManager"
	"../webServer"
	"../wrappers"
)

type configData struct {
	Address         string
	Name            string
	Security        string
	Port            string
	DatabaseCaching int
	UsersCaching    int
}

func CheckDb(address, name string) bool {
	ext := true

	_, err := os.Stat(address + name + storage.TreeExt)
	if os.IsNotExist(err) {
		ext = false
	}
	_, err = os.Stat(address + name + storage.DataExt)
	if os.IsNotExist(err) {
		ext = false
	}

	return ext
}

func main() {
	//parse config file
	config := configData{}
	confFileAddress := os.Args[1]
	file, err := ioutil.ReadFile(confFileAddress)
	if err != nil {
		log.Fatal(errors.New("Error: Config file not found"))
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(errors.New("Error: Config file structure error"))
	}

	port, err := strconv.Atoi(config.Port)
	if err != nil {
		log.Fatal("Error: Port not defined")
	}
	if port < 1024 || port > 49151 {
		log.Fatal("Error: Port is not allowed")
	}

	config.Address = config.Address + "/"
	dbCheck := CheckDb(config.Address, config.Name)
	usersCheck := CheckDb(config.Address, config.Name+"Users")

	if _, err := os.Stat(config.Address); os.IsNotExist(err) {
		err := os.Mkdir(config.Address, 0666)
		if err != nil {
			log.Fatal(errors.New("Error: Cannot create folder"))
		}
	}

	var db *wrappers.StringStringDb
	var udb *userManager.UserManager
	if dbCheck && usersCheck {
		db = wrappers.Open(config.Address, config.Name, config.DatabaseCaching, false)
		udb = userManager.Open(config.Address, config.Name+"Users", config.UsersCaching, false)
		fmt.Println("***Database successefully opened")
	} else if dbCheck && !usersCheck {
		db = wrappers.Open(config.Address, config.Name, config.DatabaseCaching, false)
		udb = userManager.Open(config.Address, config.Name+"Users", config.UsersCaching, true)
		err := udb.AddUser("admin", "admin", "A")
		if err != nil {
			log.Fatal(errors.New("Error: Cannot create admin"))
		}
		fmt.Println("***Database successefully opened, users table created with default settings")
	} else {
		db = wrappers.Open(config.Address, config.Name, config.DatabaseCaching, true)
		udb = userManager.Open(config.Address, config.Name+"Users", config.UsersCaching, true)
		err := udb.AddUser("admin", "admin", "A")
		if err != nil {
			log.Fatal(errors.New("Error: Cannot create admin"))
		}
		fmt.Println("***Database successefully created")
	}

	var pm userManager.PermissionManager
	if config.Security == "OPEN" {
		pm = userManager.CreateOpenPermission()
	} else if config.Security == "RESTRICT" {
		pm = userManager.CreateRestrictPermission(udb)
	} else if config.Security == "AUTH" {
		pm = userManager.CreateAuthPermission(udb)
	} else {
		log.Fatal(errors.New("Error: Bad security settings"))
	}

	//start server
	webServer.Start(config.Port, db, udb, pm)
}
