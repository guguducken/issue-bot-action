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

func Test_NewUpdateRangeRequest(t *testing.T) {
	cells := make([]CellData, 0, 2)
	cells = append(cells, GenTextCellData(`name is aaa`, FormatDefault))
	cells = append(cells, GenTextCellData(`name is bbb`, FormatDefault))

	row := GenRowData(cells)
	mt := NewModifyTable(`1`)
	mt.NewUpdateRangeRequest(`111`, 0, 0, row)
	mt_json, _ := json.Marshal(mt)
	fmt.Printf("mt: %v\n", string(mt_json))
}
