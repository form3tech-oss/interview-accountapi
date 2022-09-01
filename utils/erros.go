package utils

import (
	"fmt"
	"time"
)

func ShowError(name string, err error) {
	fmt.Println(name, "error at", time.Now().Format("2006-01-02T15:04:05"), err.Error())
}
