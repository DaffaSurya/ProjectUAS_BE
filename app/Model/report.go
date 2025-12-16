package model

import "time"

type AchievementStatistics struct {
	Total     int            `json:"total"`
	Submitted int            `json:"submitted"`
	Verified  int            `json:"verified"`
	Rejected  int            `json:"rejected"`
	ByType    map[string]int `json:"by_type"`
}

type StudentAchievementItem struct {
	MongoAchievementID string     `json:"mongo_achievement_id"`
	Title              string     `json:"title"`
	AchievementType    string     `json:"achievement_type"`
	Status             string     `json:"status"`
	SubmittedAt        *time.Time `json:"submitted_at"`
	VerifiedAt         *time.Time `json:"verified_at"`
}

type StudentReport struct {
	StudentID    string                   `json:"student_id"`
	Total        int                      `json:"total"`
	Submitted    int                      `json:"submitted"`
	Verified     int                      `json:"verified"`
	Rejected     int                      `json:"rejected"`
	Achievements []StudentAchievementItem `json:"achievements"`
}
