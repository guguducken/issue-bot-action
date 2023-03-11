package wecom

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_cell(t *testing.T) {
	cell := GenTextCellData(`test`, FormatDefault)
	ans, err := json.Marshal(cell)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("ans: %v\n", string(ans))

	cell = GenLinkCellData(`test`, `test`, FormatDefault)
	ans, err = json.Marshal(cell)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("ans: %v\n", string(ans))
}

func Test_UpdateTable(t *testing.T) {
	updateRR := UpdateRangeRequest{}
	ans, err := json.Marshal(updateRR)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("ans: %v\n", string(ans))
}
