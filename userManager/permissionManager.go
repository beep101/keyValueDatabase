package userManager

type PermissionManager interface {
	IsReader(uname, passwd string) bool
	IsWriter(uname, passwd string) bool
}

//OPEN - anyone can read and write
//RESTRICT - anyone can read, only authorized users can write
//AUTH - only authorized users can read and write

type OpenPermissions struct {
}

func CreateOpenPermission() *OpenPermissions {
	return &OpenPermissions{}
}

func (op *OpenPermissions) IsReader(uname, passwd string) bool {
	return true
}

func (op *OpenPermissions) IsWriter(uname, passwd string) bool {
	return true
}

type RestrictPermission struct {
	um *UserManager
}

func CreateRestrictPermission(um *UserManager) *RestrictPermission {
	return &RestrictPermission{um: um}
}

func (rp *RestrictPermission) IsReader(uname, passwd string) bool {
	return true
}

func (rp *RestrictPermission) IsWriter(uname, passwd string) bool {
	role := rp.um.IsValid(uname, passwd)
	if role == "W" || role == "A" {
		return true
	}
	return false
}

type AuthPermission struct {
	um *UserManager
}

func CreateAuthPermission(um *UserManager) *AuthPermission {
	return &AuthPermission{um: um}
}

func (ap *AuthPermission) IsReader(uname, passwd string) bool {
	role := ap.um.IsValid(uname, passwd)
	if role == "R" || role == "W" || role == "A" {
		return true
	}
	return false
}

func (ap *AuthPermission) IsWriter(uname, passwd string) bool {
	role := ap.um.IsValid(uname, passwd)
	if role == "W" || role == "A" {
		return true
	}
	return false
}
