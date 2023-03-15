package wecom

import (
	"fmt"
	"testing"
)

func Test_SendWecomFile(t *testing.T) {
	message := "# <font color=\"red\">上下标</font>\n|      算式      |   markdown   |\n| :------------: | :----------: |\n| $a_0, a_{pre}$ | a_0, a_{pre} |\n| $a^0, a^{[0]}$ | a^0, a^{[0]} |\n"
	wn, err := GenFileMessage(message, `test.md`)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("wn.SendWecomMessage(): %v\n", wn.SendWecomMessage())

}
