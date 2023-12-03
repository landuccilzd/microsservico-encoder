package repositories_test

import (
	"encoder/application/repositories"
	"encoder/domain"
	"encoder/framework/database"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoRepository_Insert(t *testing.T) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "/home"
	video.CreatedAt = time.Now()

	repository := repositories.VideoRepositoryDB{Db: db}
	videoCreated, err := repository.Insert(video)
	require.Nil(t, err)
	require.NotNil(t, videoCreated)
	require.NotEmpty(t, videoCreated.ID)
	require.Equal(t, videoCreated.ID, video.ID)
	require.Equal(t, videoCreated.ResourceID, video.ResourceID)
	require.Equal(t, videoCreated.FilePath, video.FilePath)
}

func TestVideoRepository_Find(t *testing.T) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "/home"
	video.CreatedAt = time.Now()

	repository := repositories.VideoRepositoryDB{Db: db}
	repository.Insert(video)

	videoFound, err := repository.Find(video.ID)
	require.Nil(t, err)
	require.NotEmpty(t, videoFound.ID)
	require.Equal(t, videoFound.ID, video.ID)
	require.Equal(t, videoFound.ResourceID, video.ResourceID)
	require.Equal(t, videoFound.FilePath, video.FilePath)
}
