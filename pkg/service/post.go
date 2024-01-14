package service

import (
	"context"
	"fmt"
	postProto "github.com/tom-blog-app/blog-proto/post"
	"github.com/tom-blog-app/post-service/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"time"
)

var collection = os.Getenv("POST_MONGO_DB")

type PostServer struct {
	postProto.UnimplementedPostServiceServer
	Client *mongo.Client
}

func (s *PostServer) CreatePost(ctx context.Context, req *postProto.PostRequest) (*postProto.PostResponse, error) {

	postModel := models.Post{
		ID:        primitive.NewObjectID().Hex(),
		Title:     req.GetTitle(),
		Content:   req.GetContent(),
		AuthorID:  req.GetAuthorId(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert the post into the database
	collection := s.Client.Database(collection).Collection(collection)
	result, err := collection.InsertOne(ctx, postModel)

	if err != nil {
		return nil, fmt.Errorf("could not insert post: %v", err)
	}

	id := result.InsertedID.(string)

	var insertedPost models.Post
	err = collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&insertedPost)
	if err != nil {
		return nil, fmt.Errorf("could not find inserted post: %v", err)
	}

	// Return a successful response
	res := &postProto.PostResponse{
		Post: &postProto.Post{
			Id:        id,
			Title:     insertedPost.Title,
			Content:   insertedPost.Content,
			AuthorId:  insertedPost.AuthorID,
			CreatedAt: timestamppb.New(insertedPost.CreatedAt),
			UpdatedAt: timestamppb.New(insertedPost.UpdatedAt),
		},
	}
	return res, nil
}

func (s *PostServer) GetPost(ctx context.Context, req *postProto.GetPostRequest) (*postProto.PostResponse, error) {
	log.Println("GetPost request received")
	collection := s.Client.Database(collection).Collection(collection)
	var postModel models.Post
	err := collection.FindOne(ctx, bson.M{"_id": req.GetId()}).Decode(&postModel)
	if err != nil {
		return nil, fmt.Errorf("could not find post: %v", err)
	}

	post := &postProto.Post{
		Id:        postModel.ID,
		Title:     postModel.Title,
		Content:   postModel.Content,
		AuthorId:  postModel.AuthorID,
		CreatedAt: timestamppb.New(postModel.CreatedAt),
		UpdatedAt: timestamppb.New(postModel.UpdatedAt),
	}

	return &postProto.PostResponse{Post: post}, nil
}

func (s *PostServer) UpdatePost(ctx context.Context, req *postProto.UpdatePostRequest) (*postProto.PostResponse, error) {
	//post := req.GetPost()
	id := req.GetId()

	postModel := models.Post{
		ID:       id,
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
		AuthorID: req.GetAuthorId(),
		//CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := s.Client.Database(collection).Collection(collection)
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": postModel})
	if err != nil {
		return nil, fmt.Errorf("could not update post: %v", err)
	}

	var updatedPost models.Post
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedPost)
	if err != nil {
		return nil, fmt.Errorf("could not find inserted post: %v", err)
	}

	return &postProto.PostResponse{
		Post: &postProto.Post{
			Id:        updatedPost.ID,
			Title:     updatedPost.Title,
			Content:   updatedPost.Content,
			AuthorId:  updatedPost.AuthorID,
			CreatedAt: timestamppb.New(updatedPost.CreatedAt),
			UpdatedAt: timestamppb.New(updatedPost.UpdatedAt),
		},
	}, nil
}

func (s *PostServer) DeletePost(ctx context.Context, req *postProto.GetPostRequest) (*postProto.PostDeleteResponse, error) {
	collection := s.Client.Database(collection).Collection(collection)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": req.GetId()})
	if err != nil {
		return nil, fmt.Errorf("could not delete post: %v", err)
	}

	return &postProto.PostDeleteResponse{
		Id:      req.GetId(),
		Success: true,
	}, nil
}

func (s *PostServer) ListPosts(ctx context.Context, req *postProto.GetPostListRequest) (*postProto.GetPostsListResponse, error) {
	collection := s.Client.Database(collection).Collection(collection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("could not list posts: %v", err)
	}
	defer cursor.Close(ctx)

	var posts []*postProto.Post
	for cursor.Next(ctx) {
		var postModel models.Post
		err := cursor.Decode(&postModel)
		if err != nil {
			return nil, fmt.Errorf("could not decode post: %v", err)
		}

		post := &postProto.Post{
			Id:        postModel.ID,
			Title:     postModel.Title,
			Content:   postModel.Content,
			AuthorId:  postModel.AuthorID,
			CreatedAt: timestamppb.New(postModel.CreatedAt),
			UpdatedAt: timestamppb.New(postModel.UpdatedAt),
		}
		posts = append(posts, post)
	}

	return &postProto.GetPostsListResponse{Posts: posts}, nil
}

func (s *PostServer) ListPostsByAuthor(ctx context.Context, req *postProto.GetPostListByAuthorRequest) (*postProto.GetPostsListResponse, error) {
	// Create a filter that matches posts where the AuthorId field is equal to the author ID from the request
	filter := bson.M{"author_id": req.GetAuthorId()}

	collection := s.Client.Database(collection).Collection(collection)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("could not list posts: %v", err)
	}
	defer cursor.Close(ctx)

	var posts []*postProto.Post
	for cursor.Next(ctx) {
		var postModel models.Post
		err := cursor.Decode(&postModel)
		if err != nil {
			return nil, fmt.Errorf("could not decode post: %v", err)
		}

		post := &postProto.Post{
			Id:        postModel.ID,
			Title:     postModel.Title,
			Content:   postModel.Content,
			AuthorId:  postModel.AuthorID,
			CreatedAt: timestamppb.New(postModel.CreatedAt),
			UpdatedAt: timestamppb.New(postModel.UpdatedAt),
		}
		posts = append(posts, post)
	}

	return &postProto.GetPostsListResponse{Posts: posts}, nil
}
