package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Job struct {
  ID                string
  OutputBucketPath  string
  Status            string
  Video             *Video
	VideoID						string
  Error             string
  CreatedAt         time.Time
  UpdatedAt         time.Time
}

func (job *Job) Validate() error {
	_, err := govalidator.ValidateStruct(job)

	if err != nil {
		return err
	}

	return nil
}
