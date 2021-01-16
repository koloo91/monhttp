package notifier

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"html/template"
	"time"
)

const (
	defaultUpTemplate   = "Service <b>'{{.Name}}'</b> is up again!"
	defaultDownTemplate = "Service <b>'{{.Name}}'</b> is down. Reason: '{{.Reason}}' at {{.Date}}"

	globalNotifierId = "global"
)

type NotificationSystem struct {
	notifiers         []model.Notify
	notificationQueue chan Notification
}

func NewNotificationSystem() *NotificationSystem {
	return &NotificationSystem{
		notifiers:         make([]model.Notify, 0),
		notificationQueue: make(chan Notification, 1024),
	}
}

func (n *NotificationSystem) SetupDefaultNotifier() {
	log.Info("Setting up default notifiers")
	n.notifiers = make([]model.Notify, 0)

	load := func() {
		log.Info("Adding notifiers")
		n.notifiers = append(n.notifiers, NewEMailNotifier(viper.GetViper()))
		n.notifiers = append(n.notifiers, NewTelegramNotifier(viper.GetViper()))
	}
	load()

	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Info("Config file changed: ", in.Name)
		// reset notifiers
		n.notifiers = make([]model.Notify, 0)
		load()
	})
}

func (n *NotificationSystem) Start() {
	go func() {
		for notification := range n.notificationQueue {
			if hasGlobalNotifierSet(notification.Service.Notifiers) {
				for _, notifier := range n.getEnabledNotifiers() {
					if err := sendNotification(notifier, notification); err != nil {
						log.Errorf("Unable to send notification: '%s' - '%s'", notifier.GetId(), err)
					}
				}
			} else {
				for _, notifierId := range notification.Service.Notifiers {
					notifier, err := n.getNotifierById(notifierId)
					if err != nil {
						log.Error(err)
						continue
					}

					if err := sendNotification(notifier, notification); err != nil {
						log.Errorf("Unable to send notification: '%s' - '%s'", notifier.GetId(), err)
					}
				}
			}
		}
	}()
}

func hasGlobalNotifierSet(notifierIds []string) bool {
	for _, id := range notifierIds {
		if id == globalNotifierId {
			return true
		}
	}
	return false
}

func sendNotification(notifier model.Notify, notification Notification) error {
	log.Infof("Sending notification using '%s' notifier", notifier.GetId())

	var tmpl *template.Template
	var err error

	if notification.IsUpNotification {
		tmpl, err = template.New(notifier.GetId()).Parse(notifier.GetServiceUpNotificationTemplate())
	} else {
		tmpl, err = template.New(notifier.GetId()).Parse(notifier.GetServiceDownNotificationTemplate())
	}

	if err != nil {
		log.Errorf("Unable to parse template for notifier '%s' - '%s'", notifier.GetId(), err)
		return err
	}

	data := model.TemplateData{
		Name:   notification.Service.Name,
		Date:   time.Now().Format(time.RFC3339),
		Reason: notification.Failure.Reason,
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		log.Errorf("Unable to execute template for notifier '%s' - '%s'", notifier.GetId(), err)
		return err
	}

	if err := notifier.SendNotification(notification.Service, buffer.String()); err != nil {
		log.Errorf("Unable to send notification with notifier '%s' - '%s'", notifier.GetId(), err)
	}
	return nil
}

func (n *NotificationSystem) getEnabledNotifiers() []model.Notify {
	result := make([]model.Notify, 0, len(n.notifiers))
	for _, notifier := range n.notifiers {
		if notifier.IsEnabled() {
			result = append(result, notifier)
		}
	}
	return result
}

func (n *NotificationSystem) getNotifierById(id string) (model.Notify, error) {
	for _, notifier := range n.notifiers {
		if notifier.GetId() == id {
			return notifier, nil
		}
	}
	return nil, fmt.Errorf("notifier with id '%s' not found", id)
}

func (n *NotificationSystem) AddNotification(notification Notification) {
	go func() {
		logFields := log.WithFields(log.Fields{"serviceId": notification.Service.Id, "isUpNotification": notification.IsUpNotification})
		logFields.Infof("Adding notification to notification queu")
		// add notification non blocking
		n.notificationQueue <- notification
	}()
}

func (n *NotificationSystem) GetNotifiers() []model.Notify {
	return n.notifiers
}

type Notification struct {
	Id               string
	Service          model.Service
	IsUpNotification bool
	Failure          model.Failure
}

func NewNotification(service model.Service, isUpNotification bool, failure model.Failure) Notification {
	return Notification{
		Id:               uuid.New().String(),
		Service:          service,
		IsUpNotification: isUpNotification,
		Failure:          failure,
	}
}
