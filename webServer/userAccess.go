package webServer

import (
	"encoding/json"
	"io"
)

type userBody struct {
	Uname    string
	NewUname string
	Passwd   string
	Role     string
}

func (app *application) getUser(uname, passwd string, body io.Reader) (int, string) {
	role := app.getRole(uname, passwd)
	var req userBody
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		return 400, ""
	}
	if role == "A" {
		role := app.um.GetUser(req.Uname)
		if role == "" {
			return 404, ""
		}
		return 200, role
	}
	return 401, ""
}

func (app *application) addUser(uname, passwd string, body io.Reader) (int, string) {
	role := app.getRole(uname, passwd)
	var req userBody
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		return 400, ""
	}
	if role == "A" {
		err := app.um.AddUser(req.Uname, req.Passwd, req.Role)
		if err != nil {
			return 400, ""
		} else {
			return 200, ""
		}
	} else {
		return 401, ""
	}
}

func (app *application) modUser(uname, passwd string, body io.Reader) (int, string) {
	role := app.getRole(uname, passwd)
	var req userBody
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		return 400, ""
	}
	if role == "A" {
		err := app.um.ModUser(req.Uname, req.NewUname, req.Passwd, req.Role)
		if err != nil {
			return 400, ""
		}
		return 200, ""
	} else if role == "R" || role == "W" {
		if req.Role != "" && role != req.Role {
			return 401, ""
		}
		err := app.um.ModUser(uname, req.NewUname, req.Passwd, req.Role)
		if err != nil {
			return 400, ""
		}
		return 200, ""
	} else {
		return 401, ""
	}
}

func (app *application) delUser(uname, passwd string, body io.Reader) (int, string) {
	role := app.getRole(uname, passwd)
	var req userBody
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		return 400, ""
	}
	if role == "A" {
		err := app.um.DelUser(req.Uname)
		if err != nil {
			return 400, ""
		} else {
			return 200, ""
		}
	} else {
		return 401, ""
	}
}

func (app *application) getRole(uname, passwd string) string {
	return app.um.IsValid(uname, passwd)
}
