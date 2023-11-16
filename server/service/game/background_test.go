package game

import (
	"fmt"
	"testing"
	"time"
)

// experiments
// -----------------------------------------------------------------------

func Test_time_Sub(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(5 * time.Hour)

	// ***

	fmt.Println(t2.Sub(t1))
}
