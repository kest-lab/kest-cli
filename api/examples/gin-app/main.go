// Gin + Sentry SDK 集成示例
// 演示如何将错误上报到 Trac 服务器
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	// === Step 1: 配置 Sentry SDK ===
	// DSN 格式: http://{PUBLIC_KEY}@{HOST}/{PROJECT_ID}
	// 从 Trac 服务器获取: curl http://localhost:8025/v1/projects/1/dsn
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		// 使用默认测试 DSN（项目 ID: 1, Public Key: test_public_key）
		dsn = "http://test_public_key@localhost:8025/1"
	}

	log.Printf("🔧 初始化 Sentry SDK，DSN: %s\n", dsn)

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Debug:            true, // 开启调试模式
		AttachStacktrace: true, // 附加堆栈信息
		SampleRate:       1.0,  // 100% 采样（生产环境可调低）
		TracesSampleRate: 0.2,  // 20% 性能追踪采样
		Environment:      "development",
		Release:          "gin-example@1.0.0",
		ServerName:       "gin-demo-server",
	})
	if err != nil {
		log.Fatalf("sentry.Init 失败: %v", err)
	}
	defer sentry.Flush(2 * time.Second)

	log.Println("✅ Sentry SDK 初始化成功")

	// === Step 2: 创建 Gin 引擎 ===
	r := gin.Default()

	// === Step 3: 添加 Sentry 中间件 ===
	r.Use(sentrygin.New(sentrygin.Options{
		Repanic:         true,  // 重新抛出 panic
		WaitForDelivery: false, // 不阻塞请求
	}))

	// === Step 4: 定义路由 ===

	// 健康检查
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Gin + Sentry 示例应用",
			"dsn":     dsn,
			"routes": []string{
				"GET /           - 首页",
				"GET /hello      - 正常请求",
				"GET /error      - 手动捕获错误",
				"GET /panic      - 触发 panic",
				"GET /user/:id   - 带用户上下文的错误",
			},
		})
	})

	// 正常请求
	r.GET("/hello", func(c *gin.Context) {
		// 添加面包屑
		if hub := sentrygin.GetHubFromContext(c); hub != nil {
			hub.AddBreadcrumb(&sentry.Breadcrumb{
				Category: "action",
				Message:  "User visited /hello",
				Level:    sentry.LevelInfo,
			}, nil)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// 手动捕获错误
	r.GET("/error", func(c *gin.Context) {
		hub := sentrygin.GetHubFromContext(c)
		if hub != nil {
			// 设置标签
			hub.Scope().SetTag("endpoint", "/error")
			hub.Scope().SetTag("error_type", "manual")

			// 添加额外数据
			hub.Scope().SetExtra("request_id", "req-123")
			hub.Scope().SetExtra("timestamp", time.Now().Unix())

			// 捕获错误
			eventID := hub.CaptureException(fmt.Errorf("这是一个手动触发的测试错误"))
			log.Printf("📤 错误已发送到 Trac，Event ID: %s\n", *eventID)

			c.JSON(http.StatusOK, gin.H{
				"message":  "错误已发送到 Trac",
				"event_id": eventID,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取 Sentry Hub"})
	})

	// 触发 Panic（会被 Sentry 中间件自动捕获）
	r.GET("/panic", func(c *gin.Context) {
		hub := sentrygin.GetHubFromContext(c)
		if hub != nil {
			hub.Scope().SetTag("panic_type", "intentional")
			hub.Scope().SetFingerprint([]string{"panic", "test"})
		}
		log.Println("💥 即将触发 panic...")
		panic("这是一个测试 panic！系统崩溃了！")
	})

	// 带用户上下文的错误
	r.GET("/user/:id", func(c *gin.Context) {
		userID := c.Param("id")
		hub := sentrygin.GetHubFromContext(c)

		if hub != nil {
			// 设置用户信息
			hub.Scope().SetUser(sentry.User{
				ID:        userID,
				Email:     fmt.Sprintf("user%s@example.com", userID),
				Username:  fmt.Sprintf("user_%s", userID),
				IPAddress: c.ClientIP(),
			})

			// 设置上下文
			hub.Scope().SetContext("user_action", map[string]interface{}{
				"action":    "get_profile",
				"user_id":   userID,
				"timestamp": time.Now().Format(time.RFC3339),
			})

			// 添加面包屑
			hub.AddBreadcrumb(&sentry.Breadcrumb{
				Category: "user",
				Message:  fmt.Sprintf("Fetching user %s", userID),
				Level:    sentry.LevelInfo,
			}, nil)

			// 模拟用户不存在的错误
			if userID == "0" {
				hub.Scope().SetFingerprint([]string{"user", "not_found"})
				eventID := hub.CaptureException(fmt.Errorf("用户 %s 不存在", userID))
				c.JSON(http.StatusNotFound, gin.H{
					"error":    "用户不存在",
					"user_id":  userID,
					"event_id": eventID,
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"user_id":  userID,
			"username": fmt.Sprintf("user_%s", userID),
			"email":    fmt.Sprintf("user%s@example.com", userID),
		})
	})

	// === Step 5: 启动服务器 ===
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Gin 服务器启动在 http://localhost:%s\n", port)
	log.Println("📝 测试命令:")
	log.Println("   curl http://localhost:8080/         # 首页")
	log.Println("   curl http://localhost:8080/hello    # 正常请求")
	log.Println("   curl http://localhost:8080/error    # 触发错误")
	log.Println("   curl http://localhost:8080/panic    # 触发 panic")
	log.Println("   curl http://localhost:8080/user/0   # 用户不存在")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
