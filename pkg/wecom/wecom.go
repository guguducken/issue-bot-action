package wecom

import (
	"os"

	"github.com/guguducken/issue-bot/pkg/util"
)

var (
	notice_key string
)

func init() {
	notice_key = os.Getenv(`INPUT_NOTICE_KEY`)
	if notice_key == "" {
		util.Error(`WeCom notice url setting is invalid which will make notice fail, please check again`)
	}
}
