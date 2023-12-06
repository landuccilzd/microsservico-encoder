package services_test

import (
	"encoder/application/services"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("erro ao carregar o arquivo .env")
	}
}

func TestUploadService_Upload(t *testing.T) {
	video, repositorio := prepare()
	service := services.NewVideoService()
	service.Video = video
	service.VideoRepository = &repositorio

	err := service.Fragment()
	require.Nil(t, err)

	err = service.Encode()
	require.Nil(t, err)

	err = service.Finish()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "codeeducationstest"
	videoUpload.VideoPath = os.Getenv("localStorageFragPath") + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(10, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "Upload Completed")
}
