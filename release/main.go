package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/guguducken/auto-release/pkg/github"
	"github.com/guguducken/auto-release/pkg/k8s"
	"github.com/guguducken/auto-release/pkg/util"
)

func main() {
	fmt.Println(`Info: start to get expired time....`)
	expired := os.Getenv(`INPUT_EXPIRED-TIME`)
	arr := strings.Split(expired, " ")
	if len(arr) != 5 {
		panic(`Error: the expired time is invalid, please use this format: day hour minute second millisecond`)
	}
	int_arr := make([]int, 5)
	for i := 0; i < len(arr); i++ {
		t, err := strconv.Atoi(arr[i])
		if err != nil {
			panic(err)
		}
		int_arr[i] = t
	}
	time_expired := int_arr[0]*24 + int_arr[1]
	time_expired = time_expired*60 + int_arr[2]
	time_expired = time_expired*60 + int_arr[3]
	time_expired = time_expired*1000 + int_arr[4]
	fmt.Println(`Info: the expired time is ` + arr[0] + ` days ` + arr[1] + ` hours ` + arr[2] + ` minutes ` + arr[3] + ` seconds ` + arr[4] + ` milliseconds. And the total milliseconds is ` + strconv.Itoa(time_expired))

	fmt.Println(`Info: start to get eks namespaces....`)
	nsl, err := k8s.GetNS()
	if err != nil {
		panic(err)
	}
	fmt.Println(`Info: eks namespaces' total: ` + strconv.Itoa(len(nsl.Items)))

	labels := os.Getenv(`INPUT_LABELS`)
	repo_full := os.Getenv(`GITHUB_REPOSITORY`)
	fmt.Println(`Info: start to get issues by labels: ` + labels + ` from ` + repo_full)
	oop := strings.Split(repo_full, "/")
	q_issues := github.NewIssueListQuery(oop[0], oop[1], ``, `open`, ``, ``, ``, labels, `updated`, ``, ``)
	issues, err := q_issues.GetAllIssues()
	fmt.Println(`Info: the total of issues is ` + strconv.Itoa(len(issues)))
	if err != nil {
		panic(err)
	}
	exclu := getExclusive()
	nsg := parseNS(issues)
	for i := 0; i < len(nsl.Items); i++ {
		fmt.Println(`Info: start to check namespace: ` + nsl.Items[i].Name + ` >>>>>>>>>>>>>>>>>>`)
		if _, ok := exclu[nsl.Items[i].Name]; ok {
			fmt.Println(`SKIP: skip this namespace ` + nsl.Items[i].Name + `, because this namespace is exclusive`)
			continue
		}
		if !util.ExpiredCheck(nsl.Items[i].CreationTimestamp.Time, nsl.Items[i].CreationTimestamp.Time, int64(time_expired)) {
			fmt.Println(`SKIP: skip this namespace ` + nsl.Items[i].Name + `, because this namespace is not expired`)
			continue
		}
		var ans *NS_Github
		for j := 0; j < len(nsg); j++ {
			if nsg[j].namespace == nsl.Items[i].Name {
				ans = &nsg[j]
				break
			}
		}
		if ans == nil {
			//执行删除操作
			fmt.Println(`DELETE: this namespace ` + nsl.Items[i].Name + ` will delete, because there is not corresponding issue`)
		} else {
			//检查时间
			fmt.Println(`INFO: the corresponding issue number is: ` + strconv.Itoa(ans.number) + `, title: ` + ans.title)
			if util.ExpiredCheck(ans.createdAt, ans.updatedAt, int64(time_expired)) {
				//执行删除操作
				fmt.Println(`DELETE: this namespace ` + nsl.Items[i].Name + ` will delete, because the time of issue which corresponding to this namespace is expired`)
			} else {
				fmt.Println(`SKIP: skip this namespace ` + nsl.Items[i].Name + `, because the issue which corresponding to this namespace is not expired`)
			}
		}

	}
}

type NS_Github struct {
	namespace string
	number    int
	title     string
	createdAt time.Time
	updatedAt time.Time
}

func parseNS(arr []github.Issue) []NS_Github {
	nsg := make([]NS_Github, 0, 50)
	for i := 0; i < len(arr); i++ {
		nsg = append(nsg, NS_Github{
			namespace: titleToNS(arr[i].Title),
			number:    arr[i].Number,
			title:     arr[i].Title,
			createdAt: arr[i].CreatedAt,
			updatedAt: arr[i].UpdatedAt,
		})
	}
	return nsg
}

func titleToNS(str string) string {
	i := len(str) - 1
	for ; i >= 0; i-- {
		if str[i] != ' ' {
			break
		}
	}
	end := i
	for ; i >= 0; i-- {
		if str[i] == ' ' || str[i] == ':' {
			break
		}
	}
	return str[i+1 : end+1]
}

func getExclusive() map[string]struct{} {
	exclu := os.Getenv(`INPUT_EXCLUSIVE`)
	m := make(map[string]struct{}, 10)
	left, right := 0, 0
	for left < len(exclu) {
		if exclu[right] != ' ' && exclu[right] != ',' {
			right++
			if right >= len(exclu) {
				m[exclu[left:right]] = struct{}{}
				break
			}
			continue
		}

		m[exclu[left:right]] = struct{}{}
		right++
		left = right

	}
	return m
}
