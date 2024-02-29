package vectorEmbeddingAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	openAIModel = "text-embedding-3-large"
)

type APIResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

func GetVectorEmbedding(doc, apiKey, url string) ([]float32, error) {

	// Construct the request payload
	payload := fmt.Sprintf(`{"input": "%s", "model": "%s"}`, doc, openAIModel)

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Initialize HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	// Print the response body
	//fmt.Println("Response:", string(responseBody))

	// Parse the JSON response
	var apiResponse APIResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return nil, err
	}

	// Check if the data array contains at least one item
	// TODO check embedding length?
	if len(apiResponse.Data) > 0 && len(apiResponse.Data[0].Embedding) > 0 {
		// If there's at least one item and the first item's Embedding array is not empty
	} else {
		// If there are no items or the first item's Embedding array is empty
		fmt.Println("No data found in response")
		fmt.Println("Response:", string(responseBody))
		fmt.Println("Payload:", payload)
		return nil, errors.New("No data found in response")
	}

	return apiResponse.Data[0].Embedding, nil
}
