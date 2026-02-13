package main

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

func main() {
	// 使用 Trac API 的 DSN
	// 替换为你的实际 DSN
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "http://4fbf7efef9b5b1d067dfac36d1cd891f@8.219.77.159:8025/1",
		// 可选配置
		Environment: "production",
		Release:     "my-app@1.0.0",
		Debug:       true, // 开启调试模式查看详细日志
	})
	if err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
		return
	}
	defer sentry.Flush(2 * time.Second)

	// 测试 1: 捕获简单消息
	fmt.Println("Test 1: Capturing a message...")
	sentry.CaptureMessage("Hello from Sentry Go SDK!")

	// 测试 2: 捕获异常
	fmt.Println("Test 2: Capturing an exception...")
	sentry.CaptureException(fmt.Errorf("test error: something went wrong"))

	// 测试 3: 带上下文的错误
	fmt.Println("Test 3: Capturing with context...")
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("test_type", "sdk_integration")
		scope.SetUser(sentry.User{
			ID:       "123",
			Email:    "test@example.com",
			Username: "testuser",
		})
		scope.SetContext("custom", map[string]interface{}{
			"request_id": "req-12345",
			"api_version": "v1",
		})
		sentry.CaptureMessage("Message with rich context")
	})

	// 测试 4: 模拟 panic 恢复
	fmt.Println("Test 4: Recovering from panic...")
	func() {
		defer sentry.Recover()
		panic("intentional panic for testing")
	}()

	fmt.Println("\nAll events sent! Check your Trac dashboard.")
	fmt.Println("Waiting for events to be flushed...")
	time.Sleep(3 * time.Second)
}
