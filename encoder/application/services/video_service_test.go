package services_test

import (
	"enconder/application/repositories"
	"enconder/application/services"
	"enconder/domain"
	"enconder/framework/database"
	"log"
	"os"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func prepare() (*domain.Video, repositories.VideoRepositoryDB){
	db := database.NewDbTest()
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("erro to get sql.DB: %v", err)
	}
	defer sqlDB.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "convite.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDB{DB: db}

	return video, repo
}

func TestVideoServiceDownload(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download(os.Getenv("R2_BUCKET"))
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)
}
