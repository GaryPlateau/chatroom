package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler(condition string, err error) {
	if err != nil {
		fmt.Println(condition, err)
		return
	}
}

// 返回错误信息 ErrorResponse
func ErrorResPonse(err error) Response {
	if _, ok := err.(validator.ValidationErrors); ok {
		return Response{
			Status: 400,
			Msg:    "参数错误",
			Error:  fmt.Sprint(err),
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return Response{
			Status: 400,
			Msg:    "JSON类型不匹配",
			Error:  fmt.Sprint(err),
		}
	}
	return Response{
		Status: 400,
		Msg:    "参数错误",
		Error:  fmt.Sprint(err),
	}
}
