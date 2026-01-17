package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tonysanin/brobar/pkg/helpers"
	"github.com/tonysanin/brobar/pkg/response"
	"io"
	"mime/multipart"
	"net/http"
)

func (s *ProductService) uploadFileToFileService(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	fileServiceURL :=
		helpers.GetEnv("FILE_HOST", "http://file-service") +
			":" +
			helpers.GetEnv("FILE_PORT", "3001") +
			"/files/upload"
	req, err := http.NewRequest("POST", fileServiceURL, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to file-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("file-service error: %s", respBody)
	}

	var fsResp response.FileServiceResponse
	err = json.NewDecoder(resp.Body).Decode(&fsResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse file-service response: %w", err)
	}

	if !fsResp.Success {
		return "", fmt.Errorf("file-service returned failure")
	}

	return fsResp.Data.Filename, nil
}
