package services_test

import (
	"enconder/application/repositories"
	"enconder/domain"
	"enconder/framework/database"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
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
