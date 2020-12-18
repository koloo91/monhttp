package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/repository"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func StartScheduleJob() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		go startCheckProcess()
	}
}

func startCheckProcess() {
	serviceIds, err := getNextServiceIds()
	if err != nil {
		log.Error(err)
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
	log.Infof("Processing service with id: '%s'", serviceId)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tx, err := repository.BeginnTransaction()
	if err != nil {
		log.Error(err)
		return
	}

	service, err := repository.SelectServiceByIdLocked(ctx, tx, serviceId)
	if err != nil {
		log.Error(err)
		if err := tx.Rollback(); err != nil {
			log.Error(err)
		}
		return
	}

	nextCheckTime := time.Now().Add(time.Duration(service.IntervalInSeconds) * time.Second)
	if err := repository.UpdateServiceByIdNextCheckTime(ctx, tx, serviceId, nextCheckTime); err != nil {
		log.Error(err)
		if err := tx.Rollback(); err != nil {
			log.Error(err)
		}
		return
	}

	var check *model.Check
	var failure *model.Failure
	var checkErr error

	switch service.Type {
	case model.ServiceTypeHttp:
		log.Info("Found HTTP service")
		check, failure, checkErr = handleHttpServiceType(service)
	case model.ServiceTypeIcmpPing:
		log.Info("Found ICMP Ping service")
		check, failure, checkErr = handleIcmpPingServiceType(service)
	default:
		log.Warnf("Unknown service type '%s'", service.Type)
	}

	if checkErr != nil {
		log.Error(err)
		if err := tx.Rollback(); err != nil {
			log.Error(err)
		}
		return
	}

	if check != nil {
		if err := repository.InsertCheck(ctx, tx, *check); err != nil {
			log.Error("insert check", err)
			if err := tx.Rollback(); err != nil {
				log.Error("insert check", err)
			}
			return
		}
	}

	if failure != nil {
		if err := repository.InsertFailure(ctx, tx, *failure); err != nil {
			log.Error("insert failure", err)
			if err := tx.Rollback(); err != nil {
				log.Error("insert failure", err)
			}
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error(err)
	}
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

	// TODO: add headers, add request body
	request, err := http.NewRequest(service.HttpMethod, service.Endpoint, strings.NewReader(service.HttpBody))
	if err != nil {
		return nil, nil, err
	}

	start := time.Now()
	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != service.ExpectedHttpStatusCode {
		reason := fmt.Sprintf("Expected status code '%d' but got '%d'", service.ExpectedHttpStatusCode, response.StatusCode)
		failure := model.NewFailure(service.Id, reason)

		check := model.NewCheck(service.Id, 0, true)
		return check, failure, nil
	}

	// TODO: check body

	latency := time.Since(start)
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
		return nil, nil, err
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
	strs := r.FindStringSubmatch(outputString)
	if len(strs) < 2 {
		failure := model.NewFailure(service.Id, "could not parse ping duration")
		check := model.NewCheck(service.Id, 0, true)
		return check, failure, nil
	}

	duration, _ := strconv.ParseFloat(strs[1], 64)
	check := model.NewCheck(service.Id, int64(duration), false)

	return check, nil, nil
}
