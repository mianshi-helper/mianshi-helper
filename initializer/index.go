package inInitializer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Dialogue struct {
	RequestID      string `json:"request_id"`      // 修正了字段名和JSON标签
	ConversationID string `json:"conversation_id"` // 修正了字段名和JSON标签
}

func CreateDialogue() string {
	url := "http://localhost:3099/create"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	req.Header.Add("Content-Type", "application/json")

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

	var response struct {
		SessionID string `json:"sessionId"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	return response.SessionID
}
