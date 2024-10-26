package device

import (
	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/domain/device/model/entity"
	"mossT8.github.com/device-backend/internal/domain/device/model/request"
	"mossT8.github.com/device-backend/internal/infrastructure/logger"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
)

type DeviceDomain interface {
	AddDevice(requestID string, accountID int64, payload request.Device) (*entity.Device, error)
	UpdateDevice(requestID string, accountID, deviceID int64, payload request.Device) (*entity.Device, error)
	FetchDevice(requestID string, accountID, deviceID int64) (*entity.Device, error)
	ListDevices(requestID string, accountID, page, pageSize int64) ([]entity.Device, *int64, error)
	DeleteDevice(requestID string, accountID, deviceID int64) error

	FetchSensor(requestID string, sensorID int64) (*entity.Sensor, error)
	ListSensors(requestID string, page, pageSize int64) ([]entity.Sensor, *int64, error)

	FetchUnit(requestID string, unitID int64) (*entity.Units, error)
	ListUnits(requestID string, page, pageSize int64) ([]entity.Units, *int64, error)

	FetchModel(requestID string, modelID int64) (*entity.Models, error)
	ListModels(requestID string, page, pageSize int64) ([]entity.Models, *int64, error)
}

type DeviceDomainImpl struct {
	dbConn *datastore.MySqlDataStore
}

func NewDeviceDomain(conn *datastore.MySqlDataStore) DeviceDomain {
	return &DeviceDomainImpl{
		dbConn: conn,
	}
}

// Device methods
func (d *DeviceDomainImpl) AddDevice(requestID string, accountID int64, payload request.Device) (*entity.Device, error) {
	device := entity.NewDevice(accountID, payload.ModelId, payload.Name, payload.SerialNumber, payload.ModelConfig)

	if err := device.AddDevice(*d.dbConn); err != nil {
		logger.Errorf(requestID, "unable to create device %+v", device)
		return nil, err
	}

	return &device, nil
}

func (d *DeviceDomainImpl) UpdateDevice(requestID string, accountID, deviceID int64, payload request.Device) (*entity.Device, error) {
	device := &entity.Device{}
	device.SetID(deviceID)
	if err := device.GetDeviceByID(*d.dbConn); err != nil {
		logger.Errorf(requestID, LogCantGetDeviceByID, accountID, deviceID)
		return nil, err
	}

	if device.GetAccountId() != accountID {
		logger.Errorf(requestID, LogCantViewDeviceByID, deviceID)
		return nil, domain.ErrNotOwnedDeviceByID
	}

	device.SetName(payload.Name)
	device.SetModelConfig(payload.ModelConfig)

	if err := device.UpdateDevice(*d.dbConn, nil); err != nil {
		logger.Errorf(requestID, "unable to update device %+v", device)
		return nil, err
	}

	return device, nil
}

func (d *DeviceDomainImpl) FetchDevice(requestID string, accountID, deviceID int64) (*entity.Device, error) {
	device := &entity.Device{}
	device.SetID(deviceID)
	if err := device.GetDeviceByID(*d.dbConn); err != nil {
		logger.Errorf(requestID, LogCantGetDeviceByID, accountID, deviceID)
		return nil, err
	}

	if device.GetAccountId() != accountID {
		logger.Errorf(requestID, LogCantViewDeviceByID, deviceID)
		return nil, domain.ErrNotOwnedDeviceByID
	}

	return device, nil
}

func (d *DeviceDomainImpl) DeleteDevice(requestID string, accountID, deviceID int64) error {
	device := &entity.Device{}
	device.SetID(deviceID)
	if err := device.GetDeviceByID(*d.dbConn); err != nil {
		logger.Errorf(requestID, LogCantGetDeviceByID, accountID, deviceID)
		return err
	}

	if device.GetAccountId() != accountID {
		logger.Errorf(requestID, LogCantViewDeviceByID, deviceID)
		return domain.ErrNotOwnedDeviceByID
	}

	if err := device.DeleteDevice(*d.dbConn, nil); err != nil {
		logger.Errorf(requestID, "unable to delete device by ID %d", deviceID)
		return err
	}
	return nil
}

