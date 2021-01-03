package userManager

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (um *UserManager) IsValid(uname, passwd string) string {
	role, saltHash, err := um.getData(uname)
	if err != nil {
		return ""
	}
	if compare(saltHash, passwd) {
		return role
	}
	return ""
}

func (um *UserManager) GetUser(uname string) string {
	role, _, err := um.getData(uname)
	if err != nil {
		return ""
	}
	return role
}

func (um *UserManager) AddUser(name, passwd, role string) error {
	if role != "A" && role != "W" && role != "R" {
		return errors.New("Error: Undefined role")
	}
	passwd = saltAndHash(passwd)
	return um.addUser(name, passwd, role, false)
}

func (um *UserManager) ModUser(uname, newUname, newPass, newRole string) error {
	if newRole != "A" && newRole != "W" && newRole != "R" && newRole != "" {
		return errors.New("Error: Undefined role")
	}
	if newPass != "" {
		newPass = saltAndHash(newPass)
	}
	return um.modUser(uname, newUname, newPass, newRole)
}

func (um *UserManager) DelUser(uname string) error {
	return um.delUser(uname)
}

//password related functions
func saltAndHash(passwd string) string {
	hash, err := hash([]byte(passwd))
	if err != nil {
		return ""
	}
	return string(hash)
}

func hash(pass []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
}

func compare(saltHash, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(saltHash), []byte(pass))
	if err != nil {
		return false
	}
	return true
}
