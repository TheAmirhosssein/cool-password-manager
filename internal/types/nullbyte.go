package types

import (
	"database/sql"
	"database/sql/driver"
)

type NullByte struct {
	Byte  byte `json:"byte"`
	Valid bool `json:"valid"`
}

// NullByte
func (n *NullByte) Set(value byte) {
	n.Valid = true
	n.Byte = value
}

func (n *NullByte) Scan(value interface{}) error {
	var ns sql.NullByte

	if err := ns.Scan(value); err != nil {
		n.Byte, n.Valid = 0, false
		return err
	}

	n.Byte, n.Valid = ns.Byte, ns.Valid
	return nil
}

func (n NullByte) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Byte, nil
}
