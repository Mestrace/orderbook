package order_book

import (
	"context"
	"net/http"
	"reflect"

	"github.com/cloudwego/hertz/pkg/app"
)

func handleResponse(_ context.Context, reqCtxt *app.RequestContext, resp interface{}, err error) {
	if err == nil {
		reqCtxt.JSON(http.StatusOK, resp)

		return
	}

	status := http.StatusInternalServerError

	setErrorResponse(resp, err)

	reqCtxt.JSON(status, resp)
}

func setErrorResponse(resp interface{}, err error) {
	code := int32(-1)

	if bizCode := reflect.Indirect(reflect.ValueOf(resp)).
		FieldByName("BizCode"); bizCode.CanSet() && bizCode.Kind() == reflect.Int32 {
		bizCode.SetInt(int64(code))
	}

	if errMsg := reflect.Indirect(reflect.ValueOf(resp)).
		FieldByName("ErrMsg"); errMsg.CanSet() && errMsg.Kind() == reflect.String {
		errMsg.SetString(err.Error())
	}
}
