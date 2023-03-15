package wecom

var (
	BLACK         Color      = Color{0, 0, 0, 1}
	RED           Color      = Color{255, 0, 0, 1}
	GREEN         Color      = Color{0, 255, 0, 1}
	BLUE          Color      = Color{0, 0, 255, 1}
	FormatDefault TextFormat = TextFormat{"Microsoft YaHei", 14, false, false, false, false, BLACK}
)

type TokenReply struct {
	Code         int    `json:"errcode"`
	Message      string `json:"errmsg"`
	Token        string `json:"access_token"`
	ExpiredTimes int64  `json:"expires_in"`
}

//---------------------- Create A New Doc-------------------

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

//---------------------- Create A New Doc-------------------

type RowColInfoSend struct {
	Docid string `json:"docid"`
}
type RowColInfoReply struct {
	Errcode int          `json:"errcode"`
	Errmsg  string       `json:"errmsg"`
	Data    []Properties `json:"data"`
}

type SheetDataSend struct {
	Docid   string `json:"docid"`
	SheetID string `json:"sheet_id"`
	Range   string `json:"range"`
}

type SheetDataReply struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Data    struct {
		Result GridData `json:"result"`
	} `json:"data"`
}

// ----------------------------- Modify Tables Request -------------------------------
type ModifyTable struct {
	Docid    string     `json:"docid"`
	Requests []Requests `json:"requests,omitempty"`
	total    int
}

type Requests struct {
	AddSheetRequest        *AddSheetRequest        `json:"add_sheet_request,omitempty"`
	UpdateRangeRequest     *UpdateRangeRequest     `json:"update_range_request,omitempty"`
	DeleteDimensionRequest *DeleteDimensionRequest `json:"delete_dimension_request,omitempty"`
	DeleteSheetRequest     *DeleteSheetRequest     `json:"delete_sheet_request,omitempty"`
}

type AddSheetRequest struct {
	Title       string `json:"title"`
	RowCount    int    `json:"row_count"`
	ColumnCount int    `json:"column_count"`
}

type UpdateRangeRequest struct {
	rows     int
	columns  int
	total    int
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
	Text string `json:"text,omitempty"`
	Link *Link  `json:"link,omitempty"`
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

// ----------------------------- Modify Tables Request -------------------------------

// ----------------------------- Modify Tables Response -------------------------------

type Response struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Data    Data   `json:"data"`
}
type Data struct {
	Responses []UpdateResponse `json:"responses"`
}

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

// ----------------------------- Modify Tables Response -------------------------------

// ------------------------------ WeCom Notice Types -----------------------------------
type WecomNotice struct {
	Msgtype  string    `json:"msgtype"`
	Text     *Text     `json:"text,omitempty"`
	Markdown *Markdown `json:"markdown,omitempty"`
	Image    *Image    `json:"image,omitempty"`
	News     *News     `json:"news,omitempty"`
	File     *File     `json:"file,omitempty"`
}

type File struct {
	MediaID string `json:"media_id"`
}

type Text struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

type Markdown struct {
	Content string `json:"content"`
}

type Image struct {
	Base64 string `json:"base64"`
	Md5    string `json:"md5"`
}

type News struct {
	Articles []Article `json:"articles"`
}
type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Picurl      string `json:"picurl"`
}

type FileUploadReply struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}
