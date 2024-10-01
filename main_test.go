package main

import (
	"os"
	"testing"
	"time"

	"github.com/getlantern/systray"
	"github.com/pquerna/otp/totp"
)

func TestIntegration(t *testing.T) {
	os.Setenv("TOTP_ICON_PATH", "test_icon.ico")
	os.Setenv("TOTP_SECRET", "JBSWY3DPEHPK3PXP")

	iconPath := "test_icon.ico"
	file, err := os.Create(iconPath)
	if err != nil {
		t.Fatalf("Error creating test icon file: %v", err)
	}
	file.Write(icoToByteArray("icon.ico"))
	file.Close()
	defer os.Remove(iconPath)

	go func() {
		systray.Run(onReady, onExit)
	}()

	time.Sleep(2 * time.Second)

	// TOTPコードの生成をテスト
	secret := os.Getenv("TOTP_KEY")
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Fatalf("Error generating TOTP code: %v", err)
	}

	if len(code) != 6 {
		t.Errorf("Expected TOTP code length of 6, but got %d", len(code))
	}

	// systrayの終了
	systray.Quit()
}

func TestIcoToByteArray(t *testing.T) {
	iconPath := "test_icon.ico"
	file, err := os.Create(iconPath)
	if err != nil {
		t.Fatalf("Error creating test icon file: %v", err)
	}
	expectedData := []byte{
		0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x20, 0x20, 0x00, 0x00, 0x01, 0x00,
		0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	file.Write(expectedData)
	file.Close()
	defer os.Remove(iconPath)

	// 関数をテスト
	data := icoToByteArray(iconPath)
	if data == nil {
		t.Errorf("Expected non-nil byte array, but got nil")
	}

	if len(data) != len(expectedData) {
		t.Errorf("Expected byte array length %d, but got %d", len(expectedData), len(data))
	}

	for i := range data {
		if data[i] != expectedData[i] {
			t.Errorf("Expected byte array element %d to be %v, but got %v", i, expectedData[i], data[i])
		}
	}
}
