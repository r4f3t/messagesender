package message

import (
	"gorm.io/gorm"
)

type MessageRepository interface {
	FetchMessages(limit *int, isSent bool) ([]Message, error)
	MarkMessageAsSent(messageID uint, transactionDb *gorm.DB) error
	GetDBInstance() *gorm.DB
}

// repository model that constructor returns
type repository struct {
	DbInstance *gorm.DB
}

func NewRepository(db *gorm.DB) MessageRepository {
	return &repository{
		DbInstance: db,
	}
}

// FetchUnsentMessages retrieves unsent messages from the database
func (receiver *repository) FetchMessages(limit *int, isSent bool) ([]Message, error) {
	var messages []Message
	query := receiver.DbInstance.
		Where("sent = ?", isSent).
		Order("created_at ASC")

	if limit != nil {
		query = query.
			Limit(*limit)
	}

	result := query.
		Find(&messages)
	return messages, result.Error
}

// MarkMessageAsSent updates a message's Sent status in the database
func (receiver *repository) MarkMessageAsSent(messageID uint, transactionDb *gorm.DB) error {
	result := transactionDb.Model(&Message{}).Where("id = ?", messageID).Update("sent", true)
	return result.Error
}

// Represent DB Instance for transaction management outside of repository layer
func (receiver *repository) GetDBInstance() *gorm.DB {
	return receiver.DbInstance
}
