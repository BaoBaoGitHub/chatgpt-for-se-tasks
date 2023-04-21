package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xyhelper/chatgpt-go"
)

func main() {
	// new chatgpt client
	token := uuid.New().String()
	baseURI := "https://freechat.lidong.xin"

	if testFlag := true; testFlag {
		token = "EBA1C3EB-C3AC-4D1F-B32A-005B07BD6D59"
		baseURI = "https://pluschat.lidong.xin"
	}

	cli := chatgpt.NewClient(
		chatgpt.WithDebug(false),
		chatgpt.WithTimeout(60*time.Second),
		chatgpt.WithAccessToken(token),
		chatgpt.WithBaseURI(baseURI),
		chatgpt.WithModel("gpt-3.5-turbo"),
	)

	conversationID := ""
	parentMessage := ""

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("请输入您的消息：")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("读取输入失败：", err)
			continue
		}
		stream, err := cli.GetChatStream(message, conversationID, parentMessage)
		if err != nil {
			log.Fatalf("get chat stream failed: %v\n", err)
		}

		var answer string
		for text := range stream.Stream {
			// log.Printf("stream text: %s\n", text.Content)
			print(strings.Replace(text.Content, answer, "", 1))

			answer = text.Content
			conversationID = text.ConversationID
			parentMessage = text.MessageID
		}

		if stream.Err != nil {
			log.Fatalf("stream closed with error: %v\n", stream.Err)
		}
		// 输出换行
		println()
	}
}
