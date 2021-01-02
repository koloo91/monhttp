package notifier

import (
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type NotificationSystem struct {
	notifiers         []model.Notify
	notificationQueue chan Notification
}

func NewNotificationSystem() *NotificationSystem {
	return &NotificationSystem{
		notifiers:         make([]model.Notify, 0),
		notificationQueue: make(chan Notification, 100),
	}
}

func (n *NotificationSystem) SetupDefaultNotifier() {
	n.notifiers = make([]model.Notify, 0)

	load := func() {
		n.notifiers = append(n.notifiers, NewEMailNotifier())
		n.notifiers = append(n.notifiers, NewTelegramNotifier())
	}
	load()

	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("Config files changed: ", in.Name)
		// reset notifiers
		n.notifiers = make([]model.Notify, 0)
		load()
	})
}

func (n *NotificationSystem) Start() {
	go func() {
		for notification := range n.notificationQueue {
			for _, notifier := range n.getEnabledNotifiers() {
				log.Printf("Sending notification using '%s' notifier", notifier.GetId())
				if err := notifier.SendFailure(notification.Service, notification.Failure); err != nil {
					log.Error(err)
				}
			}
		}
	}()
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

func (n *NotificationSystem) AddNotification(notification Notification) {
	go func() {
		// add notification non blocking
		n.notificationQueue <- notification
	}()
}

func (n *NotificationSystem) GetNotifiers() []model.Notify {
	return n.notifiers
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
