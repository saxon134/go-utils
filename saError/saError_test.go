package saError

import (
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	var err1 = Error{Code: SensitiveErrorCode, Msg: "Error error"}
	fmt.Println(err1)
}
