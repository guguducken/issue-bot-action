package wecom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/guguducken/issue-bot/pkg/util"
)

var (
	wecomAPI    string
	corpID      string
	corpSecret  string
	token_wecom string
	startTime   time.Time
	expiredTime int64 = 7200
	url_notice  string
)

func init() {
	wecomAPI = os.Getenv(`INPUT_WECOM_API`)
	corpID = os.Getenv(`INPUT_CORPID`)
	corpSecret = os.Getenv(`INPUT_CORPSECRET`)
	url_notice = os.Getenv(`INPUT_URL_NOTICE`)
	if wecomAPI == "" || corpID == "" || corpSecret == "" {
		panic(`Wecom Settings Error: invalid wecomAPI, corpID or corpSecret, please check again`)
	}
	if url_notice == "" {
		util.Error(`WeCom notice url setting is invalid which will make notice fail, please check again`)
	}
	getToken()
}

func getToken() {
	uri := wecomAPI + `/gettoken?corpid=` + corpID + `&corpsecret=` + corpSecret
	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	reply := TokenReply{}
	err = json.Unmarshal(body, &reply)
	if err != nil {
		panic(err)
	}

	fmt.Printf("reply: %v\n", reply)

	if reply.Code != 0 {
		util.Error(`get access toke failed, the reason is: ` + reply.Message)
		panic(`get access toke failed, the reason is: ` + reply.Message)
	}
	token_wecom = reply.Token
	expiredTime = reply.ExpiredTimes
	startTime = time.Now().In(time.FixedZone(`UTC`, 0))

}
