package models

// KinesisSequenceNumber struct models, collection sequence number
type KinesisSequenceNumber struct {
	Model
	SequenceNumber string `gorm:"column:sequence_number"`
}
