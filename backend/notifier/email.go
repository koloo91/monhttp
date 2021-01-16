package notifier

import (
	"fmt"
	"github.com/koloo91/monhttp/model"
	"github.com/spf13/viper"
	"net/smtp"
	"strings"
)

type EMailNotifier struct {
	model.Notifier
	Host string
	Port int
	From string
	To   []string
	Auth smtp.Auth
}

func NewEMailNotifier(store *viper.Viper) *EMailNotifier {
	data := make(map[string]interface{})
	data["enabled"] = store.GetBool("NOTIFIER_EMAIL_ENABLED")
	data["host"] = store.GetString("NOTIFIER_EMAIL_HOST")
	data["port"] = store.GetInt("NOTIFIER_EMAIL_PORT")
	data["from"] = store.GetString("NOTIFIER_EMAIL_FROM")
	data["password"] = store.GetString("NOTIFIER_EMAIL_PASSWORD")
	data["to"] = strings.Join(store.GetStringSlice("NOTIFIER_EMAIL_TO"), ",")
	data["password"] = store.GetString("NOTIFIER_EMAIL_PASSWORD")

	data["SERVICE_UP_TEMPLATE"] = store.GetString("NOTIFIER_EMAIL_SERVICE_UP_TEMPLATE")
	if value, exists := data["SERVICE_UP_TEMPLATE"]; !exists || len(value.(string)) == 0 {
		data["SERVICE_UP_TEMPLATE"] = defaultUpTemplate
	}

	data["SERVICE_DOWN_TEMPLATE"] = store.GetString("NOTIFIER_EMAIL_SERVICE_DOWN_TEMPLATE")
	if value, exists := data["SERVICE_DOWN_TEMPLATE"]; !exists || len(value.(string)) == 0 {
		data["SERVICE_DOWN_TEMPLATE"] = defaultDownTemplate
	}

	username := store.GetString("NOTIFIER_EMAIL_FROM")
	password := store.GetString("NOTIFIER_EMAIL_PASSWORD")
	host := store.GetString("NOTIFIER_EMAIL_HOST")
	auth := smtp.PlainAuth("", username, password, host)

	return &EMailNotifier{
		Notifier: model.Notifier{
			Id:      "email",
			Name:    "E-Mail",
			Enabled: viper.GetBool("NOTIFIER_EMAIL_ENABLED"),
			Data:    data,
			Form: []model.NotificationForm{
				{
					Type:            "switch",
					Title:           "Enabled",
					FormControlName: "enabled",
					Placeholder:     "smtp.googlemail.com",
					Required:        false,
				},
				{
					Type:            "text",
					Title:           "Host",
					FormControlName: "host",
					Placeholder:     "smtp.googlemail.com",
					Required:        true,
				},
				{
					Type:            "number",
					Title:           "Port",
					FormControlName: "port",
					Placeholder:     "587",
					Required:        true,
				},
				{
					Type:            "text",
					Title:           "From",
					FormControlName: "from",
					Placeholder:     "gululu@example.com",
					Required:        true,
				},
				{
					Type:            "password",
					Title:           "Password",
					FormControlName: "password",
					Placeholder:     "your string password",
					Required:        true,
				},
				{
					Type:            "text",
					Title:           "To",
					FormControlName: "to",
					Placeholder:     "gululu@example.com,example@example.com",
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
		Host: store.GetString("NOTIFIER_EMAIL_HOST"),
		Port: store.GetInt("NOTIFIER_EMAIL_PORT"),
		From: store.GetString("NOTIFIER_EMAIL_FROM"),
		To:   strings.Split(store.GetString("NOTIFIER_EMAIL_TO"), ","),
		Auth: auth,
	}
}

func (n *EMailNotifier) SendNotification(service model.Service, message string) error {
	return n.send(service, message)
}

func (n *EMailNotifier) send(service model.Service, message string) error {
	header := make(map[string]string)
	header["From"] = n.From
	header["To"] = strings.Join(n.To, ", ")
	header["Subject"] = fmt.Sprintf("monhttp: Service '%s' status notification", service.Name)
	header["MIME-version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	eMailMessage := ""

	for k, v := range header {
		eMailMessage += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	eMailMessage += fmt.Sprintf("\r\n%s", message)

	return smtp.SendMail(fmt.Sprintf("%s:%d", n.Host, n.Port), n.Auth, n.From, n.To, []byte(eMailMessage))
}

func (n *EMailNotifier) GetId() string {
	return n.Id
}

func (n *EMailNotifier) IsEnabled() bool {
	return n.Enabled
}

func (n *EMailNotifier) GetForms() []model.NotificationForm {
	return n.Form
}

func (n *EMailNotifier) GetName() string {
	return n.Name
}

func (n *EMailNotifier) GetData() map[string]interface{} {
	return n.Data
}

func (n *EMailNotifier) GetServiceUpNotificationTemplate() string {
	if data, exists := n.Data["SERVICE_UP_TEMPLATE"]; exists {
		return data.(string)
	}
	return defaultUpTemplate
}

func (n *EMailNotifier) GetServiceDownNotificationTemplate() string {
	if data, exists := n.Data["SERVICE_DOWN_TEMPLATE"]; exists {
		return data.(string)
	}
	return defaultDownTemplate
}
