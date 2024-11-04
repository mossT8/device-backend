package http

import (
	"bytes"
	"cmp"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"mossT8.github.com/device-backend/internal/domain"
	"mossT8.github.com/device-backend/internal/infrastructure/logger"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/constants"
	httpType "mossT8.github.com/device-backend/internal/infrastructure/transport/http/types"
)

var validate = validator.New()

// Helper function to generate a refresh token
func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func RespondWithJSON(w http.ResponseWriter, payload interface{}, code int, requestId string) {
	wappedPayload := httpType.DefaultData{
		RequestID: requestId,
		Status:    domain.SuccessCode,
		Message:   domain.SuccessMessage,
		Data:      payload,
	}
	response, err := json.Marshal(wappedPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(code)

	logger.Infof(requestId, constants.RspFormatLogging, string(response))

	if _, err := w.Write(response); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, constants.ErrFormatLogging, err)
	}
}

func RespondWithList(w http.ResponseWriter, list interface{}, page, pageSIze, total int64, code int, requestId string) {
	wrappedList := httpType.DefaultList{
		RequestID: requestId,
		Status:    domain.SuccessCode,
		Message:   domain.SuccessMessage,
		Page:      page,
		PageSize:  pageSIze,
		Total:     total,
		Data:      list,
	}
	response, err := json.Marshal(wrappedList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(code)

	logger.Infof(requestId, constants.RspFormatLogging, string(response))

	if _, err := w.Write(response); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, constants.ErrFormatLogging, err)
	}
}

func GetRequestID(ctx iris.Context) string {
	return ctx.Values().GetString(constants.CTXRequestIdKey)
}

func GetPageAndPageSize(ctx iris.Context) (*int64, *int64, error) {
	pageSize, sErr := ctx.URLParamInt64(constants.URLPageSizeKey)
	if sErr != nil && sErr.Error() == "not found" {
		pageSize = constants.DefaultPageSize
	} else if pageSize < 0 {
		return nil, nil, domain.ErrBadPageSize
	}

	pageIndex, iErr := ctx.URLParamInt64(constants.URLPageKey)
	if iErr != nil && iErr.Error() == "not found" {
		pageIndex = constants.DefaultIndex
	} else if pageIndex < 0 {
		return nil, nil, domain.ErrBadPageIndex
	}

	return &pageSize, &pageIndex, nil
}

func RespondWithMappingError(w http.ResponseWriter, reason, requestId string) {
	response, err := json.Marshal(&httpType.DefaultErrorResponse{
		Error:     fmt.Sprintf("Bad Request: %s", reason),
		RequestId: requestId,
		Code:      domain.BadPayload,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(http.StatusBadRequest)
	logger.Infof(requestId, constants.ErrFormatLogging, string(response))

	if _, err := w.Write(response); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, constants.ErrFormatLogging, err)
	}
}

func GetRequest(request *http.Request, v interface{}) error {
	if request.ContentLength == 0 {
		return nil
	}

	rs, _ := io.ReadAll(request.Body)
	_ = request.Body.Close()
	request.Body = io.NopCloser(bytes.NewBuffer(rs))

	if err := json.Unmarshal(rs, &v); err != nil {
		return err
	}

	if err := validate.Struct(v); err != nil {
		return err
	}

	return nil
}

func RespondWithError(w http.ResponseWriter, requestId string, errReason error) {

	response, err := json.Marshal(
		&httpType.DefaultErrorResponse{
			RequestId: requestId,
			Code: cmp.Or(domain.ErrCodeMap[errReason], domain.ErrCodeMap[errReason],
				domain.ErrInternalExceptionCode),
			Error: cmp.Or(domain.ErrDescriptionMap[errReason], domain.ErrDescriptionMap[errReason],
				domain.ErrInternalExceptionDesc),
		},
	)
	if err != nil {
		http.Error(w, domain.ErrInternalExceptionDesc, http.StatusInternalServerError)
		return
	}

	w.Header().Set(constants.ContentType, constants.ApplicationJson)
	w.WriteHeader(cmp.Or(domain.ErrToHTTPStatus[errReason], http.StatusBadRequest))
	logger.Infof(requestId, "unable to perform request due to: %s", errReason.Error())
	logger.Infof(requestId, "out going response : '%s'", string(response))

	if _, err := w.Write(response); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, constants.ErrFormatLogging, err)
	}
}
