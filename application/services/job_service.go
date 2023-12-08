package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"log"
	"os"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (service *JobService) Start() error {
	if err := service.ChangeJobStatus("DOWNLOADING"); err != nil {
		return service.failJob(err)
	}

	if err := service.VideoService.Download(os.Getenv("inputBucketName")); err != nil {
		return service.failJob(err)
	}

	if err := service.ChangeJobStatus("FRAGMENTING"); err != nil {
		return service.failJob(err)
	}

	if err := service.VideoService.Fragment(); err != nil {
		return service.failJob(err)
	}

	if err := service.ChangeJobStatus("ENCODING"); err != nil {
		return service.failJob(err)
	}

	if err := service.VideoService.Encode(); err != nil {
		return service.failJob(err)
	}

	if err := service.performUpload(); err != nil {
		return service.failJob(err)
	}

	if err := service.ChangeJobStatus("FINISHING"); err != nil {
		return service.failJob(err)
	}

	if err := service.VideoService.Finish(); err != nil {
		return service.failJob(err)
	}

	if err := service.ChangeJobStatus("COMPLETED"); err != nil {
		return service.failJob(err)
	}

	return nil
}

func (service *JobService) ChangeJobStatus(status string) error {
	var err error

	service.Job.Status = status
	service.Job, err = service.JobRepository.Update(service.Job)
	if err != nil {
		return service.failJob(err)
	}

	return nil
}

func (service *JobService) failJob(error error) error {
	service.Job.Status = "FAILED"
	service.Job.Error = error.Error()

	_, err := service.JobRepository.Update(service.Job)
	if err != nil {
		log.Fatalf(err.Error())
		return err
	}

	return error
}

func (service *JobService) performUpload() error {
	if err := service.ChangeJobStatus("UPLOADING"); err != nil {
		return service.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = os.Getenv("locarStorageFragPath") + "/" + service.VideoService.Video.ID
	// concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	// doneUpload := make(chan string)

	// go videoUpload.ProcessUpload(concurrency, doneUpload)
	// var uploadResult = <-doneUpload

	// if uploadResult != "Upload Completed" {
	// return service.failJob(errors.New(uploadResult))
	// }

	return nil
}
