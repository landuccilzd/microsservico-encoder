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

func TestJobRepositoryDB_Insert(t *testing.T) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "/home"
	video.CreatedAt = time.Now()

	vRepository := repositories.VideoRepositoryDB{Db: db}
	vRepository.Insert(video)

	job, _ := domain.NewJob("/home/video-output", "CREATED", video)
	repository := repositories.JobRepositoryDB{Db: db}
	jobCreated, err := repository.Insert(job)
	require.Nil(t, err)
	require.NotNil(t, jobCreated)
	require.NotEmpty(t, jobCreated.ID)
	require.Equal(t, jobCreated.ID, job.ID)
	require.Equal(t, jobCreated.OutputBucketPath, job.OutputBucketPath)
	require.Equal(t, jobCreated.Status, job.Status)
}

func TestJobRepositoryDB_Find(t *testing.T) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "/home"
	video.CreatedAt = time.Now()

	vRepository := repositories.VideoRepositoryDB{Db: db}
	vRepository.Insert(video)

	job, _ := domain.NewJob("/home/video-output", "CREATED", video)
	repository := repositories.JobRepositoryDB{Db: db}
	repository.Insert(job)

	jobFound, err := repository.Find(job.ID)
	require.Nil(t, err)
	require.NotEmpty(t, jobFound.ID)
	require.Equal(t, jobFound.ID, job.ID)
	require.Equal(t, jobFound.OutputBucketPath, job.OutputBucketPath)
	require.Equal(t, jobFound.Status, job.Status)
}

func TestJobRepositoryDB_Update(t *testing.T) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "/home"
	video.CreatedAt = time.Now()

	vRepository := repositories.VideoRepositoryDB{Db: db}
	vRepository.Insert(video)

	job, _ := domain.NewJob("/home/video-output", "CREATED", video)
	repository := repositories.JobRepositoryDB{Db: db}
	repository.Insert(job)

	job.OutputBucketPath = "/home/video-output-2"
	job.Status = "UPDATED"
	job.UpdatedAt = time.Now()

	jobUpdated, err := repository.Update(job)
	require.Nil(t, err)
	require.NotNil(t, jobUpdated)
	require.NotEmpty(t, jobUpdated.ID)
	require.Equal(t, jobUpdated.ID, job.ID)
	require.Equal(t, jobUpdated.OutputBucketPath, job.OutputBucketPath)
	require.Equal(t, jobUpdated.Status, job.Status)
}
