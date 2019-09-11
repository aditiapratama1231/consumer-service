package config

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"magento-consume-service/domain"
)

// AWSConnectKinesis function, used for connecting to aws-kinesis stream
func AWSConnectKinesis() domain.KinesisStream {
	var ac domain.KinesisStream
	stream := os.Getenv("AWS_KINESIS_STREAM_NAME")
	region := os.Getenv("AWS_KINESIS_REGION")
	endpoint := os.Getenv("AWS_KINESIS_ENDPOINT")
	accessKeyID := os.Getenv("AWS_KINESIS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_KINESIS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_KINESIS_SESSION_TOKEN")

	s := session.New(&aws.Config{
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, sessionToken),
	})

	ac.Session = s
	ac.StreamName = stream
	return ac
}
