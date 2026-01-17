package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
)

func (r *Role) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*r = Role(v)
		return nil
	case []byte:
		*r = Role(string(v))
		return nil
	default:
		return errors.New("role should be a string or []byte")
	}
}

func (r Role) Value() (driver.Value, error) {
	switch r {
	case RoleAdmin, RoleUser, RoleModerator:
		return string(r), nil
	default:
		return nil, fmt.Errorf("invalid role: %s", r)
	}
}
