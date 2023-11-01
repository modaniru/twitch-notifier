package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/modaniru/streamer-notifier-telegram/internal/entity"
)

type TwitchClient struct {
	client             *http.Client
	twitchClientId     string
	twitchClientSecret string
	token              string
}

var (
	ErrStreamerNotFound = errors.New("streamer was not found")
)

func NewTwitchClient(client *http.Client, twitchClientId, twitchClientSecret string) *TwitchClient {
	return &TwitchClient{client: client, twitchClientId: twitchClientId, twitchClientSecret: twitchClientSecret}
}

func (t *TwitchClient) GetUserIdByLogin(login string) (string, error) {
	login = strings.ToLower(login)
	token, err := t.GetToken()
	if err != nil {
		return "", err
	}
	uri := fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", login)
	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Client-Id", t.twitchClientId)
	resp, err := t.client.Do(request)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("request status code not 200")
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var userInfo entity.UserCollection
	err = json.Unmarshal(b, &userInfo)
	if err != nil {
		return "", err
	}
	if len(userInfo.Data) == 0 {
		return "", ErrStreamerNotFound
	}
	return userInfo.Data[0].Id, nil
}

func (t *TwitchClient) GetToken() (string, error) {
	ok, err := t.ValidateToken(t.token)
	if err != nil {
		return "", err
	}
	if ok {
		return t.token, nil
	}

	uri := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", t.twitchClientId, t.twitchClientSecret)
	request, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		return "", err
	}
	resp, err := t.client.Do(request)
	if err != nil {
		return "", err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	type tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	var token tokenResp
	err = json.Unmarshal(b, &token)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func (t *TwitchClient) ValidateToken(token string) (bool, error) {
	request, err := http.NewRequest(http.MethodGet, "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		return false, err
	}
	request.Header.Add("Authorization", "OAuth "+token)
	res, err := t.client.Do(request)
	if err != nil {
		return false, err
	}
	return res.StatusCode == 200, nil
}

type RegisterWebhookParams struct {
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Condition Condition `json:"condition"`
	Transport Transport `json:"transport"`
}

type Condition struct {
	BroadcasterId string `json:"broadcaster_user_id"`
}

type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

func (t *TwitchClient) RegisterStreamWebhook(callback string, userId string) error {
	condition := Condition{BroadcasterId: userId}
	transport := Transport{Method: "webhook", Callback: callback, Secret: "secretsecret"}
	requestModel := RegisterWebhookParams{Type: "stream.online", Version: "1", Transport: transport, Condition: condition}
	byte, err := json.Marshal(&requestModel)
	if err != nil {
		return err
	}
	token, err := t.GetToken()
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, "https://api.twitch.tv/helix/eventsub/subscriptions", bytes.NewReader(byte))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Client-Id", t.twitchClientId)
	request.Header.Set("Content-Type", "application/json")

	response, err := t.client.Do(request)
	fmt.Println(response.StatusCode)
	b, _ := io.ReadAll(response.Body)
	if response.StatusCode != 202 {
		return errors.New("status code is not 202 " + string(b))
	}
	return nil
}

func (t *TwitchClient) GetUsersByUserId(ids []string) ([]entity.UserInfo, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("ids size is zero")
	}
	token, err := t.GetToken()
	if err != nil {
		return nil, err
	}
	uri := "https://api.twitch.tv/helix/users?id=" + ids[0]
	for _, id := range ids[1:] {
		uri += fmt.Sprintf("&id=%s", id)
	}
	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Client-Id", t.twitchClientId)
	resp, err := t.client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("request status code not 200")
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var userInfo entity.UserCollection
	err = json.Unmarshal(b, &userInfo)
	if err != nil {
		return nil, err
	}
	if len(userInfo.Data) == 0 {
		return nil, ErrStreamerNotFound
	}
	return userInfo.Data, nil
}
