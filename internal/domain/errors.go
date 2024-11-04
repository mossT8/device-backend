package domain

import (
	"errors"
	"net/http"
)

// General errors
var (
	ErrBadPageSize  = errors.New("invalid page size")
	ErrBadPageIndex = errors.New("invalid page index")
	ErrUnauthorized = errors.New("unauthorized access")

	ErrInternalExceptionCode = "ERR_INTERNAL_EXCEPTION"
	ErrInternalExceptionDesc = "An internal server error occurred."
)

// User errors
var ErrNotFoundUserByEmail = errors.New("no user found with the given email")
var ErrNotFoundUserByID = errors.New("no user found with the given ID")

// Account errors
var ErrNotFoundAccountByEmail = errors.New("no account found with the given email")
var ErrNotFoundAccountByID = errors.New("no account found with the given ID")

// Address errors
var ErrNotFoundAddressByID = errors.New("no address found with the given ID")
var ErrNotFoundAddressByAccountID = errors.New("no address found with the given account ID")

// Model errors
var ErrNotFoundModelByID = errors.New("no model found with the given ID")

// Unit errors
var ErrNotFoundUnitByID = errors.New("no unit found with the given ID")
var ErrNotFoundUnitByName = errors.New("no unit found with the given name")

// Sensor errors
var ErrNotFoundSensorByID = errors.New("no sensor found with the given ID")
var ErrNotFoundSensorByCode = errors.New("no sensor found with the given code")

// Device errors
var ErrNotFoundDeviceByID = errors.New("no device found with the given ID")
var ErrNotOwnedDeviceByID = errors.New("no device for account ID provided")
var ErrNotFoundDeviceBySerialNumber = errors.New("no device found with the given serial number")
var ErrSerialNumberNotMatch = errors.New("serial number does not match")
var ErrModelNotMatch = errors.New("model does not match")
var ErrDeviceAndAccountNotMatch = errors.New("device and account do not match")

var BadPayload = "ERR_BAD_PAYLOAD_FIELDS"
var SuccessCode = "00"
var SuccessMessage = "Request performed successfully"

// JWT errors
var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("token has expired")
var ErrMalformedToken = errors.New("malformed token")
var ErrMissingToken = errors.New("missing token")
var ErrInvalidClaims = errors.New("invalid token claims")

