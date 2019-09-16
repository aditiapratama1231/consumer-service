package consumer

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/jinzhu/gorm"

	"magento-consumer-service/config"
	"magento-consumer-service/controller"
	"magento-consumer-service/domain"
)

type consumer struct {
	DB         *gorm.DB
	Controller controller.Controller
}

type Consumer interface {
	MainConsumer() error
	DecodeRecord(*kinesis.Record) error
}

func NewConsumer(db *gorm.DB, ctrl controller.Controller) Consumer {
	return &consumer{
		DB:         db,
		Controller: ctrl,
	}
}

// Consumer function to consume data from kinesis
func (cons *consumer) MainConsumer() error {
	ac := config.Kinesis()
	kc := kinesis.New(ac.Session)
	streamName := aws.String(ac.StreamName)
	streams, err := kc.DescribeStream(&kinesis.DescribeStreamInput{StreamName: streamName})
	if err != nil {
		return err
	}
	// // retrieve iterator
	iteratorOutput, err := kc.GetShardIterator(&kinesis.GetShardIteratorInput{
		// Shard Id is provided when making put record(s) request.
		ShardId:           aws.String(*streams.StreamDescription.Shards[0].ShardId),
		ShardIteratorType: aws.String("TRIM_HORIZON"),
		// ShardIteratorType: aws.String("AT_SEQUENCE_NUMBER"),
		// ShardIteratorType: aws.String("LATEST"),
		StreamName: streamName,
	})
	if err != nil {
		return err
	}

	shardIterator := iteratorOutput.ShardIterator
	var a *string
	for {
		// // get records use shard iterator for making request
		records, err := kc.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
		})
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		if len(records.Records) > 0 {
			for _, record := range records.Records {
				err := cons.DecodeRecord(record)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		} else if records.NextShardIterator == a || shardIterator == records.NextShardIterator || err != nil {
			log.Printf("GetRecords ERROR: %v\n", err)
			break
		}
		shardIterator = records.NextShardIterator
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (cons *consumer) DecodeRecord(record *kinesis.Record) error {
	data := &domain.Consume{SequenceNumber: record.SequenceNumber}
	err := json.Unmarshal([]byte(record.Data), &data.Data)
	if err != nil {
		return err
	}
	cons.Controller.MainController(data)
	return nil
}
