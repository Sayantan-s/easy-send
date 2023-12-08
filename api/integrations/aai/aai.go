package aai

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
        "webhook_url": "https://e655-2409-40f2-1027-653a-311b-9546-64e5-682b.ngrok.io/api/generate/transcript_CE_webhook",
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


    return transcriptionId, nil

}

func getComputedTranscriptions(transcriptId string)(map[string]string, error){
    
    AAI_TRANSCRIPT_URL :=  TRANSCRIPT_URL + "/" + transcriptId
    API_KEY := config.GetConfig("AAI_API_KEY")
    
    client := &http.Client{}
    req, _ := http.NewRequest("GET", AAI_TRANSCRIPT_URL, nil)
    req.Header.Set("content-type", "application/json")
    req.Header.Set("authorization", API_KEY)
    
    res, err := client.Do(req)

    if err != nil {
       return nil, err
    }

    defer res.Body.Close()

    if err != nil {
       return nil, err
    }

    var result map[string]string
    json.NewDecoder(res.Body).Decode(&result)

    return result, nil

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

func FetchTranscriptions(transcriptId string)(string, string, error){
    transcriptData, err := getComputedTranscriptions(transcriptId);if err != nil{
        return "", "", err
    }
    
    transcript := transcriptData["text"]
    audio_url := transcriptData["audio_url"]

    return transcript, audio_url, nil;
}