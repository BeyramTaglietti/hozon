package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	messageContentType   = "application/json"
	baseURL              = "https://api.telegram.org/bot%s/%s"
	sendMessageEndpoint  = "sendMessage"
	sendDocumentEndpoint = "sendDocument"
)

func SendMessage(token string, TgRequest telegramTextRequest) {
	ok, err := sendMessageRequest(token, TgRequest)
	if err != nil || !ok {
		log.Fatalf("Failed to send text request: %v", err)
	}
}

func SendFile(token string, TgRequest telegramDocumentRequest) {
	ok, err := sendFileRequest(token, TgRequest)
	if err != nil || !ok {
		log.Fatalf("Failed to send document request: %v", err)
	}
}

func SendGreeting(token string, chatId string) {
	SendMessage(token,
		CreateTelegramTextRequest(
			chatId,
			`*Hi there* 👋

Thanks for using *Hozon!* 🚀

Really hope you find it useful.
Visit my [website](https://beyram.dev) to learn more about me, maybe you'll find more useful stuff you can use.

*Enjoy!*
		`,
		))
}

func sendMessageRequest(token string, request telegramTextRequest) (bool, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	body, err := json.Marshal(request)
	if err != nil {
		return false, fmt.Errorf("failed to marshal message body: %v", err)
	}

	resp, err := http.Post(url, messageContentType, bytes.NewBuffer(body))

	if err != nil {
		return false, fmt.Errorf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	return validateTelegramRespose(resp)
}

func sendFileRequest(token string, request telegramDocumentRequest) (bool, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", token)

	body, contentType, err := createMultiPartForm(request)
	if err != nil {
		return false, fmt.Errorf("failed to create multipart form: %v", err)
	}

	resp, err := http.Post(url, contentType, body)

	if err != nil {
		return false, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	return validateTelegramRespose(resp)
}
