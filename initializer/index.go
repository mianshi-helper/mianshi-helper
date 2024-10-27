package inInitializer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Dialogue struct {
	RequestID      string `json:"request_id"`      // 修正了字段名和JSON标签
	ConversationID string `json:"conversation_id"` // 修正了字段名和JSON标签
}

func CreateDialogue() string {
	url := "https://qianfan.baidubce.com/v2/app/conversation"
	payload := strings.NewReader(`{"app_id":"6f7aef3e-3db1-434d-ac74-bc3199477d27"}`)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Appbuilder-Authorization", "Bearer bce-v3/ALTAK-dfpyIHGrYVav9sBP6AZp7/d81d889bc31f8af7a6cd244ee60a8e83561ce6a4")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	fmt.Println(string(body))
	dialogue := Dialogue{}
	err = json.Unmarshal(body, &dialogue)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return dialogue.ConversationID
}
