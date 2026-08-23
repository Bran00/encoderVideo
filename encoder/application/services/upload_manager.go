package services

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/text/currency"
)

type VideoUpload struct {
	Paths 					[]string
	VideoPath 			string
	OutputBucket 		string
	Errors 					[]string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(
	objectPath string,
	client *s3.Client,
	ctx context.Context,
) error {
	path := strings.Split(
		objectPath,
		os.Getenv("localStoragePath")+"/",
	)

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(vu.OutputBucket),
		Key:    aws.String(path[1]),
		Body:   f,
	})

	return nil
}

func (vu *VideoUpload) loadPaths() error {

	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {

	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.loadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go vu.uploadWorker(in, <-returnChannel, uploadClient, ctx)

	}

	go func() {
		for x:= 0; x < len(vu.Paths); x++ {
			in <- x
		}
		close(in)
	}()

}

func (vu *VideoUpload) uploadWorker(in chan int, returnChan string, uploadClient *s3.Client, ctx context.Context) {

	for x := range in {
	
	}

}

func getClientUpload() (*s3.Client, context.Context, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, ctx, err
	}

	client := s3.NewFromConfig(cfg)
	return client, ctx, err
}
