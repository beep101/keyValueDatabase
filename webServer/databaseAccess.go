package webServer

import (
	"encoding/json"
	"io"
)

type dataBody struct {
	Key   string
	Value string
}

func (app *application) getData(uname, passwd string, body io.Reader) (int, string) {
	var req dataBody
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		return 400, ""
	}
	if app.pm.IsReader(uname, passwd) {
		result := app.db.Get(req.Key)
		if result == "" {
			return 404, ""
		}
		return 200, result
	}
	return 401, ""
}

func (app *application) addData(uname, passwd string, body io.Reader) (int, string) {
	var req dataBody
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		return 400, ""
	}
	if app.pm.IsWriter(uname, passwd) {
		err := app.db.Add(req.Key, req.Value)
		if err != nil {
			return 400, ""
		}
		return 200, ""
	}
	return 401, ""
}

func (app *application) delData(uname, passwd string, body io.Reader) (int, string) {
	var req dataBody
	err := json.NewDecoder(body).Decode(&req)
	if err != nil {
		return 400, ""
	}
	if app.pm.IsWriter(uname, passwd) {
		err := app.db.Del(req.Key)
		if err != nil {
			return 400, ""
		}
		return 200, ""
	}
	return 401, ""
}
