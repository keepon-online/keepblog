package config

import (
	"os"
	"testing"
)

func TestLoadAndGet(t *testing.T) {
	if err := Load(); err != nil {
		t.Fatalf("Load 报错: %v", err)
	}
	cfg := Get()
	if cfg.Http == nil {
		t.Fatal("Load 后 Http 配置不应为 nil（有默认值）")
	}
	if cfg.Http.Port == "" {
		t.Error("http.port 应有默认值")
	}

	// 幂等：重复 Load 不报错
	if err := Load(); err != nil {
		t.Fatalf("重复 Load 报错: %v", err)
	}
}

func TestValidateConfig_RequiresJwtSecret(t *testing.T) {
	// 测试进程可能继承了已加载的配置，通过环境变量保证 jwt.secret 存在与否可控
	t.Setenv("JWT_SECRET", "")
	_ = os.Unsetenv("JWT_SECRET")

	// Load 已在其他测试中执行；这里直接构造未加载场景验证零值配置被拒绝
	snapshot := cfgPtr.Swap(new(Configs))
	defer cfgPtr.Store(snapshot)

	if err := ValidateConfig(); err == nil {
		t.Error("jwt.secret 为空时 ValidateConfig 应报错")
	}
}

func TestValidateConfig_JwtSecretFromEnv(t *testing.T) {
	// Load 幂等（首个测试已加载），env 绑定在快照解析时生效，
	// 直接重建快照验证 JWT_SECRET → jwt.secret 的映射
	t.Setenv("JWT_SECRET", "test-secret-0123456789abcdef0123456")
	t.Cleanup(func() { _ = storeSnapshot() })
	if err := storeSnapshot(); err != nil {
		t.Fatalf("storeSnapshot 报错: %v", err)
	}
	if err := ValidateConfig(); err != nil {
		t.Errorf("提供了 JWT_SECRET 时不应报错: %v", err)
	}
	if Get().Jwt == nil || Get().Jwt.Secret != "test-secret-0123456789abcdef0123456" {
		t.Errorf("JWT_SECRET 未生效: %+v", Get().Jwt)
	}
}
