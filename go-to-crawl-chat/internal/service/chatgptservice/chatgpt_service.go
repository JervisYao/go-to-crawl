package chatgptservice

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"io/ioutil"
)

var (
	api    = "https://api.openai.com/v1/chat/completions"
	apiKey = "sk-Pze776gAuHIfwMGMaUXKUT3BlbkFJKA2cSKWKb4trEDbgoauQ"
)

func requestOpenAI(question string) {

	msg := Msg{
		Role:    "user",
		Content: question,
	}
	var msgs []*Msg
	msgs = append(msgs, &msg)

	reqBody := ReqBody{
		Model:    "gpt-3.5-turbo",
		Messages: msgs,
	}

	reqJson := gjson.New(reqBody).MustToJsonString()

	response, err := g.Client().SetHeader("Authorization", "Bearer "+apiKey).ContentJson().Post(gctx.GetInitCtx(), api, reqJson)
	if err != nil {
		return
	}

	all, err := ioutil.ReadAll(response.Body)
	fmt.Println("ChatGPT: \n", string(all))

}

type ReqBody struct {
	Model    string `json:"model"`
	Messages []*Msg `json:"messages"`
}
type Msg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
