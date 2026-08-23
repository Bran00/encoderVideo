package services

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	return err
}
