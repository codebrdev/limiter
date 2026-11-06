package limiter

import (
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestMemoryLimiter_Allow(t *testing.T) {
	// 1. Configura o limitador: máximo de 3 requisições por 1 minuto
	rateLimit := 3
	window := 1 * time.Minute
	limiter := NewMemoryLimiter(rate.Limit(float64(rateLimit)/window.Seconds()), rateLimit)

	ip := "192.168.1.50"

	// =========================================================================
	// CENÁRIO 1: As primeiras 3 requisições devem ser permitidas (Allow = true)
	// =========================================================================
	for i := 1; i <= rateLimit; i++ {
		allowed := limiter.Allow(ip)
		if !allowed {
			t.Errorf("tentativa %d: esperava que a requisição fosse permitida, mas foi bloqueada", i)
		}
	}

	// =========================================================================
	// CENÁRIO 2: A 4ª requisição deve estourar o limite e ser bloqueada (Allow = false)
	// =========================================================================
	allowed := limiter.Allow(ip)
	if allowed {
		t.Error("tentativa 4: esperava que a requisição fosse bloqueada por excesso de limite, mas passou")
	}

	// =========================================================================
	// CENÁRIO 3: Outro IP diferente deve ter seu próprio contador isolado e passar
	// =========================================================================
	otherIP := "10.0.0.99"
	allowedOther := limiter.Allow(otherIP)
	if !allowedOther {
		t.Error("esperava que um IP diferente tivesse contador isolado e fosse permitido, mas foi bloqueado")
	}
}
