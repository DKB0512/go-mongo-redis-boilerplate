package models

import (
	"context"
	"fmt"
	"go-boilerplate/src/common"
	"go-boilerplate/src/core/db"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

func UsersModel() *BaseModel {
	mod := &BaseModel{
		ModelConstructor: &common.ModelConstructor{
			Collection: db.GetMongoDb().Collection(string(db.UserCollection)),
		},
	}

	return mod
}

// models definitions
type User struct {
	ID        uuid.UUID `json:"_id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password,omitempty"`
	IsDeleted bool      `json:"is_deleted,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type CreateUserForm struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Email     string `form:"email" json:"email" binding:"required"`
	FirstName string `form:"first_name" json:"first_name" binding:"required"`
	LastName  string `form:"last_name" json:"last_name" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
}

type UsersResponse struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}
type UsersFindParam struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// models methods
func (mod *BaseModel) GetOneUser(userId uuid.UUID) User {
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err := mod.Collection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)

	if err != nil {
		fmt.Println(err)
		return user
	}

	return user
}

func (mod *BaseModel) GetAllUsers(limit int, skip int, search string) ([]User, int64, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	searchPattern := fmt.Sprintf("(?i)%s", regexp.QuoteMeta(search))

	regex := bson.M{"$regex": searchPattern}

	filter := bson.M{
		"$or": []bson.M{
			{"email": regex},
			{"username": regex},
			{"first_name": regex},
			{"last_name": regex},
		},
	}

	cursor, err := mod.Collection.Find(ctx, filter, options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)))

	if err != nil {
		fmt.Println(err)
		return nil, 0, fmt.Errorf("failed to find users: %w", err)
	}

	defer cursor.Close(ctx)

	var users []User
	if err := cursor.All(ctx, &users); err != nil {
		fmt.Println(err)
		return nil, 0, fmt.Errorf("failed to decode users: %w", err)
	}

	count, err := mod.Collection.CountDocuments(ctx, filter)

	if err != nil {
		fmt.Println(err)
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	return users, count, nil
}

func (mod *BaseModel) CreateUser(body User) (User, error) {
	id, err := uuid.NewV7()

	if err != nil {
		return User{}, fmt.Errorf("failed to generate a new id: %v", err)
	}

	body.ID = id
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, fmt.Errorf("failed to generate hash password: %v", err)
	}

	body.Password = string(hashedPassword)

	body.IsDeleted = false
	body.CreatedAt = time.Now()
	body.UpdatedAt = time.Now()

	// Set a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Insert the user and throw an error if any
	_, err = mod.Collection.InsertOne(ctx, body)

	if err != nil {
		return User{}, fmt.Errorf("failed to create a new user: %v", err)
	}

	return body, nil
}

func (mod *BaseModel) UpdateUser(param UsersFindParam, body User) (User, error) {
	body.ID = param.ID

	filter := bson.M{"_id": body.ID, "is_deleted": false}

	update := bson.M{"$set": body}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body.UpdatedAt = time.Now()

	// Perform the update operation
	result, err := mod.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return User{}, fmt.Errorf("failed to update user: %v", err)
	}

	if result.MatchedCount == 0 {
		return User{}, fmt.Errorf("no user found with ID: %s", body.ID)
	}

	return body, nil
}

func (mod *BaseModel) DeleteUser(param UsersFindParam) (bool, error) {
	var user User
	user.ID = param.ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := mod.Collection.FindOne(ctx, bson.M{"_id": user.ID, "isDeleted": false}).Decode(&user)

	//Fix this with proper logic
	if err != nil {
		return false, fmt.Errorf("failed to find the user: %v", err)
	}

	filter := bson.M{"_id": user.ID}

	update := bson.M{"$set": bson.M{"isDeleted": true}}

	result, err := mod.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, fmt.Errorf("failed to update user: %v", err)
	}

	if result.MatchedCount == 0 {
		return false, fmt.Errorf("no user found with ID: %s", user.ID)
	}

	return true, nil
}
