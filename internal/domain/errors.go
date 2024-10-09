package domain

import (
	"errors"
	"net/http"
)

// General errors
var (
	ErrBadPageSize  = errors.New("invalid page size")
	ErrBadPageIndex = errors.New("invalid page index")

	ErrInternalExceptionCode = "ERR_INTERNAL_EXCEPTION"
	ErrInternalExceptionDesc = "An internal server error occurred."
)

// User errors
var ErrNotFoundUserByEmail = errors.New("no user found with given email")
var ErrNotFoundUserByID = errors.New("no user found with given ID")

// Account errors
var ErrNotFoundAccountByEmail = errors.New("no account found with given email")
var ErrNotFoundAccountByID = errors.New("no account found with given ID")

// Address errors
var ErrNotFoundAddressByID = errors.New("no address found with given ID")
var ErrNotFoundAddressByAccountID = errors.New("no address found with given account ID")

var BadPayload = "ERR_BAD_PAYLOAD_FIELDS"
var SuccessCode = "00"
var SuccessMessage = "Request performed successfully"

var (
	ErrCodeMap = map[error]string{
		ErrBadPageSize:                "ERR_BAD_PAGE_SIZE",
		ErrBadPageIndex:               "ERR_BAD_PAGE_INDEX",
		ErrNotFoundUserByEmail:        "ERR_NOT_FOUND_USER_BY_EMAIL",
		ErrNotFoundUserByID:           "ERR_NOT_FOUND_USER_BY_ID",
		ErrNotFoundAccountByEmail:     "ERR_NOT_FOUND_ACCOUNT_BY_EMAIL",
		ErrNotFoundAccountByID:        "ERR_NOT_FOUND_ACCOUNT_BY_ID",
		ErrNotFoundAddressByID:        "ERR_NOT_FOUND_ADDRESS_BY_ID",
		ErrNotFoundAddressByAccountID: "ERR_NOT_FOUND_ADDRESS_BY_ACCOUNT_ID",
	}

	ErrDescriptionMap = map[error]string{
		ErrBadPageSize:                "The page size provided is invalid.",
		ErrBadPageIndex:               "The page index provided is invalid.",
		ErrNotFoundUserByEmail:        "No user found with the given email.",
		ErrNotFoundUserByID:           "No user found with the given ID.",
		ErrNotFoundAccountByEmail:     "No account found with the given email.",
		ErrNotFoundAccountByID:        "No account found with the given ID.",
		ErrNotFoundAddressByID:        "No address found with the given ID.",
		ErrNotFoundAddressByAccountID: "No address found with the given account ID.",
	}

	ErrToHTTPStatus = map[error]int{
		ErrBadPageSize:                http.StatusBadRequest,
		ErrBadPageIndex:               http.StatusBadRequest,
		ErrNotFoundUserByEmail:        http.StatusNotFound,
		ErrNotFoundUserByID:           http.StatusNotFound,
		ErrNotFoundAccountByEmail:     http.StatusNotFound,
		ErrNotFoundAccountByID:        http.StatusNotFound,
		ErrNotFoundAddressByID:        http.StatusNotFound,
		ErrNotFoundAddressByAccountID: http.StatusNotFound,
	}
)
