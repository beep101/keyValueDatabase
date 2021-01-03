package userManager

import (
	"errors"

	ssdb "../wrappers"
)

type UserManager struct {
	db *ssdb.StringStringDb
}

func Open(address, name string, cacheCount int, rewrite bool) *UserManager {
	db := ssdb.Open(address, name, cacheCount, rewrite)
	return &UserManager{db: db}
}

func (udb *UserManager) Close() error {
	return udb.Close()
}

func (udb *UserManager) addUser(uname, passwd, role string, mod bool) error {
	if mod {
		return udb.db.Add(uname, role+passwd)
	}
	if _, _, err := udb.getData(uname); err == nil {
		return errors.New("Error: Cannot rewrite user")
	}
	return udb.db.Add(uname, role+passwd)
}

func (udb *UserManager) delUser(uname string) error {
	return udb.db.Del(uname)
}

func (udb *UserManager) modUser(uname, newUname, newPasswd, newRole string) error {
	role, passwd, err := udb.getData(uname)
	if err != nil {
		return err
	}
	if newPasswd == "" {
		newPasswd = passwd
	}
	if newRole == "" {
		newRole = role
	}
	if newUname != "" {
		udb.delUser(uname)
		udb.addUser(newUname, newPasswd, newRole, false)
		return nil
	} else {
		udb.addUser(uname, newPasswd, newRole, true)
		return nil
	}
}

func (udb *UserManager) getData(uname string) (string, string, error) {
	userData := udb.db.Get(uname)
	if userData == "" {
		return "", "", errors.New("Error: User not found")
	}
	return string(userData[0]), userData[1:], nil
}
