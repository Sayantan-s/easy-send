package openai

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sayantan-s/easy-send/config"
)

const OPEN_AI_COMPLETIONS_URL ="https://api.openai.com/v1/chat/completions"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type  CompletionsBody struct{
	Model string `json:"model"`
	Messages []Message `json:"messages"`
	// Prompt string `json:"prompt"`
	// Stream bool `json:"stream"`
    // Temperature float32 `json:"temperature"`
    // MaxTokens int32 `json:"max_tokens"`
}

func Completions(values string)(string, error){
    // jsonData, err := json.Marshal(values)
    // if err != nil {
    //     return "", err
    // }
	API_KEY := fmt.Sprintf("Bearer %s", config.GetConfig("OPENAI_API_KEY")) 
	fmt.Println(API_KEY)
    client := &http.Client{}
	req, _ := http.NewRequest("POST", OPEN_AI_COMPLETIONS_URL, bytes.NewBufferString(values))
    req.Header.Set("content-type", "application/json")
    req.Header.Set("Authorization", API_KEY)
    res, err := client.Do(req)

    if err != nil {
       return "", err
    }

    defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	result := string(responseBody)

	fmt.Println("Result", result)

	return "", nil
}