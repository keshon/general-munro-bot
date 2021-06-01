package storage

import (
	"time"

	"gorm.io/gorm"
)

type TaskRecord struct {
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

func CreateRecord(db *gorm.DB, taskID, taskUpdatedAt, taskStatus, commentID, commentUpdatedAt string) {
	db.Create(&TaskRecord{TaskID: taskID, TaskUpdatedAt: taskUpdatedAt, TaskStatus: taskStatus, CommentUpdatedAt: commentUpdatedAt, CommentID: commentID})
}

func UpdateRecord(db *gorm.DB, taskID, taskUpdatedAt, taskStatus, commentID, commentUpdatedAt string) {
	var rec TaskRecord
	db.Where("task_id=?", taskID).Find(&rec)
	rec.TaskUpdatedAt = taskUpdatedAt
	rec.TaskStatus = taskStatus
	rec.CommentUpdatedAt = commentUpdatedAt
	rec.CommentID = commentID
	db.Save(&rec)
}

func FindRecord(db *gorm.DB, taskID string) TaskRecord {
	var taskRecord TaskRecord
	db.First(&taskRecord, "task_id = ?", taskID) // find product with code D42
	return taskRecord
}
