package cache

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/robinjoseph08/redisqueue/v2"
	"github.com/spf13/cast"
)

type item struct {
	Value   string
	Expired time.Time
}

type queue chan Message

type Memory struct {
	items   *sync.Map
	queue   *sync.Map
	wait    sync.WaitGroup
	mutex   sync.RWMutex
	PoolNum uint
}

func (m *Memory) Connect() error {
	m.items = new(sync.Map)
	m.queue = new(sync.Map)
	return nil
}

func (m *Memory) makeQueue() queue {
	if m.PoolNum <= 0 {
		return make(queue)
	}
	return make(queue, m.PoolNum)
}

func (m *Memory) Get(key string) (string, error) {
	item, err := m.getItem(key)
	if err != nil {
		return "", err
	}
	return item.Value, nil
}

func (m *Memory) getItem(key string) (*item, error) {
	var err error
	i, ok := m.items.Load(key)
	if !ok {
		err = fmt.Errorf("%s not exist", key)
		return nil, err
	}
	switch i.(type) {
	case *item:
		item := i.(*item)
		if item.Expired.Before(time.Now()) {
			//过期
			_ = m.del(key)
			err = fmt.Errorf("%s is expired", key)
			//过期后删除
			return nil, err
		}
		return item, nil
	default:
		err = fmt.Errorf("value of %s type error", key)
		return nil, err
	}
}

func (m *Memory) Set(key string, val interface{}, expire int) error {
	s, err := cast.ToStringE(val)
	if err != nil {
		return err
	}
	item := &item{
		Value:   s,
		Expired: time.Now().Add(time.Duration(expire) * time.Second),
	}
	return m.setItem(key, item)
}

func (m *Memory) setItem(key string, item *item) error {
	m.items.Store(key, item)
	return nil
}

func (m *Memory) Del(key string) error {
	return m.del(key)
}

func (m *Memory) del(key string) error {
	m.items.Delete(key)
	return nil
}

func (m *Memory) HashGet(hk, key string) (string, error) {
	item, err := m.getItem(hk + key)
	if item == nil {
		err = fmt.Errorf("%s not exist", key)
		return "", err
	}
	return item.Value, err
}

func (m *Memory) HashDel(hk, key string) error {
	return m.del(hk + key)
}

func (m *Memory) Increase(key string) error {
	return m.calculate(key, 1)
}

func (m *Memory) Decrease(key string) error {
	return m.calculate(key, -1)
}

func (m *Memory) calculate(key string, num int) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	item, err := m.getItem(key)
	if err != nil {
		return err
	}

	if item == nil {
		err = fmt.Errorf("%s not exist", key)
		return err
	}
	var n int
	n, err = cast.ToIntE(item.Value)
	if err != nil {
		return err
	}
	n += num
	item.Value = strconv.Itoa(n)
	return m.setItem(key, item)
}

func (m *Memory) Expire(key string, dur time.Duration) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	item, err := m.getItem(key)
	if err != nil {
		return err
	}
	if item == nil {
		err = fmt.Errorf("%s not exist", key)
		return err
	}
	item.Expired = time.Now().Add(dur)
	return m.setItem(key, item)
}

func (m *Memory) Append(name string, message Message) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	memoryMessage := new(MemoryMessage)
	memoryMessage.SetID(message.GetID())
	memoryMessage.SetStream(message.GetStream())
	memoryMessage.SetValues(message.GetValues())
	v, ok := m.queue.Load(name)
	if !ok {
		v = m.makeQueue()
		m.queue.Store(name, v)
	}
	var q queue
	switch v.(type) {
	case queue:
		q = v.(queue)
	default:
		q = m.makeQueue()
		m.queue.Store(name, q)
	}
	go func(gm Message, gq queue) {
		gm.SetID(uuid.New().String())
		gq <- gm
	}(memoryMessage, q)
	return nil
}

func (m *Memory) Register(name string, f ConsumerFunc) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	v, ok := m.queue.Load(name)
	if !ok {
		v = m.makeQueue()
		m.queue.Store(name, v)
	}
	var q queue
	switch v.(type) {
	case queue:
		q = v.(queue)
	default:
		q = m.makeQueue()
		m.queue.Store(name, q)
	}
	go func(out queue, gf ConsumerFunc) {
		var err error
		for message := range q {
			err = gf(message)
			if err != nil {
				out <- message
				err = nil
			}
		}
	}(q, f)
}

func (m *Memory) Run() {
	m.wait.Add(1)
	m.wait.Wait()
}

func (m *Memory) Shutdown() {
	m.wait.Done()
}

type MemoryMessage struct {
	redisqueue.Message
}

func (m *MemoryMessage) GetID() string {
	return m.ID
}

func (m *MemoryMessage) GetStream() string {
	return m.Stream
}

func (m *MemoryMessage) GetValues() map[string]interface{} {
	return m.Values
}

func (m *MemoryMessage) SetID(id string) {
	m.ID = id
}

func (m *MemoryMessage) SetStream(stream string) {
	m.Stream = stream
}

func (m *MemoryMessage) SetValues(values map[string]interface{}) {
	m.Values = values
}
