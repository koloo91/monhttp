package notifier

import (
	"fmt"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TelegramNotifier struct {
	model.Notifier
	ApiToken string
	Channel  string
}

func (n *TelegramNotifier) GetId() string {
	return n.Id
}

func NewTelegramNotifier() *TelegramNotifier {
	data := make(map[string]interface{})
	data["enabled"] = viper.GetBool("notifier.telegram.enabled")
	data["apiToken"] = viper.GetString("notifier.telegram.apiToken")
	data["channel"] = viper.GetString("notifier.telegram.channel")

	return &TelegramNotifier{
		Notifier: model.Notifier{
			Id:      "telegram",
			Name:    "Telegram Notifier",
			Enabled: viper.GetBool("notifier.telegram.enabled"),
			Data:    data,
			Form: []model.NotificationForm{
				{
					Type:            "switch",
					Title:           "Enabled",
					FormControlName: "enabled",
					Placeholder:     "Enabled",
					Required:        false,
				},
				{
					Type:            "text",
					Title:           "Api Token",
					FormControlName: "apiToken",
					Placeholder:     "1404238259:AAEpdJO7az2h6s0K2WPGHSOoGAoTqKdwefrd",
					Required:        true,
				},
				{
					Type:            "text",
					Title:           "Channel",
					FormControlName: "channel",
					Placeholder:     "710144235",
					Required:        true,
				},
			},
		},
		ApiToken: viper.GetString("notifier.telegram.apiToken"),
		Channel:  viper.GetString("notifier.telegram.channel"),
	}
}

func (n *TelegramNotifier) SendSuccess(service model.Service) error {
	message := fmt.Sprintf("Service '%s' is up again", service.Name)
	return n.send(message)
}

func (n *TelegramNotifier) SendFailure(service model.Service, failure model.Failure) error {
	message := fmt.Sprintf("Service '%s' is down.\nReason: %s", service.Name, failure.Reason)
	return n.send(message)
}

func (n *TelegramNotifier) send(message string) error {

	v := url.Values{}
	v.Set("chat_id", n.Channel)
	v.Set("text", message)
	v.Encode()

	apiEndpoint := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage?%s", n.ApiToken, v.Encode())

	response, err := http.Get(apiEndpoint)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	contentByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	log.Println(string(contentByte))

	return nil
}

func (n *TelegramNotifier) IsEnabled() bool {
	return n.Enabled
}

func (n *TelegramNotifier) GetForms() []model.NotificationForm {
	return n.Form
}

func (n *TelegramNotifier) GetName() string {
	return n.Name
}

func (n *TelegramNotifier) GetData() map[string]interface{} {
	return n.Data
}
