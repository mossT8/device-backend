package entity

import "time"

type Models struct {
	ID mysqlRecordId

	Name mysqlText
	Code mysqlText

	CreatedAt  mysqlDate
	ModifiedAt mysqlDate
}

func NewModels(name, code string) Models {
	return Models{
		Name:       mysqlText(name),
		Code:       mysqlText(code),
		CreatedAt:  mysqlDate(time.Now()),
		ModifiedAt: mysqlDate(time.Now()),
	}
}

func (m *Models) GetID() int64 {
	return int64(m.ID)
}

func (m *Models) GetName() string {
	return string(m.Name)
}

func (m *Models) GetCode() string {
	return string(m.Code)
}

func (m *Models) GetCreatedAt() time.Time {
	return time.Time(m.CreatedAt)
}

func (m *Models) GetModifiedAt() time.Time {
	return time.Time(m.ModifiedAt)
}

func (m *Models) SetName(name string) {
	m.Name = mysqlText(name)
}

func (m *Models) SetCode(code string) {
	m.Code = mysqlText(code)
}

func (m *Models) SetCreatedAt(createdAt time.Time) {
	m.CreatedAt = mysqlDate(createdAt)
}

func (m *Models) SetModifiedAt(modifiedAt time.Time) {
	m.ModifiedAt = mysqlDate(modifiedAt)
}

func (m *Models) SetID(id int64) {

	m.ID = mysqlRecordId(id)
}
