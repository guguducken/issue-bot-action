package wecom

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/guguducken/issue-bot/pkg/util"
)

func (w WecomNotice) SendWecomMessage() error {
	mess, err := json.Marshal(w)
	if err != nil {
		return nil
	}
	fmt.Printf("mess: %v\n", string(mess))
	url := `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=` + notice_key
	resp, err := http.Post(url, `application/json;charset=utf-8`, strings.NewReader(string(mess)))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return util.StatusError
	}
	rr, _ := io.ReadAll(resp.Body)
	fmt.Printf("string(rr): %v\n", string(rr))
	return nil
}

func GenTextMessage(content string, mentioned_list, mentioned_mobile_list []string) WecomNotice {
	return WecomNotice{
		Msgtype: "text",
		Text: &Text{
			Content:             content,
			MentionedList:       mentioned_list,
			MentionedMobileList: mentioned_mobile_list,
		},
	}
}

func GenMarkdownMessage(content string, mentioned_list []string) WecomNotice {
	for i := 0; i < len(mentioned_list); i++ {
		content += `<@` + mentioned_list[i] + `>`
	}
	return WecomNotice{
		Msgtype: `markdown`,
		Markdown: &Markdown{
			Content: content,
		},
	}
}

func sendWecomFile(message string, fileName string) (fp FileUploadReply, err error) {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	miniHeader := textproto.MIMEHeader{}
	miniHeader.Add(`Content-Disposition`, `form-data; name="media";filename="`+fileName+`"; filelength=`+strconv.Itoa(len(message)))
	miniHeader.Add(`Content-Type`, `application/octet-stream`)
	partWriter, err := bodyWriter.CreatePart(miniHeader)
	partWriter.Write([]byte(message))
	if err != nil {
		return
	}
	fmt.Printf("bodyBuffer.String(): %v\n", bodyBuffer.String())
	err = bodyWriter.WriteField(`name`, `media`)
	if err != nil {
		return
	}

	url := `https://qyapi.weixin.qq.com/cgi-bin/webhook/upload_media?key=` + notice_key + `&type=file`
	req, err := http.NewRequest(`POST`, url, bodyBuffer)
	if err != nil {
		return
	}
	req.Header.Set(`Content-Type`, bodyWriter.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	reply, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(reply, &fp)
	if err != nil {
		return
	}
	if fp.Errcode != 0 {
		return fp, errors.New(`Upload file fail: ` + fp.Errmsg)
	}
	return
}

func GenFileMessage(message, name string) (wn WecomNotice, err error) {
	fp, err := sendWecomFile(message, name)
	if err != nil {
		return
	}
	return WecomNotice{
		Msgtype: `file`,
		File: &File{
			MediaID: fp.MediaID,
		},
	}, nil
}
