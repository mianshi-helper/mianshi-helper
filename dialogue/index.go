package Dialogue

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"mianshi-helper/config"
)

type Response struct {
	RequestID      string    `json:"request_id"`
	Date           time.Time `json:"date"`
	Answer         string    `json:"answer"`
	ConversationID string    `json:"conversation_id"`
	MessageID      string    `json:"message_id"`
	IsCompletion   *bool     `json:"is_completion"` // Changed to *bool to handle null values
	Content        []Content `json:"content"`
}

type Content struct {
	ResultType   string  `json:"result_type"`
	EventCode    int     `json:"event_code"`
	EventMessage string  `json:"event_message"`
	EventType    string  `json:"event_type"`
	EventID      string  `json:"event_id"`
	EventStatus  string  `json:"event_status"`
	ContentType  string  `json:"content_type"`
	VisibleScope string  `json:"visible_scope"`
	Outputs      Outputs `json:"outputs"`
}

type TextOutputs struct {
	Arguments     map[string]interface{} `json:"arguments"`
	ComponentCode string                 `json:"component_code"`
	ComponentName string                 `json:"component_name"`
	Text          string                 `json:"text"`
}

type Outputs struct {
	Text json.RawMessage `json:"text"` // 使用 json.RawMessage 来延迟解析
}

type NestedTextOutputs struct {
	TextOutputs
	Text string `json:"text"` // 当 text 是字符串时，这个字段将被填充
}

func parseTextOutputs(rawText json.RawMessage) (*NestedTextOutputs, error) {
	var nested NestedTextOutputs
	// 尝试将 rawText 解析为字符串
	if err := json.Unmarshal(rawText, &nested.Text); err == nil {
		return &nested, nil // 如果成功，说明 text 是一个字符串
	}
	// 如果解析为字符串失败，尝试解析为 TextOutputs 结构体
	if err := json.Unmarshal(rawText, &nested.TextOutputs); err != nil {
		return nil, err // 如果两者都失败，返回错误
	}
	return &nested, nil // 如果解析为结构体成功，返回解析后的结果
}

func SendDialogueContent(context string, conversation_id string) string {

	url := config.AIServiceUrl + "/answer"
	payload := strings.NewReader(
		`{
			"sessionId": "` + conversation_id + `",
			"query": "` + context + `"
		}`)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

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
		Response string `json:"response"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	return response.Response
}
