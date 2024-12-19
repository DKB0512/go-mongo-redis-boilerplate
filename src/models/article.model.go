package models

import (
	"context"
	"fmt"
	"go-boilerplate/src/common"
	"go-boilerplate/src/core/db"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func ArticlesModel() *BaseModel {
	mod := &BaseModel{
		ModelConstructor: &common.ModelConstructor{
			Collection: db.GetMongoDb().Collection("ArticleCollection"),
		},
	}

	return mod
}

// models definitions
type Article struct {
	ID        uuid.UUID `json:"_id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	IsDeleted bool      `json:"is_deleted,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type CreateArticleForm struct {
	Title   string `form:"title" json:"title" binding:"required,min=3,max=100"`
	Content string `form:"content" json:"content" binding:"required,min=3,max=1000"`
}

type UpdateArticleForm struct {
	Title   string `form:"title" json:"title"`
	Content string `form:"content" json:"content"`
}

type FindArticleForm struct {
	ID uuid.UUID `form:"_id" json:"_id" binding:"required"`
}

// ArticleModel ...
type ArticleModel struct{}

func (mod *BaseModel) GetAllArticles() ([]Article, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := mod.Collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, fmt.Errorf("failed to find article: %w", err)
	}

	defer cursor.Close(ctx)

	var articles []Article
	if err := cursor.All(ctx, &articles); err != nil {
		return nil, fmt.Errorf("failed to decode article: %w", err)
	}

	return articles, err
}

func (mod *BaseModel) GetOneArticle(id uuid.UUID) (Article, error) {
	var article Article

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err := mod.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&article)

	if err != nil {
		return article, fmt.Errorf("failed to find article: %w", err)
	}

	return article, nil
}

func (mod *BaseModel) CreateArticle(form CreateArticleForm) (uuid.UUID, error) {

	id, err := uuid.NewV7()

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to find article: %w", err)
	}

	var article Article

	article.ID = id
	article.Title = form.Title
	article.Content = form.Content
	article.IsDeleted = false
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	_, err = mod.Collection.InsertOne(ctx, article)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create a new article: %v", err)
	}

	return id, nil
}

func (mod *BaseModel) UpdateArticle(id uuid.UUID, form UpdateArticleForm) error {
	_, err := mod.GetOneArticle(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "is_deleted": false}
	update := bson.M{
		"$set": bson.M{
			"title":      form.Title,
			"content":    form.Content,
			"updated_at": time.Now(),
		},
	}

	result, err := mod.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func (mod *BaseModel) DeleteArticle(id uuid.UUID) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now(),
		},
	}

	result, err := mod.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}
