package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func ProcessImage(imageContent []byte) ([]byte, error) {
	modelServiceURL := "http://localhost:8080/process_image"

	resp, err := http.Post(modelServiceURL, "application/octet-stream", bytes.NewBuffer(imageContent))
	if err != nil {
		log.Println("Failed to make request to the model service:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Model service returned non-OK status: %v", resp.Status)
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read model service response:", err)
		return nil, err
	}

	return result, nil
}
