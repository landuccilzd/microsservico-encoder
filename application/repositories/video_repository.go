package repositories

import (
	"encoder/domain"
	"fmt"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDB struct {
	Db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepositoryDB {
	return &VideoRepositoryDB{Db: db}
}

func (repository VideoRepositoryDB) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}

	if err := repository.Db.Create(video).Error; err != nil {
		return nil, err
	}

	return video, nil
}

func (repository VideoRepositoryDB) Find(id string) (*domain.Video, error) {
	var video domain.Video
	repository.Db.Preload("Jobs").First(&video, "id = ?", id)

	if video.ID == "" {
		return nil, fmt.Errorf("não existe um video com o id %v", id)
	}

	return &video, nil
}
