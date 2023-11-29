package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

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
    values := map[string]string{"audio_url": url}
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

func getComputedTranscriptions(pollingUrl string, pollInterval time.Duration, wg *sync.WaitGroup, transcription chan string){
    defer wg.Done()
    
    API_KEY := config.GetConfig("AAI_API_KEY")
    client := &http.Client{}
    for{
        req, _ := http.NewRequest("POST",pollingUrl, nil)
        req.Header.Set("content-type", "application/json")
        req.Header.Set("authorization", API_KEY)
        res, _ := client.Do(req) 

        var result map[string]string
        json.NewDecoder(res.Body).Decode(&result)

        if result["status"] == "completed" {
            transcription <- result["text"]
        }
        fmt.Println("Polling::", result["error"])
        res.Body.Close()
        time.Sleep(pollInterval)
    }
}


func GetTranscriptionPollingUrl(path string)(string, error){
    uploadUrl, err := uploadFileToAssemblyAI(path);if err != nil{
        return "", err
    }
    transcriptionPollingUrl, err := startComputingTranscriptions(uploadUrl);if err != nil{
        return "", err
    }

    pollInterval := 5 * time.Second  
    resultChan := make(chan string)   
    
    var wg sync.WaitGroup
	wg.Add(1)

    go getComputedTranscriptions(transcriptionPollingUrl, pollInterval, &wg, resultChan)

    go func() {
		wg.Wait()
		close(resultChan)
	}()

    transcriptText := <-resultChan
    
    return transcriptText, nil
}