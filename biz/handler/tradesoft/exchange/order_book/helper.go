package order_book

import (
	"context"
	"reflect"

	"github.com/cloudwego/hertz/pkg/app"
)

func handleResponse(ctx context.Context, c *app.RequestContext, resp interface{}, err error) {
	if err == nil {
		c.JSON(200, resp)
		return
	}
	status := 500
	setErrorResponse(resp, err)
	c.JSON(status, resp)
}

func setErrorResponse(resp interface{}, err error) {
	code := int32(-1)
	if bizCode := reflect.Indirect(reflect.ValueOf(resp)).FieldByName("BizCode"); bizCode.CanSet() && bizCode.Kind() == reflect.Int32 {
		bizCode.SetInt(int64(code))
	}
	if errMsg := reflect.Indirect(reflect.ValueOf(resp)).FieldByName("ErrMsg"); errMsg.CanSet() && errMsg.Kind() == reflect.String {
		errMsg.SetString(err.Error())
	}
}
