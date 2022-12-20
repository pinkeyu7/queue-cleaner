package queue_management

import (
	"fmt"
	"queue-cleaner/config"
	"queue-cleaner/queue"
	"queue-cleaner/queue_api"
	"regexp"
	"strings"
)

func InactiveQueueList() ([]queue.QueueWithType, error) {
	needDeleteQueue := make([]queue.QueueWithType, 0)

	queueList, err := queue_api.ListQueue()
	if err != nil {
		return nil, err
	}

	queueWithTypes := queueWithTypesList(queueList)

	remnantQueue := remnant(queueWithTypes)
	//needDeleteQueue = append(needDeleteQueue, remnantQueue...)
	printState("Remnant Queue", remnantQueue)

	emptyConsumerQueue := emptyConsumer(queueWithTypes)
	needDeleteQueue = append(needDeleteQueue, emptyConsumerQueue...)
	printState("Empty Consumer Queue", emptyConsumerQueue)

	deleteQueueLength, deletedAmount, err := deleteQueue(needDeleteQueue)
	if err != nil {
		return nil, err
	}

	fmt.Println("deleteQueueLength:", deleteQueueLength, ", deletedAmount:", deletedAmount)

	return needDeleteQueue, nil
}

func deleteQueue(queueWithTypes []queue.QueueWithType) (int, int, error) {
	sessionMap := make(map[string]bool)
	list := make([]queue.QueueWithType, 0)
	count := 0

	for _, q := range queueWithTypes {
		if _, value := sessionMap[q.Name]; !value {
			sessionMap[q.Name] = true
			list = append(list, q)
		}
	}

	listLen := len(list)

	if config.IsDeleteMode() {
		for _, q := range list {
			err := queue_api.DeleteQueue(q.Name)
			if err != nil {
				return listLen, count, err
			}
			count++
		}
	}

	return listLen, count, nil
}

func printState(name string, queueWithTypes []queue.QueueWithType) {
	fmt.Println(name)
	for _, q := range queueWithTypes {
		fmt.Println(q)
	}
	fmt.Println("Amount:", len(queueWithTypes))
	fmt.Println("-------------------------------------------")
}

func remnant(queueWithTypes []queue.QueueWithType) []queue.QueueWithType {
	sessionMap := make(map[string]int)
	ssMap := make(map[int]int)
	for _, q := range queueWithTypes {
		sessionMap[q.SessionId] += 1
	}

	for _, times := range sessionMap {
		ssMap[times] += 1
	}

	for times, time := range ssMap {
		fmt.Println("Session Stage Amount:", times, "=>", time)
	}

	results := make([]queue.QueueWithType, 0)
	for _, q := range queueWithTypes {
		if sessionMap[q.SessionId] == 1 && q.Messages == 0 {
			results = append(results, q)
		}
	}

	return results
}

func emptyConsumer(queueWithTypes []queue.QueueWithType) []queue.QueueWithType {
	results := make([]queue.QueueWithType, 0)
	for _, q := range queueWithTypes {
		if q.Consumers == 0 {
			results = append(results, q)
		}
	}
	return results
}

func queueWithTypesList(queueList []queue.Queue) []queue.QueueWithType {
	r, _ := regexp.Compile(`.*/.*/.*.(IN|OUT)$`)

	queueWithTypes := make([]queue.QueueWithType, 0)
	for _, q := range queueList {
		match := r.FindStringSubmatch(q.Name)
		if len(match) < 2 {
			continue
		} else {
			results := strings.Split(strings.Split(q.Name, "/")[2], ".")

			queueWithTypes = append(queueWithTypes, queue.QueueWithType{
				Queue:     q,
				QueueType: match[1],
				SessionId: results[0],
				Stage:     results[1],
			})
		}
	}

	return queueWithTypes
}
