package telegram

type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

type telegramTextRequest struct {
	ChatId    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func CreateTelegramTextRequest(chatId, text string) telegramTextRequest {
	return telegramTextRequest{
		ChatId:    chatId,
		Text:      text,
		ParseMode: "Markdown",
	}
}

type telegramDocumentRequest struct {
	ChatId       string
	DocumentPath string
}

func CreateTelegramDocumentRequest(chatId, document string) telegramDocumentRequest {
	return telegramDocumentRequest{
		ChatId:       chatId,
		DocumentPath: document,
	}
}

type TelegramSettings struct {
	TGBotToken string
	TGChatID   string
}
