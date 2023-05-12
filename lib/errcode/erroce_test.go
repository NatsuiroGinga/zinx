package errcode

import (
	"fmt"
	"testing"
)

func TestError_WithDetails(t *testing.T) {
	fmt.Println(PROPERTY_NOT_FOUND.Format("test"))
	fmt.Println(PROPERTY_NOT_FOUND.Format("ddd"))
	fmt.Println(CONNECT_CLOSED.Format(1))
}
