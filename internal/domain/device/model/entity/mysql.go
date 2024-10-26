package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type mysqlRecordId int64
type mysqlText string
type mysqlDate time.Time
type mysqlJson map[string]interface{}

func (a *mysqlJson) Scan(value interface{}) error {
	val, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to int64 failed")
	}
	if err := json.Unmarshal(val, a); err != nil {
		return err
	}

	return nil
}

func (a mysqlJson) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (m *mysqlJson) Set(key string, value interface{}) {
	(*m)[key] = value
}

func (m *mysqlJson) Get(key string) (interface{}, error) {
	if value, exists := (*m)[key]; exists {
		return value, nil
	}
	return nil, errors.New("key not found")
}

func (m *mysqlJson) Delete(key string) {
	delete(*m, key)
}

func (m *mysqlJson) Map() map[string]interface{} {
	return *m
}

func (m *mysqlJson) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), m)
}

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
