package storage

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID               uint `gorm:"primaryKey"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	TaskID           string
	TaskUpdatedAt    string
	TaskStatus       string
	CommentID        string
	CommentUpdatedAt string
}

type Attachment struct {
	ID                  uint `gorm:"primaryKey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	AttachmentID        string
	AttachmentUpdatedAt string
	AttachmentStatus    string
}

func CreateTask(db *gorm.DB, taskID, taskUpdatedAt, taskStatus, commentID, commentUpdatedAt string) {
	db.Create(&Task{TaskID: taskID, TaskUpdatedAt: taskUpdatedAt, TaskStatus: taskStatus, CommentUpdatedAt: commentUpdatedAt, CommentID: commentID})
}

func UpdateTask(db *gorm.DB, taskID, taskUpdatedAt, taskStatus, commentID, commentUpdatedAt string) {
	var rec Task
	db.Where("task_id=?", taskID).Find(&rec)
	rec.TaskUpdatedAt = taskUpdatedAt
	rec.TaskStatus = taskStatus
	rec.CommentUpdatedAt = commentUpdatedAt
	rec.CommentID = commentID
	db.Save(&rec)
}

func FindTask(db *gorm.DB, taskID string) Task {
	var Task Task
	db.First(&Task, "task_id = ?", taskID) // find product with code D42
	return Task
}

func CreateAttachment(db *gorm.DB, attachmentID, attachmentUpdatedAt, attachmentStatus string) {
	db.Create(&Attachment{AttachmentID: attachmentID, AttachmentUpdatedAt: attachmentUpdatedAt, AttachmentStatus: attachmentStatus})
}

func UpdateAttachment(db *gorm.DB, attachmentID, attachmentUpdatedAt, attachmentStatus string) {
	var rec Attachment
	db.Where("attachment_id=?", attachmentID).Find(&rec)
	rec.AttachmentUpdatedAt = attachmentUpdatedAt
	rec.AttachmentStatus = attachmentStatus

	db.Save(&rec)
}

func FindAttachment(db *gorm.DB, attachmentID string) Attachment {
	var Attachment Attachment
	db.First(&Attachment, "attachment_id = ?", attachmentID) // find product with code D42
	return Attachment
}
