package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
)

type MoodDTO map[string]int
type MoodSave struct {
	Time time.Time `json:"time"`
	Entries MoodDTO `json:"entries"`
}

var baseUrl string
var loginUrl string
var moodUrl string

var saveDir string

var interval time.Duration

func main() {
	passwd := os.Getenv("ADM_PASSWD")
	baseUrl = strings.TrimSuffix(os.Getenv("BASE_URL"), "/")

	saveDir = strings.TrimSuffix(os.Getenv("SAVE_DIR"), "/")

	loginUrl = fmt.Sprintf("%s/login", baseUrl)
	moodUrl = fmt.Sprintf("%s/adm/mood.json", baseUrl)

	var err error
	if interval, err = time.ParseDuration(os.Getenv("INTERVAL")); err != nil {
		interval = 24 * time.Hour
	}

	slog.Info("scheduling mood backups", "from", baseUrl, "to", saveDir, "every", interval)

	ticker := time.NewTicker(interval)

	for {
		backup(passwd)
		<-ticker.C
	}
}

func backup(passwd string) {
	slog.Info("backing up latest data")

	data := fetchMoods(passwd)
	if data == nil {
		return
	}

	var dto MoodDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		slog.Error("failed to parse dto", "err", err)
		return
	}

	save(dto)
}
