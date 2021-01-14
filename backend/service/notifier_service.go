package service

import (
	"bytes"
	"fmt"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/notifier"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"text/template"
	"time"
)

func GetNotifiers() []model.Notify {
	return notificationSystem.GetNotifiers()
}

func UpdateNotifier(id string, body map[string]interface{}) error {
	log.Infof("Updating notififier with id '%s'", id)
	for k, v := range body {
		viper.Set(fmt.Sprintf("NOTIFIER_%s_%s", strings.ToUpper(id), strings.ToUpper(k)), v)
	}

	notificationSystem.SetupDefaultNotifier()

	return viper.WriteConfig()
}

func TestNotifierUpTemplate(id string, body map[string]interface{}) error {
	log.Infof("Test notififiers up template with id '%s'", id)

	testNotify, err := setupTestNotifier(id, body)
	if err != nil {
		return err
	}

	tmpl, err := template.New(id).Parse(testNotify.GetServiceUpNotificationTemplate())
	if err != nil {
		log.Errorf("Unable to parse template for test notifier '%s' - '%s'", id, err)
		return err
	}

	data := model.TemplateData{
		Name:   "Test Service Name",
		Date:   time.Now().Format(time.RFC3339),
		Reason: "",
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		log.Errorf("Unable to execute template for test notifier '%s' - '%s'", id, err)
		return err
	}

	return testNotify.SendNotification(model.Service{}, buffer.String())
}

func TestNotifierDownTemplate(id string, body map[string]interface{}) error {
	log.Infof("Test notififiers down template with id '%s'", id)

	testNotify, err := setupTestNotifier(id, body)
	if err != nil {
		return err
	}

	tmpl, err := template.New(id).Parse(testNotify.GetServiceDownNotificationTemplate())
	if err != nil {
		log.Errorf("Unable to parse template for test notifier '%s' - '%s'", id, err)
		return err
	}

	data := model.TemplateData{
		Name:   "Test Service Name",
		Date:   time.Now().Format(time.RFC3339),
		Reason: "This is just a test",
	}

	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, data); err != nil {
		log.Errorf("Unable to execute template for test notifier '%s' - '%s'", id, err)
		return err
	}

	return testNotify.SendNotification(model.Service{}, buffer.String())
}

func setupTestNotifier(id string, body map[string]interface{}) (model.Notify, error) {
	testStore := viper.New()
	for k, v := range body {
		testStore.Set(fmt.Sprintf("NOTIFIER_%s_%s", strings.ToUpper(id), strings.ToUpper(k)), v)
	}

	var testNotify model.Notify
	switch id {
	case "email":
		testNotify = notifier.NewEMailNotifier(testStore)
	case "telegram":
		testNotify = notifier.NewTelegramNotifier(testStore)
	default:
		return nil, fmt.Errorf("notifier with id '%s' is unknown", id)
	}

	return testNotify, nil
}
