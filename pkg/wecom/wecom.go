package wecom

import (
	"os"

	"github.com/guguducken/issue-bot/pkg/util"
)

var (
	url_notice string
)

func init() {
	url_notice = os.Getenv(`INPUT_URL_NOTICE`)
	if url_notice == "" {
		util.Error(`WeCom notice url setting is invalid which will make notice fail, please check again`)
	}
}
