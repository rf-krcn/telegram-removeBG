# Telegram Background Removal Bot

This Telegram bot allows users to send images, and it will remove the background and send the edited image back to the user. The bot is implemented using Golang for managing and handling communication, and it leverages a pre-existing Python model for background removal.

## Features

-   **Image Background Removal**: Users can send images to the bot, and it will process the image and remove the background.
## Getting Started
### Prerequisites

-   Telegram Bot Token: Obtain a token by creating a new bot on Telegram via BotFather.
### Installation

1.  Clone the repository:

```
	https://github.com/rf-krcn/telegram-removeBG.git
```

2.  Install the Golang application:

```
	cd bot-service
	go install ./...
```
3. Install Python dependencies:


```
	cd model-service
	pip install -r requirements.txt
``` 

Make sure to have the required dependencies, including torch, installed.

## Usage  
1. Open the `main.go` file in the `bot-service/cmd` directory and replace the placeholder `token` with your actual Telegram bot token: 
```go 
// main.go 
const token = "YOU_BOT_TOKEN"
```
2. Download the model weights from [here](https://drive.usercontent.google.com/download?id=1ao1ovG1Qtx4b7EoskHXmi2E9rp5CHLcZ&authuser=0) and put the file in `./saved_models/u2net/u2net.pth`.

3. Run both Golang and Python services:

-   **Golang Service:**

```bash
	cd bot-service
	go run ./cmd 
```
-   **Python Service:**
```bash
	cd model-service
	python u2net_test.py
```
4.  In your Telegram app, start a chat with the bot.
    
5.  Send an image file to the bot.

6.  The bot will process the image and send back the edited version with the background removed.
## Acknowledgments
This project integrates a background removal model implemented in Python, sourced from [model-repository-link](https://github.com/xuebinqin/U-2-Net) .
