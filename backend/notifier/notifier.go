package notifier

import (
	"github.com/google/uuid"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
)

type NotificationSystem struct {
	notifiers         []Notify
	notificationQueue chan Notification
}

func NewNotificationSystem() *NotificationSystem {
	return &NotificationSystem{
		notifiers:         make([]Notify, 0),
		notificationQueue: make(chan Notification, 100),
	}
}

func (n *NotificationSystem) SetupDefaultNotifier() {
	n.notifiers = append(n.notifiers, NewEMailNotifier())
}

func (n *NotificationSystem) Start() {
	go func() {
		for notification := range n.notificationQueue {
			for _, notifier := range n.getEnabledNotifiers() {
				log.Printf("Sending notification using '%s' notifier", notifier.Id())
				if err := notifier.SendFailure(notification.Service, notification.Failure); err != nil {
					log.Error(err)
				}
			}
		}
	}()
}

func (n *NotificationSystem) getEnabledNotifiers() []Notify {
	result := make([]Notify, 0, len(n.notifiers))
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
	Id() string
	SendSuccess(model.Service) error
	SendFailure(model.Service, model.Failure) error
	IsEnabled() bool
}

type Notifier struct {
	Id      string
	Name    string
	Enabled bool
}
