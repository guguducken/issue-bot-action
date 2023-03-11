package wecom

func GenRowData(cells ...CellData) (row RowData) {
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
