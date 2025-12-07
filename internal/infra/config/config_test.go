package config

import "testing"

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config.Server.Host != "localhost" {
		t.Errorf("Expected host localhost, got %s", config.Server.Host)
	}

	if config.Server.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", config.Server.Port)
	}
}