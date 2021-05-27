package storage

import (
	"time"

	"gorm.io/gorm"
)

type TaskRecord struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	TaskID         string
	TaskStatus     string
	TaskLastUpdate string
}

func CreateRecord(db *gorm.DB, taskID, taskStatus, lastUpdate string) {
	db.Create(&TaskRecord{TaskID: taskID, TaskStatus: taskStatus, TaskLastUpdate: lastUpdate})
}

func UpdateRecord(db *gorm.DB, taskID, taskStatus, lastUpdate string) {
	var taskRecord TaskRecord
	db.Where("task_id=?", taskID).Find(&taskRecord)
	taskRecord.TaskStatus = taskStatus
	taskRecord.TaskLastUpdate = lastUpdate
	db.Save(&taskRecord)
}

func FindRecord(db *gorm.DB, taskID string) TaskRecord {
	var taskRecord TaskRecord
	db.First(&taskRecord, "task_id = ?", taskID) // find product with code D42
	return taskRecord
}
