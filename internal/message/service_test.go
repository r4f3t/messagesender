package message

import (
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) FetchMessages(limit *int, sent bool) ([]Message, error) {
	args := m.Called(limit, sent)
	return args.Get(0).([]Message), args.Error(1)
}

func (m *MockMessageRepository) MarkMessageAsSent(messageID uint, tx *gorm.DB) error {
	args := m.Called(messageID, tx)
	return args.Error(0)
}

func (m *MockMessageRepository) GetDBInstance() *gorm.DB {
	return nil // Not used in unit tests
}

func TestStartProcessingMessages(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockRedis := &redis.Client{}
	svc := NewService(mockRepo, mockRedis)

	mockRepo.On("FetchMessages", mock.Anything, false).Return([]Message{
		{ID: 1, To: "user1", Content: "Hello", Sent: false},
	}, nil)

	mockRepo.On("MarkMessageAsSent", 1, mock.Anything).Return(nil)

	err := svc.StartProcessingMessages()
	assert.NoError(t, err)

	time.Sleep(500 * time.Millisecond) // Allow goroutine to process messages

	err = svc.StopProcessingMessages()
	assert.NoError(t, err)
}

func TestStartProcessingMessages_AlreadyRunning(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockRedis := &redis.Client{}
	svc := NewService(mockRepo, mockRedis)

	err := svc.StartProcessingMessages()
	assert.NoError(t, err)

	err = svc.StartProcessingMessages()
	assert.EqualError(t, err, "Message processing is already running.")

	err = svc.StopProcessingMessages()
	assert.NoError(t, err)
}

func TestStopProcessingMessages_NotRunning(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockRedis := &redis.Client{}
	svc := NewService(mockRepo, mockRedis)

	err := svc.StopProcessingMessages()
	assert.EqualError(t, err, "Message processing is not running.")
}
