package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("should return error if JWT_SECRET is missing", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "")
		_, err := Load()
		if err == nil {
			t.Error("expected error, got nil")
		}
        if err.Error() != "JWT_SECRET environment variable is not set" {
            t.Errorf("expected error 'JWT_SECRET environment variable is not set', got '%v'", err)
        }
	})

	t.Run("should return error if JWT_SECRET is default", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "your-super-secret-key-change-in-production")
		_, err := Load()
		if err == nil {
			t.Error("expected error, got nil")
		}
        if err.Error() != "JWT_SECRET is set to the default insecure value. Please change it in production" {
            t.Errorf("expected error 'JWT_SECRET is set to the default insecure value. Please change it in production', got '%v'", err)
        }
	})

	t.Run("should load config successfully if JWT_SECRET is valid", func(t *testing.T) {
		t.Setenv("JWT_SECRET", "valid-secret-key")
		cfg, err := Load()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if cfg.JWT.Secret != "valid-secret-key" {
			t.Errorf("expected JWT secret 'valid-secret-key', got '%s'", cfg.JWT.Secret)
		}
	})
}
