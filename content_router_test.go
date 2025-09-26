package contentrouter_test

import (
	"context"
	"testing"

	"github.com/aomirun/content-router"
)

func TestContentRouterAPI(t *testing.T) {
	// 测试创建路由器
	router := contentrouter.NewRouter()
	if router == nil {
		t.Error("Failed to create router")
	}

	// 测试创建缓冲区
	buf := contentrouter.NewBuffer()
	if buf == nil {
		t.Error("Failed to create buffer")
	}

	// 测试创建上下文
	ctx := contentrouter.NewContext(context.Background(), buf)
	if ctx == nil {
		t.Error("Failed to create context")
	}

	// 验证类型兼容性
	var _ contentrouter.Router = router
	var _ contentrouter.Buffer = buf
	var _ contentrouter.Context = ctx
}