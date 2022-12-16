package queue_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"queue-cleaner/config"
	"queue-cleaner/queue"
)

func ListQueue() ([]queue.Queue, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/queues/", config.GetManagementUrl()), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(config.GetManagementBasicAuth())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return nil, err
	}

	res := make([]queue.Queue, 0)
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("get queue list error")
	} else {
		err = json.Unmarshal(buf.Bytes(), &res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func CloseQueue(queueName string) error {
	dto := CloseQueueBody{
		Properties:      Properties{Headers: Headers{ControlMsg: "close"}},
		RoutingKey:      queueName,
		DeliveryMode:    "1",
		Payload:         "",
		PayloadEncoding: "string",
	}

	payload, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/exchanges/%%2f/amq.default/publish", config.GetManagementUrl()), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(config.GetManagementBasicAuth())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("close queue error")
	} else {
		fmt.Println(buf.Bytes())
	}

	return nil
}

func DeleteQueue(queueName string) error {
	queueNameEncoded := url.QueryEscape(queueName)

	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/queues/%%2F/%s", config.GetManagementUrl(), queueNameEncoded), nil)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(config.GetManagementBasicAuth())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete queue error: response status code: %d", response.StatusCode)
	}

	return nil
}
