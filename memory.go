package limiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type clientRecord struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type MemoryLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientRecord
	rate    rate.Limit
	burst   int
}

// NewMemoryLimiter cria um limitador que roda inteiramente na memória RAM da instância
func NewMemoryLimiter(rate rate.Limit, burst int) *MemoryLimiter {
	mrl := &MemoryLimiter{
		clients: make(map[string]*clientRecord),
		rate:    rate,
		burst:   burst,
	}

	// Rotina de limpeza para evitar vazamento de memória (goroutine em background)
	go mrl.cleanupRoutine()

	return mrl
}

func (mrl *MemoryLimiter) cleanupRoutine() {
	for {
		time.Sleep(2 * time.Minute)
		mrl.mu.Lock()
		for ip, client := range mrl.clients {
			if time.Since(client.lastSeen) > 5*time.Minute {
				delete(mrl.clients, ip)
			}
		}
		mrl.mu.Unlock()
	}
}

func (mrl *MemoryLimiter) Allow(ip string) bool {
	mrl.mu.Lock()
	defer mrl.mu.Unlock()

	client, exists := mrl.clients[ip]
	if !exists {
		client = &clientRecord{
			limiter: rate.NewLimiter(mrl.rate, mrl.burst),
		}
		mrl.clients[ip] = client
	}
	client.lastSeen = time.Now()
	return client.limiter.Allow()
}
