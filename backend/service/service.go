package service

import "github.com/koloo91/monhttp/notifier"

var (
	notificationSystem *notifier.NotificationSystem
)

func SetNotificationSystem(n *notifier.NotificationSystem) {
	notificationSystem = n
}
