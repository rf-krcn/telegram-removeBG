package handlers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleTelegramUpdates(bot *tgbotapi.BotAPI) {
	updates, err := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go HandleTelegramUpdate(bot, update)
	}
}

func HandleTelegramUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			SendMessage(bot, update.Message.Chat.ID, "Welcome! Send me an image as a file.")
		case "help":
			SendMessage(bot, update.Message.Chat.ID, "I can help you remove the background from images. Just send me an image.")
		}
	} else if update.Message.Document != nil {

		fileID := update.Message.Document.FileID

		if !isImage(update.Message.Document.FileName) {
			SendMessage(bot, update.Message.Chat.ID, "The file should be an image!")
			return
		}
		file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
		if err != nil {
			SendMessage(bot, update.Message.Chat.ID, "Failed to get file information! Try another file, if still not working try later.")
			return
		}

		message := SendMessage(bot, update.Message.Chat.ID, "Processing...")

		fileURL := file.Link(bot.Token)

		response, err := http.Get(fileURL)
		if err != nil {
			bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: update.Message.Chat.ID, MessageID: message.MessageID})
			SendMessage(bot, update.Message.Chat.ID, "Failed to download image! Try another file, if still not working try later.")
			return
		}
		defer response.Body.Close()

		imageBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: update.Message.Chat.ID, MessageID: message.MessageID})
			SendMessage(bot, update.Message.Chat.ID, "Failed to read image content!")
			return
		}

		processedImage, err := ProcessImage(imageBytes)
		if err != nil {
			bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: update.Message.Chat.ID, MessageID: message.MessageID})
			SendMessage(bot, update.Message.Chat.ID, "Failed to process image!")
			return
		}

		fileName, err := saveFile(processedImage)
		if err != nil {
			bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: update.Message.Chat.ID, MessageID: message.MessageID})
			SendMessage(bot, update.Message.Chat.ID, "Internal error! Try again later.")
			return
		}

		err = SendPhoto(bot.Token, update.Message.Chat.ID, fileName, update.Message.Document.FileName)
		bot.DeleteMessage(tgbotapi.DeleteMessageConfig{ChatID: update.Message.Chat.ID, MessageID: message.MessageID})
		if err != nil {
			SendMessage(bot, update.Message.Chat.ID, "Internal error! Try again later.")
		}
		_ = os.Remove(fileName)
	} else if update.Message.Photo != nil {
		SendMessage(bot, update.Message.Chat.ID, "The image must be sent as a file to preserve the size!")
	}
}

func SendPhoto(token string, chatID int64, path, fileName string) error {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", token)

	photoFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer photoFile.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("document", addSuffixBeforeExtension(fileName, "_no_BG"))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, photoFile)
	if err != nil {
		return err
	}

	writer.WriteField("chat_id", strconv.FormatInt(chatID, 10))
	writer.WriteField("caption", "Background removed successfully!")

	writer.Close()

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatID, text)
	message, err := bot.Send(msg)
	if err != nil {
		log.Println("Failed to send message:", err)
		return message
	}
	return message
}
