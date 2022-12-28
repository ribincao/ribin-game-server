package test

import (
	"fmt"
	errs "ribin-server/error"
	"testing"
)

func TestError(t *testing.T) {
	eb := errs.New(errs.ConfigErrorCode, "BusinuessError")
	ef := errs.NewFrameworkError(errs.MsgErrorCode, "FrameworkError")
	fmt.Println("eb", eb.Code, eb.Message)
	fmt.Println("ef", ef.Code, ef.Message)
}
