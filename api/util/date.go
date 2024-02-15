package util

import (
	"fmt"
	"time"
)

func GenerateTimestamp() string {

	yy, mm, dd := time.Now().Date()
	hh, MM, ss := time.Now().Clock()

	return fmt.Sprintf("%d-%s-%d %d:%d:%d", yy, mm, dd, hh, MM, ss)
}
