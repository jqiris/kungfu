package validator

import (
	"github.com/jqiris/kungfu/v2/logger"
)

type ValidateItem struct {
	Con        bool
	Msg        string
	Code       int
	PrintError bool
	PrintFmt   string
	PrintArgs  []any
}

type Validate struct {
	StdOk    int
	StdError int
	list     []ValidateItem
}

func (v *Validate) Condition(con bool, msg string, args ...int) *Validate {
	code := v.StdError
	if len(args) > 0 {
		code = args[0]
	}
	v.list = append(v.list, ValidateItem{Con: con, Msg: msg, Code: code})
	return v
}
func (v *Validate) PrintError(fmt string, args ...any) {
	if len(v.list) > 0 {
		k := len(v.list) - 1
		v.list[k].PrintError = true
		v.list[k].PrintFmt = fmt
		v.list[k].PrintArgs = args
	}
}
func (v *Validate) Verify() (bool, string, int) {
	hasError := false
	errMsg := ""
	errCode := v.StdOk
	for _, item := range v.list {
		if item.Con {
			hasError = true
			errMsg = item.Msg
			errCode = item.Code
			if item.PrintError {
				logger.Errorf(item.PrintFmt, item.PrintArgs...)
			}
			break
		}
	}
	v.list = []ValidateItem{}
	return hasError, errMsg, errCode
}

func (v *Validate) IsPass() bool {
	isPass := true
	for _, item := range v.list {
		if item.Con {
			isPass = false
			break
		}
	}
	return isPass
}

func NewValidator(stdOk, stdError int) *Validate {
	return &Validate{
		StdOk:    stdOk,
		StdError: stdError,
		list:     make([]ValidateItem, 0),
	}
}
