package model

// AdminAchievementFilter - DTO untuk filter
type AdminAchievementFilter struct {
	Status          string `json:"status,omitempty"`          // 'draft', 'submitted', 'verified', 'rejected'
	AchievementType string `json:"achievement_type,omitempty"` // 'academic', 'competition', etc.
	StudentID       string `json:"student_id,omitempty"`      // Filter by student
	AdvisorID       string `json:"advisor_id,omitempty"`      // Filter by advisor
	ProgramStudy    string `json:"program_study,omitempty"`   // Filter by program study
}

// AdminAchievementSort - DTO untuk sorting
type AdminAchievementSort struct {
	Field string `json:"field"` // 'created_at', 'updated_at', 'status', 'student_name'
	Order string `json:"order"` // 'asc', 'desc'
}

// AchievementSummary - DTO untuk summary statistics
type AchievementSummary struct {
	Total     int `json:"total"`
	Draft     int `json:"draft"`
	Submitted int `json:"submitted"`
	Verified  int `json:"verified"`
	Rejected  int `json:"rejected"`
}