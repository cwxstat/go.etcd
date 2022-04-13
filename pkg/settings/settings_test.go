package settings

import (
	"fmt"
	"testing"
)

func Test_Read(t *testing.T) {
	TestRead()
}

func TestCreateDefault(t *testing.T) {
	CreateDefault()
}

func TestReadConfig(t *testing.T) {
	r, err := ReadConfig()
	if err != nil {
		t.Fatalf("err: %v\n", err)
	}
	fmt.Println(r)
}
