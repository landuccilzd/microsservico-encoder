package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
)

type JobRepository interface {
	Insert(job *domain.Job) (*domain.Job, error)
	Find(id string) (*domain.Job, error)
	Update(job *domain.Job) (*domain.Job, error)
}

type JobRepositoryDB struct {
	Db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepositoryDB {
	return &JobRepositoryDB{Db: db}
}

func (repository *JobRepositoryDB) Insert(job *domain.Job) (*domain.Job, error) {
	if err := repository.Db.Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (repository *JobRepositoryDB) Find(id string) (*domain.Job, error) {
	var job domain.Job
	repository.Db.Preload("Video").First(&job, "id = ?", id)

	if job.ID == "" {
		return nil, fmt.Errorf("não existe um job com o id %v", id)
	}

	return &job, nil
}

func (repository *JobRepositoryDB) Update(job *domain.Job) (*domain.Job, error) {
	if err := repository.Db.Save(&job).Error; err != nil {
		return nil, err
	}
	return job, nil
}
