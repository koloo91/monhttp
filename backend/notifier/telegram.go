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

func NewTelegramNotifier(store *viper.Viper) *TelegramNotifier {
	data := make(map[string]interface{})
	data["enabled"] = store.GetBool("NOTIFIER_TELEGRAM_ENABLED")
	data["apiToken"] = store.GetString("NOTIFIER_TELEGRAM_APITOKEN")
	data["channel"] = store.GetString("NOTIFIER_TELEGRAM_CHANNEL")

	data["SERVICE_UP_TEMPLATE"] = store.GetString("NOTIFIER_TELEGRAM_SERVICE_UP_TEMPLATE")
	if value, exists := data["SERVICE_UP_TEMPLATE"]; !exists || len(value.(string)) == 0 {
		data["SERVICE_UP_TEMPLATE"] = defaultUpTemplate
	}

	data["SERVICE_DOWN_TEMPLATE"] = store.GetString("NOTIFIER_TELEGRAM_SERVICE_DOWN_TEMPLATE")
	if value, exists := data["SERVICE_DOWN_TEMPLATE"]; !exists || len(value.(string)) == 0 {
		data["SERVICE_DOWN_TEMPLATE"] = defaultDownTemplate
	}

	return &TelegramNotifier{
		Notifier: model.Notifier{
			Id:      "telegram",
			Name:    "Telegram Notifier",
			Enabled: viper.GetBool("NOTIFIER_TELEGRAM_ENABLED"),
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
				{
					Type:            "textarea",
					Title:           "Up template",
					FormControlName: "SERVICE_UP_TEMPLATE",
					Placeholder:     "Service {{.Name}} is up",
					Required:        true,
				},
				{
					Type:            "textarea",
					Title:           "Down template",
					FormControlName: "SERVICE_DOWN_TEMPLATE",
					Placeholder:     "Service {{.Name}} is down",
					Required:        true,
				},
			},
		},
		ApiToken: store.GetString("NOTIFIER_TELEGRAM_APITOKEN"),
		Channel:  store.GetString("NOTIFIER_TELEGRAM_CHANNEL"),
	}
}

func (n *TelegramNotifier) SendNotification(service model.Service, message string) error {
	return n.send(message)
}

func (n *TelegramNotifier) send(message string) error {

	v := url.Values{}
	v.Set("chat_id", n.Channel)
	v.Set("text", message)
	v.Encode()

	apiEndpoint := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage?%s&parse_mode=HTML", n.ApiToken, v.Encode())

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

func (n *TelegramNotifier) GetId() string {
	return n.Id
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

func (n *TelegramNotifier) GetServiceUpNotificationTemplate() string {
	if data, exists := n.Data["SERVICE_UP_TEMPLATE"]; exists {
		return data.(string)
	}
	return "Service <b>'{{.Name}}'</b> is up again!"
}

func (n *TelegramNotifier) GetServiceDownNotificationTemplate() string {
	if data, exists := n.Data["SERVICE_DOWN_TEMPLATE"]; exists {
		return data.(string)
	}
	return "Service <b>'{{.Name}}'</b> is down. Reason: '{{.Reason}}' at {{.Date}}"
}
