package services

import (
	"encoding/json"
	"enconder/application/repositories"
	"enconder/domain"
	"enconder/framework/queue"
	"log"
	"os"
	"strconv"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type JobManager struct {
	DB										*gorm.DB
	Domain								domain.Job
	MessageChannel chan 	amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ							*queue.RabbitMQ
}

type JobNotificationError struct {
	Message string	`json:"messages"`
	Error		string	`json:"error"`
}

func NewJobManager(db *gorm.DB, rabbitMQ *queue.RabbitMQ, jobReturnChannel chan JobWorkerResult, messageChannel chan amqp.Delivery) *JobManager {

	return &JobManager{
		DB: db,
		Domain: domain.Job{},
		MessageChannel: messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ: rabbitMQ,
	}

}

func (j *JobManager) Start(ch *amqp.Channel) {

	videoService := NewVideoService()
	videoService.VideoRepository = repositories.VideoRepositoryDB{DB: j.DB}

	jobService := JobService{
		JobRepository: repositories.JobRepositoryDb{Db: j.DB},
		VideoService: videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))
	if err != nil {
		log.Fatalf("error loading var: CONCURRENCY_WORKERS!")
	}

	
	for qtdProcesses := 0; qtdProcesses < concurrency; qtdProcesses++ {
		go JobWorker(j.MessageChannel, j.JobReturnChannel, jobService, j.Domain, qtdProcesses)
	}

	for jobResult := range j.JobReturnChannel {
		if jobResult.Error != nil {
			err = j.checkParseErrors(jobResult)
		} else {
			err = j.notifySucess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}

}

func (j *JobManager) notifySucess(jobResult JobWorkerResult, ch *amqp.Channel) error {
	jobJson, err := json.Marshal(jobResult.Job)

	if err != nil {
		return err
	}

	err = j.notify(jobJson)

	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error parsing job: #{jobResult.Job.ID}")
	} else {
		log.Printf("MessageID #{jobResult.Message.DeliveryTag}. Error with messaage: #{jobResult.Job.Error}")
	}

	// errorMsg := JobNotificationError{
	// 	Message: string(jobResult.Message.Body),
	// 	Error: jobResult.Error.Error(),
	// }
	//
	// jobJson, err := json.Marshal(errorMsg)

	// Is missing return notification
	
	return nil
}

func (j *JobManager) notify(jobJson []byte) error {
	err := j.RabbitMQ.Notify(
		string(jobJson), 
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBIT_NOTIFICATION_ROUTING_KEY"),
	)

	if err != nil {
		return err
	}

	return nil
}
