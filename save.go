package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"
)

var timeFormat = fmt.Sprintf("%s_%s", time.DateOnly, time.TimeOnly)

func latestSave() *MoodSave {
	path := latestPath()

	data, err := os.ReadFile(path)
	if err != nil && err != os.ErrNotExist {
		slog.Error("failed to read latest save", "err", err, "path", path)
		return nil
	}

	moodSave := new(MoodSave)

	if err := json.Unmarshal(data, moodSave); err != nil {
		slog.Error("failed to parse latest save", "err", err, "path", path)
		return nil
	}
	
	return moodSave
}

func save(dto MoodDTO) {
	save := MoodSave{
		Time: time.Now(),
		Entries: dto,
	}

	saveData, err := json.Marshal(save)
	if err != nil {
		slog.Error("failed to serialize new save", "err", err)
		return
	}

	oldLatest := latestSave()

	if oldLatest != nil {
		metric := existingChangedCount(oldLatest.Entries, dto)
		if metric > 3 {
			slog.Info("metric threshold reached, creating new save", "metric", metric)
			oldLatestPath := fmt.Sprintf("%s/ms_%s.json", saveDir, oldLatest.Time.Format(timeFormat))

			data, err := json.Marshal(oldLatest)
			if err != nil {
				slog.Error("failed to serialize latest save", "err", err)
				return
			}

			if err := os.WriteFile(oldLatestPath, data, 0644); err != nil {
				slog.Error("failed to back up last save", "err", err, "path", oldLatestPath)
				return
			}
		} else {
			slog.Info("overwriting last save")
		}
	}

	path := latestPath()
	slog.Info("saving new data", "path", path)
	if err := os.WriteFile(path, saveData, 0644); err != nil {
		slog.Error("failed to save new data", "err", err, "path", path)
	}
}

func existingChangedCount(oldDto, newDto MoodDTO) (count int) {
	for entryTime, oldValue := range oldDto {
		newValue, ok := newDto[entryTime]
		if !ok || newValue != oldValue {
			count++
		}
	}
	return
}

func latestPath() string {
	return fmt.Sprintf("%s/ms_latest.json", saveDir)
}
