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

func (n *EMailNotifier) GetId() string {
	return n.Id
}

func NewEMailNotifier() *EMailNotifier {

	username := viper.GetString("notifier.email.from")
	password := viper.GetString("notifier.email.password")
	host := viper.GetString("notifier.email.host")
	auth := smtp.PlainAuth("", username, password, host)

	return &EMailNotifier{
		Notifier: model.Notifier{
			Id:      viper.GetString("notifier.email.id"),
			Name:    viper.GetString("notifier.email.name"),
			Enabled: viper.GetBool("notifier.email.enabled"),
			Form: []model.NotificationForm{
				{
					Type:            "text",
					Title:           "Host",
					FormControlName: "host",
					Placeholder:     "Host",
					Required:        true,
				},
				{
					Type:            "number",
					Title:           "Port",
					FormControlName: "port",
					Placeholder:     "Port",
					Required:        true,
				},
				{
					Type:            "text",
					Title:           "From",
					FormControlName: "from",
					Placeholder:     "From",
					Required:        true,
				},
				{
					Type:            "password",
					Title:           "Password",
					FormControlName: "password",
					Placeholder:     "Password",
					Required:        true,
				},
				{
					Type:            "text",
					Title:           "To",
					FormControlName: "to",
					Placeholder:     "To",
					Required:        true,
				},
			},
		},
		Host: viper.GetString("notifier.email.host"),
		Port: viper.GetInt("notifier.email.port"),
		From: viper.GetString("notifier.email.from"),
		To:   viper.GetStringSlice("notifier.email.to"),
		Auth: auth,
	}
}

func (n *EMailNotifier) SendSuccess(service model.Service) error {
	message := fmt.Sprintf("Service '%s' is up again", service.Name)
	return n.send(service, message)
}

func (n *EMailNotifier) SendFailure(service model.Service, failure model.Failure) error {
	message := fmt.Sprintf("Service '%s' is down.\nReason: %s", service.Name, failure.Reason)
	return n.send(service, message)
}

func (n *EMailNotifier) send(service model.Service, message string) error {
	header := make(map[string]string)
	header["From"] = n.From
	header["To"] = strings.Join(n.To, ", ")
	header["Subject"] = fmt.Sprintf("monhttp: Service '%s' is down", service.Name)
	header["Content-Type"] = "text/plain; charset=\"utf-8\""

	eMailMessage := ""

	for k, v := range header {
		eMailMessage += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	eMailMessage += fmt.Sprintf("\r\n%s", message)

	return smtp.SendMail(fmt.Sprintf("%s:%d", n.Host, n.Port), n.Auth, n.From, n.To, []byte(eMailMessage))
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
