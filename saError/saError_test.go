package saError

import (
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	var err1 = Error{Code: SensitiveErrorCode, Msg: "Error error"}
	var err2 = StackError(err1)
	fmt.Println(err2)
}
