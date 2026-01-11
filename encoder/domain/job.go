package domain

type Job struct {
  ID                string
  OutputBucketPath  string
  Status            string
  Video             *Video
}

