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
		notificationQueue: make(chan Notification, 1024),
	}
}

func (n *NotificationSystem) SetupDefaultNotifier() {
	log.Info("Setting up default notifiers")
	n.notifiers = make([]model.Notify, 0)

	load := func() {
		log.Info("Adding notifiers")
		n.notifiers = append(n.notifiers, NewEMailNotifier())
		n.notifiers = append(n.notifiers, NewTelegramNotifier())
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
			for _, notifier := range n.getEnabledNotifiers() {
				log.Infof("Sending notification using '%s' notifier", notifier.GetId())
				if notification.IsUpNotification {
					if err := notifier.SendServiceIsUpNotification(notification.Service); err != nil {
						log.Errorf("Unable to send notification with notifier '%s' - '%s'", notifier.GetId(), err)
					}
				} else {
					if err := notifier.SendServiceIsDownNotification(notification.Service, notification.Failure); err != nil {
						log.Errorf("Unable to send notification with notifier '%s' - '%s'", notifier.GetId(), err)
					}
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
