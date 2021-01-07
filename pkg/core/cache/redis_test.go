package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v7"
	"github.com/robinjoseph08/redisqueue/v2"
)

func TestRedis_Append(t *testing.T) {
	type fields struct {
		ConnectOption   *redis.Options
		ConsumerOptions *redisqueue.ConsumerOptions
		ProducerOptions *redisqueue.ProducerOptions
		client          *redis.Client
		consumer        *redisqueue.Consumer
		producer        *redisqueue.Producer
		mutex           *redislock.Client
	}
	type args struct {
		name    string
		message Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"test01",
			fields{
				ConnectOption: &redis.Options{
					Addr: "127.0.0.1:6379",
				},
				ConsumerOptions: &redisqueue.ConsumerOptions{
					VisibilityTimeout: 60 * time.Second,
					BlockingTimeout:   5 * time.Second,
					ReclaimInterval:   1 * time.Second,
					BufferSize:        100,
					Concurrency:       10,
				},
				ProducerOptions: &redisqueue.ProducerOptions{
					StreamMaxLength:      100,
					ApproximateMaxLength: true,
				},
			},
			args{
				name: "test",
				message: &RedisMessage{redisqueue.Message{
					Stream: "test",
					Values: map[string]interface{}{
						"key": "value",
					},
				}},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Redis{
				ConnectOption: tt.fields.ConnectOption,
			}
			if err := r.Connect(); err != nil {
				t.Errorf("SetQueue() error = %v", err)
			}
			if err := r.Append(tt.args.name, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("SetQueue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedis_Register(t *testing.T) {
	type fields struct {
		ConnectOption   *redis.Options
		ConsumerOptions *redisqueue.ConsumerOptions
		ProducerOptions *redisqueue.ProducerOptions
		client          *redis.Client
		consumer        *redisqueue.Consumer
		producer        *redisqueue.Producer
		mutex           *redislock.Client
	}
	type args struct {
		name string
		f    ConsumerFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test01",
			fields{
				ConnectOption: &redis.Options{
					Addr: "127.0.0.1:6379",
				},
				ConsumerOptions: &redisqueue.ConsumerOptions{
					VisibilityTimeout: 60 * time.Second,
					BlockingTimeout:   5 * time.Second,
					ReclaimInterval:   1 * time.Second,
					BufferSize:        100,
					Concurrency:       10,
				},
				ProducerOptions: &redisqueue.ProducerOptions{
					StreamMaxLength:      100,
					ApproximateMaxLength: true,
				},
			},
			args{
				name: "test",
				f: func(message Message) error {
					fmt.Println(message.GetValues())
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Redis{
				ConnectOption: tt.fields.ConnectOption,
				ConsumerOptions: &redisqueue.ConsumerOptions{
					VisibilityTimeout: 60 * time.Second,
					BlockingTimeout:   5 * time.Second,
					ReclaimInterval:   1 * time.Second,
					BufferSize:        100,
					Concurrency:       10,
				},
				ProducerOptions: &redisqueue.ProducerOptions{
					StreamMaxLength:      100,
					ApproximateMaxLength: true,
				},
			}
			if err := r.Connect(); err != nil {
				t.Error(err)
			}
			r.Register(tt.args.name, tt.args.f)
			r.Run()
		})
	}
}

func TestRedis_newConsumer(t *testing.T) {
	type fields struct {
		ConnectOption   *redis.Options
		ConsumerOptions *redisqueue.ConsumerOptions
		ProducerOptions *redisqueue.ProducerOptions
		client          *redis.Client
		consumer        *redisqueue.Consumer
		producer        *redisqueue.Producer
		mutex           *redislock.Client
	}
	type args struct {
		client *redis.Client
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"test01",
			fields{
				ConnectOption: &redis.Options{
					Addr: "127.0.0.1:6379",
				},
				ConsumerOptions: &redisqueue.ConsumerOptions{
					VisibilityTimeout: 60 * time.Second,
					BlockingTimeout:   5 * time.Second,
					ReclaimInterval:   1 * time.Second,
					BufferSize:        100,
					Concurrency:       10,
				},
				ProducerOptions: &redisqueue.ProducerOptions{
					StreamMaxLength:      100,
					ApproximateMaxLength: true,
				},
			},
			args{client: nil},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Redis{
				ConnectOption:   tt.fields.ConnectOption,
				ConsumerOptions: tt.fields.ConsumerOptions,
				ProducerOptions: tt.fields.ProducerOptions,
				client:          tt.fields.client,
				consumer:        tt.fields.consumer,
				producer:        tt.fields.producer,
				mutex:           tt.fields.mutex,
			}
			if err := r.Connect(); err != nil {
				t.Error(err)
				return
			}
			tt.args.client = r.client
			got, err := r.newConsumer(tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("newConsumer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.Register("test", func(message *redisqueue.Message) error {
				fmt.Println(message)
				return nil
			})
			got.Run()
		})
	}
}
