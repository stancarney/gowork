package gowork

import "time"

type EventCode string

const (
	ACCESS EventCode = "ACCESS"
	CREATE EventCode = "CREATE"
	DELETE EventCode = "DELETE"
	DENIED EventCode = "DENIED"
	ERROR EventCode = "ERROR"
	READ EventCode = "READ"
	UPDATE EventCode = "UPDATE"

	PERM_LOGIN Permission = "LOGIN"
	PERM_SUPER_USER Permission = "SUPER_USER" //Has all permissions except LOGIN
)

type Session struct {
	Id         string             `json:"id"`
	Created    time.Time          `json:"created" validate:"nonzero" datastore:"cr"`
	LastAccess time.Time          `json:"access" validate:"nonzero" datastore:"access"`
	Values     map[string]string  `json:"values" validate:"nonzero"` //TODO:Stan investigate changing value to interface or []byte
	UserId     string             `json:"userid" validate:"nonzero"`
	Version    int                `json:"v" validate:"min=0" datastore:"v"`
}

type Config struct {
	Id          string    `json:"id" datastore:"id"`
	Created     time.Time `json:"created,omitempty" validate:"nonzero" datastore:"cr"`
	Description string    `json:"desc,omitempty" validate:"nonzero" datastore:"descr"`
	Value       string    `json:"value,omitempty" validate:"nonzero"`
	Version     int       `json:"v,omitempty" validate:"min=0" datastore:"v"`
}

type User interface {
	HasPermission(p interface{}) bool
}

type Permission string
type Permissions []Permission

func (u *Permissions) HasPermission(perm Permission) bool {
	for _, p := range *u {
		if p == perm {
			return true
		}

		//SUPER_USER's can still have their LOGIN permission revoked
		if perm == PERM_LOGIN {
			continue;
		}

		if p == PERM_SUPER_USER {
			return true
		}
	}
	return false
}
