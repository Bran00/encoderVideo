package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Job struct {
	ID                string		`valid:"uuid"`
	OutputBucketPath  string		`valid:"notnull"`
	Status            string		`valid:"notnull"`
	Video             *Video		`valid:"_"`
	VideoID						string		`valid:"_"`
  Error             string		`valid:"_"`
  CreatedAt         time.Time `valid:"_"`
  UpdatedAt         time.Time `valid:"_"`
}

func (job *Job) Validate() error {
	_, err := govalidator.ValidateStruct(job)

	if err != nil {
		return err
	}

	return nil
}
