package users

import (
	"apous-films-rest-api/models"
	"apous-films-rest-api/test"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Test_Repository(t *testing.T) {
	asserts := assert.New(t)

	db := test.DB(t)
	test.ResetCollections(t, db, "users")
	repo := NewRepository(db)

	ctx := context.Background()

	// Count
	count, err := repo.Count(ctx)
	asserts.Nil(err)
	asserts.Equal(0, count)

	// Insert
	id, err := repo.Insert(ctx, &models.User{
		Email:        "example@mail.com",
		PasswordHash: "hashed_password",
		Provider:     "my_provider",
	})

	asserts.Nil(err)
	count2, _ := repo.Count(ctx)
	asserts.Equal(1, count2-count)

	// Get by Id
	user, err := repo.GetById(ctx, id)
	asserts.Nil(err)
	asserts.Equal("example@mail.com", user.Email)
	_, err = repo.GetById(ctx, primitive.NewObjectID().Hex())
	asserts.Equal(mongo.ErrNoDocuments, err)

	// Get By Email
	user, err = repo.GetByEmail(ctx, "example@mail.com")
	asserts.Nil(err)
	asserts.Equal(id, user.ID)
	_, err = repo.GetByEmail(ctx, "nobody@mail.com")
	asserts.Equal(mongo.ErrNoDocuments, err)
}
