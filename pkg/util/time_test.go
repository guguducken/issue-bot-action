package util

import (
	"fmt"
	"testing"
)

func Test_time(t *testing.T) {
	str := `2023-03-05T03:00:00Z`
	work, holiday, _ := GetPassedTimeWithoutWeekend_String(str)
	fmt.Printf("work: %v\n", work)
	fmt.Printf("holiday: %v\n", holiday)
}
