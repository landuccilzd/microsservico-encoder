package services_test

import (
	"encoder/application/repositories"
	"encoder/application/services"
	"encoder/domain"
	"encoder/framework/database"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("erro ao carregar o arquivo .env")
	}
}

func prepare() (*domain.Video, repositories.VideoRepository) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = os.Getenv("localStoragePath") + "/video-danca.mp4"
	video.CreatedAt = time.Now()

	repositorio := repositories.VideoRepositoryDB{Db: db}
	return video, repositorio
}

func TesteVideoService_Download(t *testing.T) {
	video, repositorio := prepare()
	service := services.NewVideoService()
	service.Video = video
	service.VideoRepository = &repositorio

	err := service.Download("nome_do_bucket")
	require.Nil(t, err)
}

func TestVideoService_Fragment(t *testing.T) {
	video, repositorio := prepare()
	service := services.NewVideoService()
	service.Video = video
	service.VideoRepository = &repositorio

	err := service.Fragment()
	require.Nil(t, err)
}

func TestVideoService_Encode(t *testing.T) {
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
}
