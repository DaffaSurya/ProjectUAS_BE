package repository

import (
	model "PROJECTUAS_BE/app/Model"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	GetAll() ([]model.Achievement, error)
	GetAchievementByID(id string) (*model.Achievement, error)
	Create(ctx context.Context, achieve *model.Achievement) error
	FindById(ctx context.Context, id string) (*model.Achievement, error)
	Update(ctx context.Context, id string, update bson.M) error
	Delete(ctx context.Context, id string) error
	GetStudentByAchievement(ctx context.Context, studentID string) ([]*model.Achievement, error)
	AddAttachment(ctx context.Context, achievementID string, attachment model.Attachment) error
}

type AchievementMongoDB struct {
	Collection *mongo.Collection
}

func NewAchievementMongo(db *mongo.Database) AchievementRepository {
	return &AchievementMongoDB{
		Collection: db.Collection("achievements"),
	}
}

func (r *AchievementMongoDB) GetAll() ([]model.Achievement, error) {
	ctx := context.TODO()

	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []model.Achievement
	if err := cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *AchievementMongoDB) GetAchievementByID(id string) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var achievement model.Achievement

	// filter id bson
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}

	err = r.Collection.FindOne(ctx, filter).Decode(&achievement)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &achievement, nil
}

func (r *AchievementMongoDB) Create(ctx context.Context, achieve *model.Achievement) error {
	achieve.CreatedAt = time.Now()
	_, err := r.Collection.InsertOne(ctx, achieve)
	return err
}

func (r *AchievementMongoDB) FindById(ctx context.Context, id string) (*model.Achievement, error) {
	var achievement model.Achievement
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&achievement)
	if err != nil {
		return nil, err
	}
	return &achievement, nil
}

func (r *AchievementMongoDB) Update(ctx context.Context, id string, update bson.M) error {
	_, err := r.Collection.UpdateByID(ctx, id, bson.M{
		"$set": update,
	})
	return err
}

func (r *AchievementMongoDB) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.Collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *AchievementMongoDB) GetStudentByAchievement(ctx context.Context, studentID string) ([]*model.Achievement, error) {
	filter := bson.M{
		"studentId": studentID,
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var achievements []*model.Achievement
	if err := cursor.All(ctx, &achievements); err != nil {
		return nil, err
	}

	return achievements, nil
}

func (r *AchievementMongoDB) AddAttachment(ctx context.Context, achievementID string, attachment model.Attachment) error {
	update := bson.M{
		"$push": bson.M{
			"attachments": attachment,
		},
		"$set": bson.M{
			"updatedAt": attachment.UploadedAt,
		},
	}

	_, err := r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": achievementID},
		update,
	)

	return err
}
