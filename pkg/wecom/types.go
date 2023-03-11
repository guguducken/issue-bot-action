package wecom

type TokenReply struct {
	Code         int    `json:"errcode"`
	Message      string `json:"errmsg"`
	Token        string `json:"access_token"`
	ExpiredTimes int64  `json:"expires_in"`
}

type NewDocSend struct {
	Spaceid    string   `json:"spaceid"`
	Fatherid   string   `json:"fatherid"`
	DocType    string   `json:"doc_type"`
	DocName    string   `json:"doc_name"`
	AdminUsers []string `json:"admin_users"`
}

type NewDocReply struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	URL     string `json:"url"`
	Docid   string `json:"docid"`
}

// ----------------------------- Modify Tables Request Start -------------------------------
type ModifyTable struct {
	Docid    string     `json:"docid"`
	Requests []Requests `json:"requests"`
}

type Requests struct {
	AddSheetRequest        AddSheetRequest        `json:"add_sheet_request,omitempty"`
	UpdateRangeRequest     UpdateRangeRequest     `json:"update_range_request,omitempty"`
	DeleteDimensionRequest DeleteDimensionRequest `json:"delete_dimension_request,omitempty"`
	DeleteSheetRequest     DeleteSheetRequest     `json:"delete_sheet_request,omitempty"`
}

type AddSheetRequest struct {
	Title       string `json:"title"`
	RowCount    int    `json:"row_count"`
	ColumnCount int    `json:"column_count"`
}

type UpdateRangeRequest struct {
	SheetID  string   `json:"sheet_id"`
	GridData GridData `json:"grid_data"`
}

// Dimension only two selection: ROW or COLUMN
// Index are all start at one
type DeleteDimensionRequest struct {
	SheetID    string `json:"sheet_id"`
	Dimension  string `json:"dimension"`
	StartIndex int    `json:"start_index"`
	EndIndex   int    `json:"end_index"`
}

type DeleteSheetRequest struct {
	SheetID string `json:"sheet_id"`
}

type GridData struct {
	StartRow    int       `json:"start_row"`
	StartColumn int       `json:"start_column"`
	Rows        []RowData `json:"rows"`
}

type RowData struct {
	Values []CellData `json:"values"`
}

type CellData struct {
	CellValue  CellValue  `json:"cell_value"`
	CellFormat CellFormat `json:"cell_format"`
}

type CellFormat struct {
	TextFormat TextFormat `json:"text_format"`
}

// only support Text or Link, only one
type CellValue struct {
	Text string `json:"text"`
	Link Link   `json:"link"`
}

type TextFormat struct {
	Font          string `json:"font"`
	FontSize      int    `json:"font_size"`
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Color         Color  `json:"color"`
}

type Link struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

type Color struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
	Alpha int `json:"alpha"`
}

// ----------------------------- Modify Tables Request End -------------------------------

// ----------------------------- Modify Tables Response Start -------------------------------

type UpdateResponse struct {
	AddSheetResponse        AddSheetResponse        `json:"add_sheet_response"`
	UpdateRangeResponse     UpdateRangeResponse     `json:"update_range_response"`
	DeleteDimensionResponse DeleteDimensionResponse `json:"delete_dimension_response"`
	DeleteSheetResponse     DeleteSheetResponse     `json:"delete_sheet_response"`
}
type Properties struct {
	SheetID     string `json:"sheet_id"`
	Title       string `json:"title"`
	RowCount    int    `json:"row_count"`
	ColumnCount int    `json:"column_count"`
}
type AddSheetResponse struct {
	Properties Properties `json:"properties"`
}
type UpdateRangeResponse struct {
	UpdatedCells int `json:"updated_cells"`
}
type DeleteDimensionResponse struct {
	Deleted int `json:"deleted"`
}
type DeleteSheetResponse struct {
	SheetID string `json:"sheet_id"`
}

// ----------------------------- Modify Tables Response End -------------------------------
