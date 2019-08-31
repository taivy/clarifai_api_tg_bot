package clarifai_api

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"strings"
)

const CLARIFAI_TOKEN = "a5264ebf156a4cdc92b5d36f4f229e52"

type ClarifaiResponse struct {
    Status    StatusFields   `json:"status,omitempty"`
    Outputs []Output `json:"outputs,omitempty"`
}

type StatusFields struct {
	Code int `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

type Output struct {
    Data DataField `json:"data,omitempty"`
}

type DataField struct {
    Concepts []Concept `json:"concepts,omitempty"`
}

type Concept struct {
    Id string `json:"id,omitempty"`
    Name string `json:"name,omitempty"`
    Value float32 `json:"value,omitempty"`
}

var Language string = "en"

func GetClarifaiResp(imageUrl string) string {
	url := "https://api.clarifai.com/v2/models/aaa03c23b3724a16a56b629203edc62c/versions/aa7f35c01e0642fda5cf400f543e7c40/outputs"
    var jsonStr = []byte(fmt.Sprintf(
`{
  "inputs": [
    {
      "data": {
        "image": {
          "url": "%s"
        }
      }
    }
  ],
  "model":{
    "output_info":{
      "output_config":{
        "language":"%s"
      }
    }
  }
}
`, imageUrl, Language))
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Key a5264ebf156a4cdc92b5d36f4f229e52")
    
	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
	    panic(e)
	}
	
	defer resp.Body.Close()
	
	body, _ := ioutil.ReadAll(resp.Body)
	
	var apiResp ClarifaiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
	  panic(err)
	}
		
	if apiResp.Status.Code != 10000 {
		return fmt.Sprintf("Api response isn't ok: %s, %s", strconv.Itoa(apiResp.Status.Code), apiResp.Status.Description)
	}
	
	conceptsLabels := make([]string, 0)
	for _, concept := range apiResp.Outputs[0].Data.Concepts {
		conceptsLabels = append(conceptsLabels, concept.Name)
	}
	
	return strings.Join(conceptsLabels, "\n")
}

