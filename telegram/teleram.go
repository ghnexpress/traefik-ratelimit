package telegram

import (
	"fmt"
	"github.com/ghnexpress/traefik-ratelimit/log"
	"github.com/ghnexpress/traefik-ratelimit/utils"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TelegramService struct {
	host   string
	token  string
	chatID string
}

func NewTelegramService(host, token, chatID string) *TelegramService {
	return &TelegramService{host: host, token: token, chatID: chatID}
}

func (s *TelegramService) SendError(errToSend error) error {
	log.Log("send err")
	sendMessagePath := "sendMessage"
	apiUrl, err := utils.GetUrl(s.host, s.token, sendMessagePath)
	if err != nil {
		return err
	}
	queries := url.Values{}
	queries.Add("chat_id", s.chatID)
	queries.Add("text", errToSend.Error())
	queries.Add("parse_mode", "HTML")

	apiUrl.RawQuery = queries.Encode()
	log.Log(apiUrl.String())
	res, err := http.Get(apiUrl.String())
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		log.Log(res.StatusCode, res.Body)
		body, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return err
		}
		return fmt.Errorf(string(body))
	}
	return nil
}
