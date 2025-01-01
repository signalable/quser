// quser/internal/repository/mongodb/user_repository.go
package mongodb

import (
	"context"
	"time"

	"github.com/signalable/quser/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewUserRepository MongoDB 유저 레포지토리 생성자
func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{
		db:         db,
		collection: db.Collection("users"),
	}
}

// Create 새로운 사용자 생성
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Status = domain.UserStatusPending

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByEmail 이메일로 사용자 찾기
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}

// FindByID ID로 사용자 찾기
func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrUserNotFound
	}
	return &user, err
}

// ExistsByEmail 이메일 존재 여부 확인
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update 사용자 정보 업데이트
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

// UpdateProfile 프로필 업데이트
func (r *userRepository) UpdateProfile(ctx context.Context, userID string, profile *domain.UserProfile) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"profile":    profile,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}

// UpdateVerificationStatus 이메일 인증 상태 업데이트
func (r *userRepository) UpdateVerificationStatus(ctx context.Context, userID string, isVerified bool) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"is_verified": isVerified,
				"status":      domain.UserStatusActive,
				"updated_at":  time.Now(),
			},
		},
	)
	return err
}

// UpdateStatus 사용자 상태 업데이트
func (r *userRepository) UpdateStatus(ctx context.Context, userID string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)
	return err
}
