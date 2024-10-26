package entity

import (
	"time"
)

type Sensor struct {
	ID mysqlRecordId

	UnitId         mysqlRecordId
	Code           mysqlText
	Name           mysqlText
	ConfigRequried mysqlJson
	DefaultConfig  mysqlJson

	CreatedAt  mysqlDate
	ModifiedAt mysqlDate
}

func NewSensor(code, name string) Sensor {
	return Sensor{
		Code: mysqlText(code),
		Name: mysqlText(name),
	}
}

func (s *Sensor) GetID() int64 {
	return int64(s.ID)
}

func (s *Sensor) GetUnit() int64 {
	return int64(s.UnitId)
}

func (s *Sensor) GetCode() string {
	return string(s.Code)
}

func (s *Sensor) GetName() string {
	return string(s.Name)
}

func (s *Sensor) GetConfigRequried() map[string]interface{} {
	return s.ConfigRequried.Map()
}

func (s *Sensor) GetDefaultConfig() map[string]interface{} {
	return s.DefaultConfig.Map()
}

func (s *Sensor) GetCreatedAt() time.Time {
	return time.Time(s.CreatedAt)
}

func (s *Sensor) GetModifiedAt() time.Time {
	return time.Time(s.ModifiedAt)
}

func (s *Sensor) SetID(id int64) {
	s.ID = mysqlRecordId(id)
}

func (s *Sensor) SetUnitId(unitId int64) {
	s.UnitId = mysqlRecordId(unitId)
	s.ModifiedAt = mysqlDate(time.Now())
}

func (s *Sensor) SetCode(code string) {
	s.Code = mysqlText(code)
	s.ModifiedAt = mysqlDate(time.Now())
}

func (s *Sensor) SetName(name string) {
	s.Name = mysqlText(name)
	s.ModifiedAt = mysqlDate(time.Now())
}

func (s *Sensor) SetConfigRequried(configRequried map[string]interface{}) {
	s.ConfigRequried = mysqlJson(configRequried)
	s.ModifiedAt = mysqlDate(time.Now())
}

func (s *Sensor) SetDefaultConfig(defaultConfig map[string]interface{}) {
	s.DefaultConfig = mysqlJson(defaultConfig)
	s.ModifiedAt = mysqlDate(time.Now())
}

func (s *Sensor) SetCreatedAt(createdAt time.Time) {
	s.CreatedAt = mysqlDate(createdAt)
}

func (s *Sensor) SetModifiedAt(modifiedAt time.Time) {
	s.ModifiedAt = mysqlDate(modifiedAt)
}
