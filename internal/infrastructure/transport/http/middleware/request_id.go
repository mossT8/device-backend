package middleware

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/constants"
)

func RequestIDMiddleware(ctx iris.Context) {
	requestId := ctx.Request().Header.Get(constants.CTXRequestIdKey)
	if len(requestId) == 0 {
		requestId = uuid.NewString()
	}
	ctx.Values().Set(constants.CTXRequestIdKey, requestId)
	ctx.Next()
}
