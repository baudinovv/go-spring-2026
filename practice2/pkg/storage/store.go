package storage

import "task-api/internal/models"

var Tasks []models.Task = make([]models.Task, 0)

// NextID is the auto-increment ID counter
var NextID int = 1
