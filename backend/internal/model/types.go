package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONStrings JSON 字符串数组类型
type JSONStrings []string

func (j JSONStrings) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONStrings) Scan(value interface{}) error {
	if value == nil {
		*j = JSONStrings{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan source for JSONStrings")
	}
	return json.Unmarshal(bytes, j)
}
