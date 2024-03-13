package main

import (
	"fmt"
	"testing"
)

func TestGetVt(t *testing.T) {
	vt := getVt("0x904c35b8694b846aa878d078535f259a")
	fmt.Println(vt)
}
