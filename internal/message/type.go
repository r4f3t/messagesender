package message

import "time"

// Message represents the database structure for a message type
type Message struct {
	ID        uint   `gorm:"primaryKey"`
	To        string `gorm:"not null"`
	Content   string `gorm:"not null;size:255`
	Sent      bool   `gorm:"default:false"`
	CreatedAt time.Time
}
