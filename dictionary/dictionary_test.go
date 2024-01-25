package dictionary_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redis/v8"
)

func SetDefinition(ctx context.Context, rdb *redis.Client, word string, definition string) error {
	if len(word) == 0 {
		return errors.New("invalid word length")
	}
	if len(definition) == 0 {
		return errors.New("invalid definition length")
	}
	return rdb.Set(ctx, word, definition, 0).Err()
}

func TestAdd(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	err := SetDefinition(ctx, rdb, "test_word", "test_definition")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = SetDefinition(ctx, rdb, "", "test_definition")
	if err == nil {
		t.Error("Expected error for invalid word length, but got nil")
	}
}
func TestGet(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	err := rdb.Set(ctx, "test_word", "test_definition", 0).Err()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	definition, err := rdb.Get(ctx, "test_word").Result()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if definition != "test_definition" {
		t.Errorf("Expected definition 'test_definition', but got '%s'", definition)
	}

	_, err = rdb.Get(ctx, "nonexistent_word").Result()
	if err == nil {
		t.Error("Expected error for nonexistent word, but got nil")
	}
}

func TestRemove(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()

	err := rdb.Set(ctx, "test_word", "test_definition", 0).Err()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = rdb.Del(ctx, "test_word").Err()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	_, err = rdb.Get(ctx, "nonexistent_word").Result()
	if err == nil {
		t.Error("Expected error for nonexistent word, but got nil")
	}
}
