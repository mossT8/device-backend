package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type mysqlRecordId int64
type mysqlText string
type mysqlBool bool
type mysqlDate time.Time

func (a *mysqlRecordId) Scan(value interface{}) error {
	val, ok := value.(int64)
	if !ok {
		return errors.New("type assertion to int64 failed")
	}
	*a = mysqlRecordId(val)
	return nil
}

func (a mysqlRecordId) Value() (driver.Value, error) {
	return int64(a), nil
}

func (a *mysqlText) Scan(value interface{}) error {
	if value == nil {
		*a = ""
		return nil
	}
	val, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	*a = mysqlText(string(val))
	return nil
}

func (a mysqlText) Value() (driver.Value, error) {
	return string(a), nil
}

func (a *mysqlBool) Scan(value interface{}) error {
	if value == nil {
		*a = false
		return nil
	}

	switch v := value.(type) {
	case bool:
		*a = mysqlBool(v)
	case int64:
		*a = mysqlBool(v != 0)
	case string:
		if v == "true" {
			*a = mysqlBool(true)
		} else {
			*a = mysqlBool(false)
		}
	default:
		return errors.New("type assertion to bool failed")
	}
	return nil
}

func (a mysqlBool) Value() (driver.Value, error) {
	return bool(a), nil
}

func (a *mysqlDate) Scan(value interface{}) error {
	if value == nil {
		*a = mysqlDate(time.Time{})
		return nil
	}
	val, ok := value.(time.Time)
	if !ok {
		return errors.New("type assertion to time.Time failed")
	}
	*a = mysqlDate(val)
	return nil
}

func (a mysqlDate) Value() (driver.Value, error) {
	return time.Time(a), nil
}

func (a mysqlDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(a))
}

func (a *mysqlDate) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	*a = mysqlDate(t)
	return nil
}
