package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

func fetchMoods(passwd string) []byte {
	slog.Info("fetching mood list", "url", baseUrl)

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			slog.Info("not following redirect", "url", loginUrl)
			return http.ErrUseLastResponse
    },
	}

	slog.Info("making request", "method", "POST", "url", loginUrl)
	loginForm := url.Values{}
	loginForm.Set("password", passwd)
	resp, err := client.PostForm(loginUrl, loginForm)
	if err != nil {
		slog.Error("failed to login", "err", err)
		return nil
	}
	cookies := resp.Cookies()

	slog.Info("making request", "method", "GET", "url", moodUrl)
	req, err := http.NewRequest("GET", moodUrl, nil)
	if err != nil {
		slog.Error("failed to create request", "err", err)
		return nil
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err = client.Do(req)
	if err != nil {
		slog.Error("failed to fetch mood", "err", err)
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read body", "err", err)
		return nil
	}

	return data
}
