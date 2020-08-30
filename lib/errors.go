package lib

import (
	"fmt"
	"runtime/debug"
)

type FameError struct {
	ErrorCode    string
	ErrorMessage string
	Caption      string
	CaptionData  map[string]interface{}
	StackTrace   []byte
}

func (err *FameError) IsObjectNotFoundError() bool {
	return err.ErrorCode == "ObjectNotFound"
}

func (err *FameError) String() string {
	return fmt.Sprintf("%s: %s", err.ErrorCode, err.ErrorMessage)
}

func InvalidParamsError(e error) *FameError {
	return &FameError{
		ErrorCode:    "InvalidParams",
		ErrorMessage: e.Error(),
		Caption:      "ERR_INVALID_PARAMS",
		CaptionData:  nil,
		StackTrace:   debug.Stack(),
	}
}

func ObjectNotFoundError(e error) *FameError {
	return &FameError{
		ErrorCode:    "ObjectNotFound",
		ErrorMessage: e.Error(),
		Caption:      "ERR_OBJECT_NOT_FOUND",
		CaptionData:  nil,
		StackTrace:   debug.Stack(),
	}
}

func PrivilegeError(e error) *FameError {
	return &FameError{
		ErrorCode:    "PrivilegeError",
		ErrorMessage: e.Error(),
		Caption:      "ERR_INVALID_PRIVILEGES",
		CaptionData:  nil,
		StackTrace:   debug.Stack(),
	}
}

func DataCorruptionError(e error) *FameError {
	return &FameError{
		ErrorCode:    "DataError",
		ErrorMessage: e.Error(),
		Caption:      "ERR_DATA_CORRUPTION",
		CaptionData:  nil,
		StackTrace:   debug.Stack(),
	}
}

func InternalError(e error) *FameError {
	return &FameError{
		ErrorCode:    "InternalError",
		ErrorMessage: e.Error(),
		Caption:      "ERR_INTERNAL",
		CaptionData:  nil,
		StackTrace:   debug.Stack(),
	}
}

func WorkflowError(e error) *FameError {
	return &FameError{
		ErrorCode:    "WorkflowError",
		ErrorMessage: e.Error(),
		Caption:      "ERR_WORKFLOW",
		CaptionData:  nil,
		StackTrace:   debug.Stack(),
	}
}

func AuthenticationError(e error) *FameError {
	return &FameError{
		ErrorCode:    "AuthenticationError",
		ErrorMessage: e.Error(),
		Caption:      "ERR_AUTHENTICATION",
		CaptionData:  nil,
		StackTrace:   debug.Stack(),
	}
}
