package cache

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/robinjoseph08/redisqueue/v2"
)

func TestMemory_Append(t *testing.T) {
	type fields struct {
		items *sync.Map
		queue *sync.Map
		wait  sync.WaitGroup
		mutex sync.RWMutex
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
			fields{},
			args{
				name: "test",
				message: &MemoryMessage{redisqueue.Message{
					ID:     "",
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
			m := &Memory{
				items: tt.fields.items,
				queue: tt.fields.queue,
				wait:  tt.fields.wait,
				mutex: tt.fields.mutex,
			}
			if err := m.Connect(); err != nil {
				t.Errorf("Connect() error = %v", err)
				return
			}
			if err := m.Append(tt.args.name, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMemory_Register(t *testing.T) {
	log.SetFlags(19)
	type fields struct {
		items *sync.Map
		queue *sync.Map
		wait  sync.WaitGroup
		mutex sync.RWMutex
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
			fields{},
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
			m := &Memory{
				items: tt.fields.items,
				queue: tt.fields.queue,
				wait:  tt.fields.wait,
				mutex: tt.fields.mutex,
			}
			if err := m.Connect(); err != nil {
				t.Error(err)
				return
			}
			m.Register(tt.name, tt.args.f)
			if err := m.Append(tt.name, &MemoryMessage{redisqueue.Message{
				ID:     "",
				Stream: "test",
				Values: map[string]interface{}{
					"key": "value",
				},
			}}); err != nil {
				t.Error(err)
				return
			}
			go func() {
				m.Run()
			}()
			time.Sleep(3 * time.Second)
			m.Shutdown()
		})
	}
}

func TestMemory_Get(t *testing.T) {
	type fields struct {
		items   *sync.Map
		queue   *sync.Map
		wait    sync.WaitGroup
		mutex   sync.RWMutex
		PoolNum uint
	}
	type args struct {
		key    string
		value  string
		expire int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"test01",
			fields{},
			args{
				key:    "test",
				value:  "test",
				expire: 10,
			},
			"test",
			false,
		},
		{
			"test02",
			fields{},
			args{
				key:    "test",
				value:  "test1",
				expire: 1,
			},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				items:   tt.fields.items,
				queue:   tt.fields.queue,
				wait:    tt.fields.wait,
				mutex:   tt.fields.mutex,
				PoolNum: tt.fields.PoolNum,
			}
			if err := m.Connect(); err != nil {
				t.Errorf("Connect() error = %v", err)
				return
			}
			if err := m.Set(tt.args.key, tt.args.value, tt.args.expire); err != nil {
				t.Errorf("Set() error = %v", err)
				return
			}
			time.Sleep(2 * time.Second)
			got, err := m.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}
