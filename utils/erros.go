package utils

import (
	"fmt"
	"time"
)

func ShowError(name string, err error) {
	fmt.Println(name, "error at", time.Now(), err.Error())
}
