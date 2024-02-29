package vector_db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ZillisOp struct {
	collectionName string
	zillisAPIKey   string
	endpoint       string
}

type addRequestBody struct {
	CollectionName string `json:"collectionName"`
	Data           struct {
		Vector []float32 `json:"doc_vector"`
		UserID int64     `json:"user_id"`
	} `json:"data"`
}

type addResponseBody struct {
	Code int `json:"code"`
	Data struct {
		InsertCount int      `json:"insertCount"`
		InsertIds   []string `json:"insertIds"`
	} `json:"data"`
}

func (zillisOp *ZillisOp) AddToDb(userId int64, docVector []float32) error {

	// Initialize the request body
	body := addRequestBody{
		CollectionName: zillisOp.collectionName,
		Data: struct {
			Vector []float32 `json:"doc_vector"`
			UserID int64     `json:"user_id"`
		}{
			Vector: docVector,
			UserID: userId,
		},
	}

	// Marshal the request body into JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling the request body:", err)
		return err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", zillisOp.endpoint+"insert", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating the request:", err)
		return err
	}

	// Set the necessary headers
	req.Header.Set("Authorization", "Bearer "+zillisOp.zillisAPIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending the request:", err)
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return err
	}

	// Unmarshal the response body into the responseBody struct
	var response addResponseBody
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		fmt.Println("Error unmarshalling the response body:", err)
		return err
	}

	if response.Code != 200 {
		return errors.New("add request to zillis failed")
	}

	return nil
}

type SearchRequest struct {
	CollectionName string    `json:"collectionName"`
	OutputFields   []string  `json:"outputFields"`
	Vector         []float32 `json:"vector"`
	Limit          int       `json:"limit"`
}

type SuccessResponse struct {
	Code int `json:"code"`
	Data []struct {
		Distance float64 `json:"distance"`
		UserID   int64   `json:"user_id"`
	} `json:"data"`
}

// Define a struct for the failure response
type FailureResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (zillisOp *ZillisOp) SearchInDb(queryVector []float32) ([]int64, error) {

	// Dynamic variables
	clusterEndpoint := zillisOp.endpoint // Replace YOUR_CLUSTER_ENDPOINT with your actual endpoint
	yourToken := zillisOp.zillisAPIKey   // Replace YOUR_TOKEN with your actual token

	searchReq := SearchRequest{
		CollectionName: zillisOp.collectionName,
		OutputFields:   []string{"user_id"},
		Vector:         queryVector,
		Limit:          10,
	}

	// Marshal the request body to JSON
	reqBody, err := json.Marshal(searchReq)
	if err != nil {
		fmt.Printf("Error marshaling request body: %v\n", err)
		return nil, err

	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", clusterEndpoint+"search", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err

	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+yourToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return nil, err

	}
	defer resp.Body.Close()

	//Read and print the response body
	// TODO check WEIRD doesn't work when duplicated ReadAll??
	//respBody, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	fmt.Printf("Error reading response body: %v\n", err)
	//	return nil, err
	//}
	//fmt.Println("Response:", string(respBody))

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return nil, err
	}

	// Attempt to unmarshal the response into the successResponse struct
	var sResp SuccessResponse
	err = json.Unmarshal(responseBody, &sResp)
	fmt.Println("sResp.Code=>", sResp.Code)
	if err != nil || sResp.Code != 200 { // Also check if the code is not 200, then it's not a success
		// If there's an error or the code is not 200, try to unmarshal into the failureResponse struct
		var fResp FailureResponse
		err = json.Unmarshal(responseBody, &fResp)
		if err != nil {
			fmt.Println("Error unmarshalling the response body:", err)
			return nil, errors.New(fResp.Message)
		}
		// Handle failure response here
		fmt.Printf("Failed: %d - %s\n", fResp.Code, fResp.Message)
		return nil, err
	} else {
		// Handle success response here
		fmt.Printf("Success: %d\n", sResp.Code)
		var userIDs []int64
		for _, data := range sResp.Data {
			userIDs = append(userIDs, data.UserID)
		}
		return userIDs, err
	}
}

func NewZillisOp(collectionName, zillisAPIKey, endpoint string) VectorOp {
	milvusOp := &ZillisOp{
		collectionName: collectionName,
		zillisAPIKey:   zillisAPIKey,
		endpoint:       endpoint,
	}

	return milvusOp
}
