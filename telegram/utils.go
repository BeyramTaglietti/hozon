package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func createMultiPartForm(tgRequest telegramDocumentRequest) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	writer.WriteField("chat_id", tgRequest.ChatId)

	part, err := writer.CreateFormFile("document", tgRequest.DocumentPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form file: %w", err)
	}

	file, err := os.Open(tgRequest.DocumentPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Copy the file content into the form file
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to copy file content: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("failed to close writer: %w", err)
	}

	return &buf, writer.FormDataContentType(), nil
}

func validateTelegramRespose(resp *http.Response) (bool, error) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	var TGResponse TelegramResponse
	if err := json.Unmarshal(respBody, &TGResponse); err != nil {
		return false, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	if !TGResponse.Ok {
		return false, fmt.Errorf("telegram response not ok: %s", TGResponse.Description)
	}

	return TGResponse.Ok, nil
}
