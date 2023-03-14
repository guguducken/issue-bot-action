package wecom

import (
	"encoding/json"
	"net/http"

	"github.com/guguducken/issue-bot/pkg/util"
)

func (w WecomNotice) SendWecomMessage() error {
	mess, err := json.Marshal(w)
	if err != nil {
		return nil
	}
	resp, err := http.Post(url_notice, string(mess), nil)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return util.StatusError
	}
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
