package wecom

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/guguducken/issue-bot/pkg/util"
)

func GetTableInfo(docid string) (rcinfo RowColInfoReply, err error) {
	send := RowColInfoSend{
		Docid: docid,
	}
	url := wecomAPI + `/wedoc/spreadsheet/get_sheet_properties?access_token=` + token_wecom
	send_json, err := json.Marshal(send)
	if err != nil {
		return
	}
	resp, err := http.Post(url, string(send_json), nil)
	if err != nil {
		return
	}
	reply, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	err = json.Unmarshal(reply, &rcinfo)
	return
}

func NewModifyTable(docid string) (mt ModifyTable) {
	return ModifyTable{
		Docid:    docid,
		Requests: make([]Requests, 0, 4),
		total:    0,
	}
}

func (mt *ModifyTable) NewUpdateRangeRequest(sheet_id string, start_row, start_column int, rowdatas ...RowData) {
	if mt.total >= 4 {
		util.Error(`the max of request in one ModifyTable is 4`)
		return
	}
	mt.Requests = append(mt.Requests, Requests{
		UpdateRangeRequest: &UpdateRangeRequest{
			SheetID: sheet_id,
			GridData: GridData{
				StartRow:    start_row,
				StartColumn: start_column,
				Rows:        rowdatas,
			},
		},
	})
	mt.total++
}

func GenRowData(cells []CellData) (row RowData) {
	return RowData{
		Values: cells,
	}
}

func GenTextCellData(message string, format TextFormat) (cell CellData) {
	return CellData{
		CellFormat: CellFormat{
			TextFormat: format,
		},
		CellValue: CellValue{
			Text: message,
		},
	}
}

func GenLinkCellData(url, message string, format TextFormat) (cell CellData) {
	return CellData{
		CellFormat: CellFormat{
			TextFormat: format,
		},
		CellValue: CellValue{
			Link: &Link{
				URL:  url,
				Text: message,
			},
		},
	}
}
