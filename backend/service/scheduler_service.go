package service

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/notifier"
	"github.com/koloo91/monhttp/repository"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func StartScheduleJob() {
	log.Info("Starting job scheduler")

	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		go startCheckProcess()
	}
}

func startCheckProcess() {
	log.Info("Looking for next services to process")
	serviceIds, err := getNextServiceIds()
	if err != nil {
		log.Errorf("Unable to get next service ids to process '%s'", err)
		return
	}

	for _, service := range serviceIds {
		go processService(service)
	}

}

func getNextServiceIds() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return repository.GetNextScheduledServiceIds(ctx)
}

func processService(serviceId string) {
	processLogger := log.New()
	processLogger.SetFormatter(&log.TextFormatter{})
	logger := processLogger.WithField("serviceId", serviceId)

	logger.Infof("Processing service with id: '%s'", serviceId)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tx, err := repository.BeginnTransaction()
	if err != nil {
		logger.Errorf("Unable to start transaction: '%s'", err)
		return
	}

	service, err := repository.SelectServiceByIdLocked(ctx, tx, serviceId)
	if err != nil {
		logger.Errorf("Unable to lock service: '%s'", err)
		if err := tx.Rollback(); err != nil {
			logger.Errorf("Error rolling back transaction: '%s'", err)
		}
		return
	}

	logger.Infof("Locked service '%s'", service.Name)

	nextCheckTime := time.Now().Add(time.Duration(service.IntervalInSeconds) * time.Second)
	logger.Infof("Set next check time of service '%s' to %s", service.Name, nextCheckTime.String())

	if err := repository.UpdateServiceByIdNextCheckTime(ctx, tx, serviceId, nextCheckTime); err != nil {
		logger.Errorf("Unable to update service '%s' next check time: '%s'", service.Name, err)
		if err := tx.Rollback(); err != nil {
			logger.Errorf("Error rolling back transaction: '%s'", err)
		}
		return
	}

	var check *model.Check
	var failure *model.Failure
	var checkErr error

	switch service.Type {
	case model.ServiceTypeHttp:
		logger.Infof("Processing service '%s' as type HTTP", service.Name)
		check, failure, checkErr = handleHttpServiceType(service)
	case model.ServiceTypeIcmpPing:
		logger.Infof("Processing service '%s' as type ICMP Ping", service.Name)
		check, failure, checkErr = handleIcmpPingServiceType(service)
	default:
		logger.Warnf("Unknown service type '%s'", service.Type)
	}

	if checkErr != nil {
		logger.Errorf("Error handling service type: '%s' - '%s'", service.Name, err)
		if err := tx.Rollback(); err != nil {
			logger.Errorf("Error rolling back transaction: '%s'", err)
		}
		return
	}

	if check != nil {
		if err := repository.InsertCheck(ctx, tx, *check); err != nil {
			logger.Errorf("Unable to insert check for service '%s' - '%s'", service.Name, err)
			if err := tx.Rollback(); err != nil {
				logger.Errorf("Error rolling back transaction: '%s'", err)
			}
			return
		}
	}

	if failure != nil {
		if service.EnableNotifications {
			logger.Infof("Notifications for service '%s' enabled", service.Name)
			sendNotification, err := shouldSendNotification(ctx, tx, service)
			if err != nil {
				logger.Errorf("Unable to determine if we should send a notfication for service '%s' - '%s'", service.Name, err)
				if err := tx.Rollback(); err != nil {
					logger.Errorf("Error rolling back transaction: '%s'", err)
				}
				return
			}

			if sendNotification {
				logger.Infof("Sending notification for service '%s'", service.Name)
				notificationSystem.AddNotification(notifier.NewNotification(service, *failure))
			}
		}

		if err := repository.InsertFailure(ctx, tx, *failure); err != nil {
			logger.Errorf("Unable to insert failure for service '%s' - '%s'", service.Name, err)
			if err := tx.Rollback(); err != nil {
				logger.Errorf("Error rolling back transaction: '%s'", err)
			}
			return
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Errorf("Error commiting transaction: '%s'", err)
	}
}

func shouldSendNotification(ctx context.Context, tx *sql.Tx, service model.Service) (bool, error) {
	checks, err := repository.GetLastNChecks(ctx, tx, service.Id, service.NotifyAfterNumberOfFailures)
	if err != nil {
		return false, err
	}

	counter := 0
	for _, check := range checks {
		if check.IsFailure {
			counter++
		}
	}

	sendNotification := false
	if service.ContinuouslySendNotifications {
		sendNotification = counter+1 >= service.NotifyAfterNumberOfFailures
	} else {
		sendNotification = counter+1 == service.NotifyAfterNumberOfFailures
	}
	return sendNotification, nil
}

func handleHttpServiceType(service model.Service) (*model.Check, *model.Failure, error) {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !service.VerifySsl,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !service.FollowRedirects {
				return fmt.Errorf("i am not allowed to follow redirects")
			}
			return nil
		},
		Timeout: time.Duration(service.RequestTimeoutInSeconds) * time.Second,
	}

	request, err := http.NewRequest(service.HttpMethod, service.Endpoint, strings.NewReader(service.HttpBody))
	if err != nil {
		return nil, nil, err
	}

	headers := strings.Split(service.HttpHeaders, ";")
	for _, header := range headers {
		headerValues := strings.Split(header, ":")
		if len(headerValues) != 2 {
			continue
		}

		headerKey := headerValues[0]
		headerValue := headerValues[1]

		request.Header.Add(headerKey, headerValue)
	}

	start := time.Now()
	response, err := client.Do(request)
	if err != nil {
		return model.NewCheck(service.Id, 0, true), model.NewFailure(service.Id, err.Error()), nil
	}
	defer response.Body.Close()

	latency := time.Since(start)

	if response.StatusCode != service.ExpectedHttpStatusCode {
		reason := fmt.Sprintf("Expected status code '%d' but got '%d'", service.ExpectedHttpStatusCode, response.StatusCode)
		failure := model.NewFailure(service.Id, reason)

		check := model.NewCheck(service.Id, 0, true)
		return check, failure, nil
	}

	if len(service.ExpectedHttpResponseBody) > 0 {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			reason := fmt.Sprintf("Unable to read response body: %s", err.Error())
			failure := model.NewFailure(service.Id, reason)
			return model.NewCheck(service.Id, 0, true), failure, nil
		} else {
			matched, err := regexp.Match(service.ExpectedHttpResponseBody, bodyBytes)
			if err != nil {
				reason := fmt.Sprintf("Unable to read response body: %s", err.Error())
				failure := model.NewFailure(service.Id, reason)
				return model.NewCheck(service.Id, 0, true), failure, nil
			}

			if !matched {
				reason := fmt.Sprintf("Body did not match '%s'", service.ExpectedHttpResponseBody)
				failure := model.NewFailure(service.Id, reason)
				return model.NewCheck(service.Id, 0, true), failure, nil
			}
		}

	}

	return model.NewCheck(service.Id, latency.Milliseconds(), false), nil, nil
}

func handleIcmpPingServiceType(service model.Service) (*model.Check, *model.Failure, error) {
	ping, err := exec.LookPath("ping")
	if err != nil {
		return nil, nil, err
	}

	command := exec.Command(ping, service.Endpoint, "-c", "1", "-W", strconv.Itoa(service.RequestTimeoutInSeconds*1000))
	outputBytes, err := command.CombinedOutput()
	if err != nil {
		log.Error(err)
	}

	outputString := string(outputBytes)
	if strings.Contains(outputString, "Unknown host") {
		failure := model.NewFailure(service.Id, "unknown host")
		return nil, failure, nil
	}

	if strings.Contains(outputString, "100.0% packet loss") {
		failure := model.NewFailure(service.Id, "destination host unreachable")
		return nil, failure, nil
	}

	r := regexp.MustCompile(`time=(.*) ms`)
	submatches := r.FindStringSubmatch(outputString)
	if len(submatches) < 2 {
		failure := model.NewFailure(service.Id, "could not parse ping duration")
		check := model.NewCheck(service.Id, 0, true)
		return check, failure, nil
	}

	duration, _ := strconv.ParseFloat(submatches[1], 64)
	check := model.NewCheck(service.Id, int64(duration), false)

	return check, nil, nil
}
