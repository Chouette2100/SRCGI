package main

import (
	// "fmt"
	"log"
	// "net"
	// "net/http"
	// "strings"
	"sync"
	"time"
)

// SimpleRateLimiter はIPアドレスごとのリクエスト数をカウントし、レート制限を行う簡易リミッターです。
type SimpleRateLimiter struct {
	mu      sync.Mutex
	counts  map[string]int
	lastReq map[string]time.Time
	limit   int           // 許容リクエスト数
	window  time.Duration // 時間枠
}

// NewSimpleRateLimiter は新しいSimpleRateLimiterを作成します。
func NewSimpleRateLimiter(limit int, window time.Duration) *SimpleRateLimiter {
	rl := &SimpleRateLimiter{
		counts:  make(map[string]int),
		lastReq: make(map[string]time.Time),
		limit:   limit,
		window:  window,
	}
	// バックグラウンドで古いエントリを定期的にクリーンアップ
	go rl.cleanupLoop()
	return rl
}

// Allow は指定されたIPアドレスからのリクエストを許可するかどうかを判定します。
// 制限を超過している場合は false を返します。
func (rl *SimpleRateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// 時間枠内のリクエスト数をカウント
	lastTime, ok := rl.lastReq[ip]
	if !ok || now.Sub(lastTime) > rl.window {
		// 時間枠を過ぎたか、最初のリクエストの場合はカウントをリセット
		rl.counts[ip] = 1
		rl.lastReq[ip] = now
		return true // 許可
	}

	// 時間枠内のリクエスト数をインクリメント
	rl.counts[ip]++

	// 制限を超えているか判定
	if rl.counts[ip] > rl.limit {
		// 制限超過
		log.Printf("Rate limit exceeded for IP: %s (count: %d)", ip, rl.counts[ip])
		return false // 拒否
	}

	// 許可
	rl.lastReq[ip] = now // 最後のアクセス時刻を更新
	return true
}

// cleanupLoop は古いエントリをマップから定期的に削除します。
func (rl *SimpleRateLimiter) cleanupLoop() {
	// クリーンアップ間隔はウィンドウの半分など、適当な値に設定
	interval := rl.window / 2
	if interval < time.Second {
		interval = time.Second // 最低1秒
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cleanedCount := 0
		for ip, lastTime := range rl.lastReq {
			if now.Sub(lastTime) > rl.window {
				delete(rl.counts, ip)
				delete(rl.lastReq, ip)
				cleanedCount++
			}
		}
		// log.Printf("Rate limiter cleanup: removed %d old entries", cleanedCount) // デバッグ用
		rl.mu.Unlock()
	}
}
