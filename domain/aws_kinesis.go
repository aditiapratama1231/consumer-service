package domain

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// KinesisStream , collection to AWS Kinesis Stream Connection
type KinesisStream struct {
	StreamName string
	Session    *session.Session
}

// Consume struct
type Consume struct {
	SequenceNumber *string
	Data           struct {
		Head struct {
			Service    string
			Domain     string
			ActionType uint16
			Action     string
		}
		Body struct {
			Payload map[string]interface{}
		}
	}
}
