package http

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"mossT8.github.com/device-backend/internal/domain/customer"
	"mossT8.github.com/device-backend/internal/domain/device"
	"mossT8.github.com/device-backend/internal/domain/device/model/request"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/constants"
)

type DeviceController struct {
	customerDomain customer.CustomerDomain
	deviceDomain   device.DeviceDomain
}

func NewDeviceController(conn *datastore.MySqlDataStore, server *iris.Application, devDomain device.DeviceDomain, custDomain customer.CustomerDomain) DeviceController {
	dc := DeviceController{
		deviceDomain:   devDomain,
		customerDomain: custDomain,
	}

	server.Post(constants.ApiPrefix+"/account/{accountID:int64}/device", dc.HandlePostDevice)
	server.Put(constants.ApiPrefix+"/account/{accountID:int64}/device/{deviceID:int64}/update", dc.HandlePutDevice)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/device/{deviceID:int64}/fetch", dc.HandleGetDevice)
	server.Get(constants.ApiPrefix+"/account/{accountID:int64}/device/list", dc.HandleGetDevices)

	server.Get(constants.ApiPrefix+"/sensor/{sensorID:int64}/fetch", dc.HandleGetSensor)
	server.Get(constants.ApiPrefix+"/sensor/list", dc.HandleGetSensors)

	server.Get(constants.ApiPrefix+"/unit/{unitID:int64}/fetch", dc.HandleGetUnit)
	server.Get(constants.ApiPrefix+"/unit/list", dc.HandleGetUnits)

	server.Get(constants.ApiPrefix+"/model/{modelID:int64}/fetch", dc.HandleGetModel)
	server.Get(constants.ApiPrefix+"/model/list", dc.HandleGetModels)

	return dc
}

// Device handlers
func (dc *DeviceController) HandlePostDevice(ctx iris.Context) {
	var req request.Device
	requestId := GetRequestID(ctx)

	if err := ctx.ReadJSON(&req); err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	_, err = dc.customerDomain.FetchAccount(requestId, accountID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	device, err := dc.deviceDomain.AddDevice(requestId, accountID, req)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), device, http.StatusCreated, requestId)
}

func (dc *DeviceController) HandlePutDevice(ctx iris.Context) {
	var req request.Device
	requestId := GetRequestID(ctx)

	if err := ctx.ReadJSON(&req); err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	_, err = dc.customerDomain.FetchAccount(requestId, accountID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	deviceID, err := ctx.Params().GetInt64("deviceID")
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	device, err := dc.deviceDomain.UpdateDevice(requestId, accountID, deviceID, req)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), device, http.StatusOK, requestId)
}

func (dc *DeviceController) HandleGetDevice(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	_, err = dc.customerDomain.FetchAccount(requestId, accountID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	deviceID, err := ctx.Params().GetInt64("deviceID")
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	device, err := dc.deviceDomain.FetchDevice(requestId, accountID, deviceID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), device, http.StatusOK, requestId)
}

func (dc *DeviceController) HandleGetDevices(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	pageSize, page, err := GetPageAndPageSize(ctx)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	accountID, err := ctx.Params().GetInt64("accountID")
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	account, err := dc.customerDomain.FetchAccount(requestId, accountID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	paginatedList, total, err := dc.deviceDomain.ListDevices(requestId, account.GetID(), *page, *pageSize)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithList(ctx.ResponseWriter(), paginatedList, *page, *pageSize, *total, http.StatusOK, requestId)
}

// Sensor handlers
func (dc *DeviceController) HandleGetSensor(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	sensorID, err := ctx.Params().GetInt64("sensorID")
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	sensor, err := dc.deviceDomain.FetchSensor(requestId, sensorID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), sensor, http.StatusOK, requestId)
}

func (dc *DeviceController) HandleGetSensors(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	pageSize, page, err := GetPageAndPageSize(ctx)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	paginatedList, total, err := dc.deviceDomain.ListSensors(requestId, *page, *pageSize)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithList(ctx.ResponseWriter(), paginatedList, *page, *pageSize, *total, http.StatusOK, requestId)
}

// Unit handlers
func (dc *DeviceController) HandleGetUnit(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	unitID, err := ctx.Params().GetInt64("unitID")
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	unit, err := dc.deviceDomain.FetchUnit(requestId, unitID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), unit, http.StatusOK, requestId)
}

func (dc *DeviceController) HandleGetUnits(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	pageSize, page, err := GetPageAndPageSize(ctx)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	paginatedList, total, err := dc.deviceDomain.ListUnits(requestId, *page, *pageSize)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithList(ctx.ResponseWriter(), paginatedList, *page, *pageSize, *total, http.StatusOK, requestId)
}

// Model handlers
func (dc *DeviceController) HandleGetModel(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	modelID, err := ctx.Params().GetInt64("modelID")
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	model, err := dc.deviceDomain.FetchModel(requestId, modelID)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithJSON(ctx.ResponseWriter(), model, http.StatusOK, requestId)
}

func (dc *DeviceController) HandleGetModels(ctx iris.Context) {
	requestId := GetRequestID(ctx)
	pageSize, page, err := GetPageAndPageSize(ctx)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	paginatedList, total, err := dc.deviceDomain.ListModels(requestId, *page, *pageSize)
	if err != nil {
		RespondWithError(ctx.ResponseWriter(), requestId, err)
		return
	}

	RespondWithList(ctx.ResponseWriter(), paginatedList, *page, *pageSize, *total, http.StatusOK, requestId)
}
