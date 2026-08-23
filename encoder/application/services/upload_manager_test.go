package services_test

import (
	"enconder/application/services"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVideoServiceUpload(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download(os.Getenv("R2_BUCKET"))
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "codeeducationtest"
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(20, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "upload completed")
}
