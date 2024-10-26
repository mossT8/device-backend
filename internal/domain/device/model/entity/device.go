package entity

import (
	"time"
)

type Device struct {
	ID mysqlRecordId

	AccountId    mysqlRecordId
	Name         mysqlText
	SerialNumber mysqlText
	ModelId      mysqlRecordId
	ModelConfig  mysqlJson

	CreatedAt  mysqlDate
	ModifiedAt mysqlDate
}

func NewDevice(accountId, modelID int64, name, serialNumber string, config map[string]interface{}) Device {
	return Device{
		AccountId:    mysqlRecordId(accountId),
		ModelId:      mysqlRecordId(modelID),
		Name:         mysqlText(name),
		SerialNumber: mysqlText(serialNumber),
		ModelConfig:  mysqlJson(config),
		CreatedAt:    mysqlDate(time.Now()),
		ModifiedAt:   mysqlDate(time.Now()),
	}
}

func (d *Device) GetID() int64 {
	return int64(d.ID)
}

func (d *Device) GetAccountId() int64 {
	return int64(d.AccountId)
}

func (d *Device) GetName() string {
	return string(d.Name)
}

func (d *Device) GetSerialNumber() string {
	return string(d.SerialNumber)
}

func (d *Device) GetModelId() int64 {
	return int64(d.ModelId)
}

func (d *Device) GetModelConfig() map[string]interface{} {
	return d.ModelConfig.Map()
}

func (d *Device) GetCreatedAt() time.Time {
	return time.Time(d.CreatedAt)
}

func (d *Device) GetModifiedAt() time.Time {
	return time.Time(d.ModifiedAt)
}

func (d *Device) SetID(id int64) {
	d.ID = mysqlRecordId(id)
}

func (d *Device) SetAccountId(accountId int64) {
	d.AccountId = mysqlRecordId(accountId)
	d.ModifiedAt = mysqlDate(time.Now())
}

func (d *Device) SetName(name string) {
	d.Name = mysqlText(name)
	d.ModifiedAt = mysqlDate(time.Now())
}

func (d *Device) SetSerialNumber(serialNumber string) {
	d.SerialNumber = mysqlText(serialNumber)
	d.ModifiedAt = mysqlDate(time.Now())
}

func (d *Device) SetModelId(modelId int64) {
	d.ModelId = mysqlRecordId(modelId)
	d.ModifiedAt = mysqlDate(time.Now())
}

func (d *Device) SetModelConfig(modelConfig map[string]interface{}) {
	d.ModelConfig = mysqlJson(modelConfig)
	d.ModifiedAt = mysqlDate(time.Now())
}

func (d *Device) SetCreatedAt(createdAt time.Time) {
	d.CreatedAt = mysqlDate(createdAt)
}

func (d *Device) SetModifiedAt(modifiedAt time.Time) {
	d.ModifiedAt = mysqlDate(modifiedAt)
}
