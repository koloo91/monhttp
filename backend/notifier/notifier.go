package notifier

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/smtp"
	"strings"
)

type NotificationSystem struct {
	notifier          []Notify
	notificationQueue chan Notification
}

func NewNotificationSystem() *NotificationSystem {
	return &NotificationSystem{
		notifier:          make([]Notify, 0),
		notificationQueue: make(chan Notification, 100),
	}
}

func (n *NotificationSystem) SetupDefaultNotifier() {
	n.notifier = append(n.notifier, NewEMailNotifier())
}

func (n *NotificationSystem) Start() {
	go func() {
		for notification := range n.notificationQueue {
			log.Println("Sending notification: ", notification)
		}
	}()
}

func (n *NotificationSystem) AddNotification(notification Notification) {
	go func() {
		// add notification non blocking
		n.notificationQueue <- notification
	}()
}

type Notification struct {
	Id      string
	Service model.Service
	Failure model.Failure
}

func NewNotification(service model.Service, failure model.Failure) Notification {
	return Notification{
		Id:      uuid.New().String(),
		Service: service,
		Failure: failure,
	}
}

type Notify interface {
	SendSuccess(model.Service) error
	SendFailure(model.Service, model.Failure) error
}

type Notifier struct {
	Id      string
	Name    string
	Enabled bool
}

type EMailNotifier struct {
	Notifier
	Host string
	Port int
	From string
	To   []string
	Auth smtp.Auth
}

func NewEMailNotifier() *EMailNotifier {

	username := viper.GetString("notifier.email.from")
	password := viper.GetString("notifier.email.password")
	host := viper.GetString("notifier.email.host")
	auth := smtp.PlainAuth("", username, password, host)

	return &EMailNotifier{
		Notifier: Notifier{
			Id:      viper.GetString("notifier.email.id"),
			Name:    viper.GetString("notifier.email.name"),
			Enabled: viper.GetBool("notifier.email.enabled"),
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
	return n.send(message)
}

func (n *EMailNotifier) SendFailure(service model.Service, failure model.Failure) error {
	message := fmt.Sprintf("Service '%s' is down.\nReason: %s", service.Name, failure.Reason)
	return n.send(message)
}

func (n *EMailNotifier) send(message string) error {
	message = fmt.Sprintf(`From:%s\n
								  To:%s\n
								  Subject:%s\n\n
  								  %s`,
		n.From, strings.Join(n.To, ", "), "Service fucked up. monhttp", message)
	return smtp.SendMail(fmt.Sprintf("%s:%d", n.Host, n.Port), n.Auth, n.From, n.To, []byte(message))
}
