package null

import (
	"database/sql"
	"encoding/json"
)

type String sql.NullString

func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (ns *String) Scan(value interface{}) error {
	if value == nil {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	ns.String = value.(string)
	return nil
}
