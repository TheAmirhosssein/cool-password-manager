package types

import (
	"database/sql"
	"database/sql/driver"
)

type NullString struct {
	String string `json:"string"`
	Valid  bool   `json:"valid"`
}

// NullString
func (n *NullString) Set(value string) {
	n.Valid = true
	n.String = value
}

func (n *NullString) Scan(value interface{}) error {
	var ns sql.NullString

	if err := ns.Scan(value); err != nil {
		n.String, n.Valid = "", false
		return err
	}

	n.String, n.Valid = ns.String, ns.Valid
	return nil
}

func (n NullString) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String, nil
}
