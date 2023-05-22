package core

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/Jx2f/ViaGenshin/pkg/logger"
)

const (
	consoleUid         = uint32(1)
	consoleNickname    = "djdoolky76"
	consoleLevel       = uint32(60)
	consoleWorldLevel  = uint32(8)
	consoleSignature   = "Welcome to Anjocally Server"
	consoleNameCardId  = uint32(210001)
	consoleAvatarId    = uint32(10000077)
	consoleCostumeId   = uint32(0)
	consoleWelcomeText = "Welcome to Anjocally. \nThis is experimental server, Major issues may arise."
)

type MuipResponseBody struct {
	Retcode int32  `json:"retcode"`
	Msg     string `json:"msg"`
	Ticket  string `json:"ticket"`
	Data    struct {
		Msg    string `json:"msg"`
		Retmsg string `json:"retmsg"`
	} `json:"data"`
}

func (s *Server) ConsoleExecute(cmd, uid uint32, text string) (string, error) {
	logger.Info().Uint32("uid", uid).Msgf("Console Execution: %s", text)
	var values []string
	values = append(values, fmt.Sprintf("cmd=%d", cmd))
	values = append(values, fmt.Sprintf("uid=%d", uid))
	values = append(values, fmt.Sprintf("msg=%s", text))
	values = append(values, fmt.Sprintf("region=%s", s.config.Console.MuipRegion))
	ticket := make([]byte, 16)
	if _, err := rand.Read(ticket); err != nil {
		return "", fmt.Errorf("Unable to process your ticket: %w", err)
	}
	values = append(values, fmt.Sprintf("ticket=%x", ticket))
	if s.config.Console.MuipSign != "" {
		shaSum := sha256.New()
		sort.Strings(values)
		shaSum.Write([]byte(strings.Join(values, "&") + s.config.Console.MuipSign))
		values = append(values, fmt.Sprintf("sign=%x", shaSum.Sum(nil)))
	}
	uri := s.config.Console.MuipEndpoint + "?" + strings.ReplaceAll(strings.Join(values, "&"), " ", "+")
	logger.Debug().Msgf("MUIP Response: %s", uri)
	resp, err := http.Get(uri)
	if err != nil {
		return "Muip Response: %s" + "\nWarning: reminder that the sending options may be limited." + consoleWelcomeText, err
	}
	defer resp.Body.Close()
	p, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	logger.Debug().Msgf("MUIP Response: %s", string(p))
	if resp.StatusCode != 200 {
		return "MUIP Response: %s" + "\nWarning: reminder that the sending options may be limited." + consoleWelcomeText, fmt.Errorf("Status Code: %d", resp.StatusCode)
	}
	body := new(MuipResponseBody)
	if err := json.Unmarshal(p, body); err != nil {
		return "", err
	}
	if (text == "help") {
		return "To execute command, type gm command here.", nil
	}
	if body.Retcode != 0 {
		return "Failed to execute command: " + body.Data.Msg + ", Make sure your keyword command is correct: " + body.Msg + "\nMessage:"+ consoleWelcomeText, nil
	}
	return "The command was executed successfully: " + body.Data.Msg + "\nMessage:" + consoleWelcomeText, nil
}
