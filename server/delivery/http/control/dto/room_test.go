package dto

import (
	"fmt"
	"testing"
	"time"
)

// experiments
// -----------------------------------------------------------------------

func Test_time(t *testing.T) {
	currentTime := time.Now()
	fmt.Println(currentTime)
	fmt.Println(currentTime.UTC())

	fmt.Println(currentTime.Unix())
	fmt.Println(currentTime.UTC().Unix())
}
