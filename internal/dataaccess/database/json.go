package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSON struct {
	Data any
}

func (j *JSON) Scan(src any) error {
	if src == nil {
		return nil
	}

	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, &j.Data)

	case string:
		return json.Unmarshal([]byte(src), &j.Data)

	default:
		return fmt.Errorf("unsupported type for json scan: %T", src)
	}
}

func (j JSON) Value() (driver.Value, error) {
	if j.Data == nil {
		return nil, nil
	}

	return json.Marshal(j.Data)
}
