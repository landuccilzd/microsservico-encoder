package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/queue"
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	DB               *gorm.DB
	Job              domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(db *gorm.DB, rabbitMQ *queue.RabbitMQ, jobReturnChannel chan JobWorkerResult, messageChannel chan amqp.Delivery) *JobManager {
	return &JobManager{
		DB:               db,
		Job:              domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ:         rabbitMQ,
	}
}

func (jm *JobManager) Start(channel *amqp.Channel) {
	videoService := NewVideoService()
	videoService.VideoRepository = repositories.VideoRepositoryDB{Db: jm.DB}

	jobService := JobService{
		JobRepository: repositories.JobRepositoryDB{Db: jm.DB},
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))
	if err != nil {
		log.Fatalf("error loading environment var CONCURRENCY_WORKERS")
	}

	for qtProcs := 0; qtProcs < concurrency; qtProcs++ {
		go JobWorker(jm.MessageChannel, jm.JobReturnChannel, jobService, jm.Job, qtProcs)
	}

	for jobResult := range jm.JobReturnChannel {
		if jobResult.Error != nil {
			err = jm.checkParseErrors(jobResult)
		} else {
			err = jm.notifySuccess(jobResult, channel)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}

	}
}

func (jm *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("123 MessageID: %v. Error with job: %v with video %v. Error: %v",
			jobResult.Message.DeliveryTag, jobResult.Job.ID, jobResult.Job.Video.ID, jobResult.Error)
	} else {
		log.Printf("MessageID: %v. Error parsing message: %v", jobResult.Message.DeliveryTag, jobResult.Error.Error())
	}

	errorMessage := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}

	jobJson, _ := json.Marshal(errorMessage)
	if err := jm.notify(jobJson); err != nil {
		return err
	}

	if err := jobResult.Message.Reject(false); err != nil {
		return err
	}

	return nil
}

func (jm *JobManager) notifySuccess(jobResult JobWorkerResult, channel *amqp.Channel) error {
	jobJson, err := json.Marshal(jobResult.Job)
	if err != nil {
		return err
	}

	err = jm.notify(jobJson)
	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)
	if err != nil {
		return err
	}

	return nil
}

func (jm *JobManager) notify(jobJson []byte) error {
	err := jm.RabbitMQ.Notify(
		string(jobJson),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)

	if err != nil {
		return err
	}

	return nil
}
