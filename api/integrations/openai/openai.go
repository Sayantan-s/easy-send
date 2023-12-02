package openai

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sayantan-s/easy-send/config"
)

const UPLOAD_URL = "https://api.assemblyai.com/v2/upload"
const TRANSCRIPT_URL = "https://api.assemblyai.com/v2/transcript"


func uploadFileToAssemblyAI(path string) (string,error){
    data, err := ioutil.ReadFile(path)
    API_KEY := config.GetConfig("AAI_API_KEY")

    if err != nil {
        return "", err
    }
    
    client := &http.Client{}
    req, _ := http.NewRequest("POST", UPLOAD_URL, bytes.NewBuffer(data))
    req.Header.Set("authorization", API_KEY)
    res, err := client.Do(req)

    if err != nil {
        return "", err
    }

    var result map[string]string
    json.NewDecoder(res.Body).Decode(&result)

    return result["upload_url"], nil
}

func startComputingTranscriptions(url string)(string, error){
    values := map[string]string{
        "audio_url": url,
        "webhook_url": "https://c7f0-2409-40f2-1008-7132-c5a8-a577-96f5-daa8.ngrok.io/transcript_CE_webhook",
    }
    jsonData, err := json.Marshal(values)
    if err != nil {
        return "", err
    }
    API_KEY := config.GetConfig("AAI_API_KEY")
    client := &http.Client{}
    req, _ := http.NewRequest("POST", TRANSCRIPT_URL, bytes.NewBuffer(jsonData))
    req.Header.Set("content-type", "application/json")
    req.Header.Set("authorization", API_KEY)
    res, err := client.Do(req)

    if err != nil {
       return "", err
    }

    defer res.Body.Close()

    if err != nil {
       return "", err
    }

    var result map[string]string
    json.NewDecoder(res.Body).Decode(&result)

    transcriptionId := result["id"]

    POLLING_URL := TRANSCRIPT_URL + "/" + transcriptionId

    return POLLING_URL, nil

}

func SetUpTranscriptions(path string)(string, error){
    uploadUrl, err := uploadFileToAssemblyAI(path);if err != nil{
        return "", err
    }
    transcriptionPollingUrl, err := startComputingTranscriptions(uploadUrl);if err != nil{
        return "", err
    }
    
    return transcriptionPollingUrl, nil
}