package model

import "time"

type Attachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FileUrl    string    `bson:"fileUrl" json:"fileUrl"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}

type CompetitionDetails struct {
	CompetitionName  string    `bson:"competitionName" json:"competitionName"`
	CompetitionLevel string    `bson:"competitionLevel" json:"competitionLevel"`
	Rank             int       `bson:"rank" json:"rank"`
	MedalType        string    `bson:"medalType" json:"medalType"`
	EventDate        time.Time `bson:"eventDate" json:"eventDate"`
	Location         string    `bson:"location" json:"location"`
	Organizer        string    `bson:"organizer" json:"organizer"`
}

type Achievement struct {
	ID              string `bson:"_id,omitempty" json:"id"`
	StudentID       string `bson:"studentId" json:"studentId"`
	AchievementType string `bson:"achievementType" json:"achievementType"`
	Title           string `bson:"title" json:"title"`
	Description     string `bson:"description" json:"description"`

	Details     CompetitionDetails `bson:"details" json:"details"`
	Attachments []Attachment       `bson:"attachments" json:"attachments"`
	Tags        []string           `bson:"tags" json:"tags"`
	Points      int                `bson:"points" json:"points"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
