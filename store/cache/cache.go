package cache

import (
	"context"
	"sync"
	"time"
)

type Config struct {
	DefaultTTL      time.Duration
	CleanupInterval time.Duration
	MaxItems        int
}

type Cache struct {
	config     Config               // 缓存配置
	items      map[string]cacheItem // 缓存项存储，键为字符串，值为缓存项
	mu         sync.RWMutex         // 读写锁，用于并发安全
	stopCh     chan struct{}        // 停止信号通道
	ctx        context.Context      // 上下文，用于控制生命周期
	cancelFunc context.CancelFunc   // 取消函数，用于关闭缓存
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

func New(config Config) *Cache {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Cache{
		config:     config,
		items:      make(map[string]cacheItem),
		stopCh:     make(chan struct{}),
		ctx:        ctx,
		cancelFunc: cancel,
	}
	go c.cleanupLoop()
	return c
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(c.config.DefaultTTL),
	}
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	if time.Now().After(item.expiresAt) {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Delete(ctx context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache) Close() {
	c.cancelFunc()
	close(c.stopCh)
}

func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(c.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for key, item := range c.items {
				if now.After(item.expiresAt) {
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		case <-c.stopCh:
			return
		case <-c.ctx.Done():
			return
		}
	}
}