func (d *DeviceDomainImpl) ListDevices(requestID string, accountID, page, pageSize int64) ([]entity.Device, *int64, error) {
	queryDevice := entity.Device{}
	queryDevice.SetAccountId(accountID)
	devices, err := queryDevice.ListDevices(*d.dbConn, page, pageSize)
	if err != nil {
		logger.Errorf(requestID, "unable to list devices for account ID %d", accountID)
		return nil, nil, err
	}

	total, err := queryDevice.CountDevices(*d.dbConn)
	if err != nil {
		logger.Errorf(requestID, "unable to count all devices for account ID %d", accountID)
		return nil, nil, err
	}

	return devices, total, nil
}

// Sensor methods
func (d *DeviceDomainImpl) FetchSensor(requestID string, sensorID int64) (*entity.Sensor, error) {
	sensor := &entity.Sensor{}
	sensor.SetID(sensorID)
	if err := sensor.GetSensorByID(*d.dbConn); err != nil {
		logger.Errorf(requestID, "unable to get sensor by ID %d", sensorID)
		return nil, err
	}
	return sensor, nil
}

func (d *DeviceDomainImpl) ListSensors(requestID string, page, pageSize int64) ([]entity.Sensor, *int64, error) {
	querySensor := entity.Sensor{}
	sensors, err := querySensor.ListSensors(*d.dbConn, page, pageSize)
	if err != nil {
		logger.Errorf(requestID, "unable to list sensors")
		return nil, nil, err
	}

	total, err := querySensor.CountSensors(*d.dbConn)
	if err != nil {
		logger.Errorf(requestID, "unable to count all sensors")
		return nil, nil, err
	}

	return sensors, total, nil
}

// Units methods
func (d *DeviceDomainImpl) FetchUnit(requestID string, unitID int64) (*entity.Units, error) {
	unit := &entity.Units{}
	unit.SetID(unitID)
	if err := unit.GetUnitByID(*d.dbConn); err != nil {
		logger.Errorf(requestID, "unable to get unit by ID %d", unitID)
		return nil, err
	}
	return unit, nil
}

func (d *DeviceDomainImpl) ListUnits(requestID string, page, pageSize int64) ([]entity.Units, *int64, error) {
	queryUnit := entity.Units{}
	units, err := queryUnit.ListUnits(*d.dbConn, page, pageSize)
	if err != nil {
		logger.Errorf(requestID, "unable to list units")
		return nil, nil, err
	}

	total, err := queryUnit.CountUnits(*d.dbConn)
	if err != nil {
		logger.Errorf(requestID, "unable to count all units")
		return nil, nil, err
	}

	return units, total, nil
}

// Model methods
func (d *DeviceDomainImpl) FetchModel(requestID string, modelID int64) (*entity.Models, error) {
	model := &entity.Models{}
	model.SetID(modelID)
	if err := model.GetModelByID(*d.dbConn); err != nil {
		logger.Errorf(requestID, "unable to get model by ID %d", modelID)
		return nil, err
	}
	return model, nil
}

func (d *DeviceDomainImpl) ListModels(requestID string, page, pageSize int64) ([]entity.Models, *int64, error) {
	queryModel := entity.Models{}
	models, err := queryModel.ListModels(*d.dbConn, page, pageSize)
	if err != nil {
		logger.Errorf(requestID, "unable to list models")
		return nil, nil, err
	}

	total, err := queryModel.CountModels(*d.dbConn)
	if err != nil {
		logger.Errorf(requestID, "unable to count all models")
		return nil, nil, err
	}

	return models, total, nil
}

var LogCantGetDeviceByID = "unable to get device by ID %d"
var LogCantViewDeviceByID = "account ID %d mismatch for device ID %d"
