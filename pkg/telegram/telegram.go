package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TelegramClient struct {
	BotToken string
	ChatID   string
	Client   *http.Client
}

type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

func NewTelegramClient(botToken, chatID string) *TelegramClient {
	return &TelegramClient{
		BotToken: botToken,
		ChatID:   chatID,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (t *TelegramClient) SendMessage(message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)

	payload := TelegramMessage{
		ChatID:    t.ChatID,
		Text:      message,
		ParseMode: "HTML",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal telegram message: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.Client.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to send telegram message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram API returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
