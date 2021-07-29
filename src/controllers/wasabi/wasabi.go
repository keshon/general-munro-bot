package wasabi

import (
	"app/src/utils/config"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadFile(filename string, content string, conf config.Config) {
	bucket := aws.String(conf.Backup.S3.BucketName)
	key := aws.String(filename)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(conf.Backup.S3.AccessKey, conf.Backup.S3.SecretKey, ""),
		Endpoint:         aws.String(conf.Backup.S3.Endpoint),
		Region:           aws.String(conf.Backup.S3.Region),
		S3ForcePathStyle: aws.Bool(conf.Backup.S3.S3ForcePathStyle),
	}
	newSession := session.New(s3Config)

	s3Client := s3.New(newSession)

	_, err := s3Client.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(content),
		Bucket: bucket,
		Key:    key,
	})

	if err != nil {
		fmt.Printf("Failed to upload object %s%s, %s\n", *bucket, *key, err.Error())
		return
	}
	fmt.Printf("Successfully uploaded key %s\n", *key)
}
