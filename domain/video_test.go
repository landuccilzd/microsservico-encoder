package domain_test

import (
	"encoder/domain"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()
	require.Error(t, err)
}

func TestVideoWhenIDIsNotUUID(t *testing.T) {
	video := domain.NewVideo()
	video.ID = "not an UUID"
	video.ResourceID = "123"
	video.FilePath = "/home"
	video.CreatedAt = time.Now()
	err := video.Validate()
	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.ResourceID = "123"
	video.FilePath = "/home"
	video.CreatedAt = time.Now()
	err := video.Validate()
	require.Nil(t, err)
}