var (
	ErrCodeMap = map[error]string{
		ErrInvalidToken:                 "ERR_BAD_TOKEN",
		ErrExpiredToken:                 "ERR_EXPIRED_TOKEN",
		ErrMalformedToken:               "ERR_MALFORMED_TOKEN",
		ErrMissingToken:                 "ERR_MISSING_TOKEN",
		ErrInvalidClaims:                "ERR_BAD_TOKEN_CLAIMS",
		ErrBadPageSize:                  "ERR_BAD_PAGE_SIZE",
		ErrBadPageIndex:                 "ERR_BAD_PAGE_INDEX",
		ErrNotFoundUserByEmail:          "ERR_NOT_FOUND_USER_BY_EMAIL",
		ErrNotFoundUserByID:             "ERR_NOT_FOUND_USER_BY_ID",
		ErrNotFoundAccountByEmail:       "ERR_NOT_FOUND_ACCOUNT_BY_EMAIL",
		ErrNotFoundAccountByID:          "ERR_NOT_FOUND_ACCOUNT_BY_ID",
		ErrNotFoundAddressByID:          "ERR_NOT_FOUND_ADDRESS_BY_ID",
		ErrNotFoundAddressByAccountID:   "ERR_NOT_FOUND_ADDRESS_BY_ACCOUNT_ID",
		ErrNotFoundModelByID:            "ERR_NOT_FOUND_MODEL_BY_ID",
		ErrNotFoundUnitByID:             "ERR_NOT_FOUND_UNIT_BY_ID",
		ErrNotFoundUnitByName:           "ERR_NOT_FOUND_UNIT_BY_NAME",
		ErrNotFoundSensorByID:           "ERR_NOT_FOUND_SENSOR_BY_ID",
		ErrNotFoundSensorByCode:         "ERR_NOT_FOUND_SENSOR_BY_CODE",
		ErrNotFoundDeviceByID:           "ERR_NOT_FOUND_DEVICE_BY_ID",
		ErrNotOwnedDeviceByID:           "ERR_NOT_OWNED_DEVICE_BY_ID",
		ErrNotFoundDeviceBySerialNumber: "ERR_NOT_FOUND_DEVICE_BY_SERIAL_NUMBER",
		ErrSerialNumberNotMatch:         "ERR_SERIAL_NUMBER_NOT_MATCH",
		ErrModelNotMatch:                "ERR_MODEL_NOT_MATCH",
		ErrDeviceAndAccountNotMatch:     "ERR_DEVICE_AND_ACCOUNT_NOT_MATCH",
		ErrUnauthorized:                 "ERR_UNAUTHORIZED",
	}

	ErrDescriptionMap = map[error]string{
		ErrUnauthorized:                 "Unauthorized access.",
		ErrInvalidToken:                 "The token provided is invalid.",
		ErrExpiredToken:                 "The token provided has expired.",
		ErrMalformedToken:               "The token provided is malformed.",
		ErrMissingToken:                 "No token provided.",
		ErrInvalidClaims:                "The token claims are invalid.",
		ErrBadPageSize:                  "The page size provided is invalid.",
		ErrBadPageIndex:                 "The page index provided is invalid.",
		ErrNotFoundUserByEmail:          "No user found with the given email.",
		ErrNotFoundUserByID:             "No user found with the given ID.",
		ErrNotFoundAccountByEmail:       "No account found with the given email.",
		ErrNotFoundAccountByID:          "No account found with the given ID.",
		ErrNotFoundAddressByID:          "No address found with the given ID.",
		ErrNotFoundAddressByAccountID:   "No address found with the given account ID.",
		ErrNotFoundModelByID:            "No model found with the given ID.",
		ErrNotFoundUnitByID:             "No unit found with the given ID.",
		ErrNotFoundUnitByName:           "No unit found with the given name.",
		ErrNotFoundSensorByID:           "No sensor found with the given ID.",
		ErrNotFoundSensorByCode:         "No sensor found with the given code.",
		ErrNotFoundDeviceByID:           "No device found with the given ID.",
		ErrNotOwnedDeviceByID:           "No device found for the provided account ID.",
		ErrNotFoundDeviceBySerialNumber: "No device found with the given serial number.",
		ErrSerialNumberNotMatch:         "The serial number does not match.",
		ErrModelNotMatch:                "The model does not match.",
		ErrDeviceAndAccountNotMatch:     "The device and account do not match.",
	}

	ErrToHTTPStatus = map[error]int{
		ErrUnauthorized:                 http.StatusUnauthorized,
		ErrBadPageSize:                  http.StatusBadRequest,
		ErrBadPageIndex:                 http.StatusBadRequest,
		ErrNotFoundUserByEmail:          http.StatusNotFound,
		ErrNotFoundUserByID:             http.StatusNotFound,
		ErrNotFoundAccountByEmail:       http.StatusNotFound,
		ErrNotFoundAccountByID:          http.StatusNotFound,
		ErrNotFoundAddressByID:          http.StatusNotFound,
		ErrNotFoundAddressByAccountID:   http.StatusNotFound,
		ErrNotFoundModelByID:            http.StatusNotFound,
		ErrNotFoundUnitByID:             http.StatusNotFound,
		ErrNotFoundUnitByName:           http.StatusNotFound,
		ErrNotFoundSensorByID:           http.StatusNotFound,
		ErrNotFoundSensorByCode:         http.StatusNotFound,
		ErrNotFoundDeviceByID:           http.StatusNotFound,
		ErrNotOwnedDeviceByID:           http.StatusNotFound,
		ErrNotFoundDeviceBySerialNumber: http.StatusNotFound,
		ErrSerialNumberNotMatch:         http.StatusBadRequest,
		ErrModelNotMatch:                http.StatusBadRequest,
		ErrDeviceAndAccountNotMatch:     http.StatusBadRequest,
	}
)
