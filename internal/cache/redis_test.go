package cache

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisClient(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectPing().SetVal("PONG")

	client := &RedisClient{client: db}

	ctx := context.Background()
	err := client.HealthCheck(ctx)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisClient_Set(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectSet("test_key", "test_value", 1*time.Hour).SetVal("OK")

	client := &RedisClient{client: db}

	ctx := context.Background()
	err := client.Set(ctx, "test_key", "test_value", 1*time.Hour)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisClient_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectGet("test_key").SetVal("test_value")

	client := &RedisClient{client: db}

	ctx := context.Background()
	val, err := client.Get(ctx, "test_key")
	assert.NoError(t, err)
	assert.Equal(t, "test_value", val)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisClient_HealthCheck(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectPing().SetVal("PONG")

	client := &RedisClient{client: db}

	ctx := context.Background()
	err := client.HealthCheck(ctx)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisClient_HealthCheck_Error(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectPing().SetErr(assert.AnError)

	client := &RedisClient{client: db}

	ctx := context.Background()
	err := client.HealthCheck(ctx)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
