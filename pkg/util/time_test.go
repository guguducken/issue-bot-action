package util

import (
	"fmt"
	"testing"
)

func Test_time(t *testing.T) {
	str := `2023-03-05T03:00:00Z`
	work, holiday, _ := GetPassedTimeWithoutWeekend_String(str)
	t1, _ := MStoTime(work)
	fmt.Printf("work: %v\n", t1)
	t2, _ := MStoTime(holiday)
	fmt.Printf("holiday: %v\n", t2)
}
