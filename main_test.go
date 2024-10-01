package main

import (
	"os"
	"testing"
	"time"

	"github.com/atotto/clipboard"
	"github.com/pquerna/otp/totp"
)

type MockMenuItem struct {
	title   string
	tooltip string
	clicked chan struct{}
}

func (m *MockMenuItem) Click() {
	m.clicked <- struct{}{}
}

// Mock systray
type MockSystray struct {
	icon    []byte
	title   string
	tooltip string
	items   []*MockMenuItem
}

func (m *MockSystray) SetIcon(icon []byte) {
	m.icon = icon
}

func (m *MockSystray) SetTitle(title string) {
	m.title = title
}

func (m *MockSystray) SetTooltip(tooltip string) {
	m.tooltip = tooltip
}

func (m *MockSystray) AddMenuItem(title, tooltip string) *MockMenuItem {
	item := &MockMenuItem{
		title:   title,
		tooltip: tooltip,
		clicked: make(chan struct{}, 1),
	}
	m.items = append(m.items, item)
	return item
}

func TestIntegration(t *testing.T) {
	// Set up environment variables
	os.Setenv("TOTP_SECRET", "JBSWY3DPEHPK3PXP") // Test secret
	os.Setenv("TOTP_ICON_PATH", "testicon.ico")

	// Create a test icon file
	createTestIconFile(t)
	defer os.Remove("testicon.ico")

	// Create mock systray
	mockSystray := &MockSystray{}

	// Initialize the application with mock systray
	onReady := func() {
		mockSystray.SetIcon(icoToByteArray(os.Getenv("TOTP_ICON_PATH")))
		mockSystray.SetTitle("TOTP")
		mockSystray.SetTooltip("TOTP")

		copyMenuItem := mockSystray.AddMenuItem("Copy", "Copy the TOTP code to the clipboard")
		exitMenuItem := mockSystray.AddMenuItem("Exit", "Exit the application")

		go func() {
			for range copyMenuItem.clicked {
				secret := os.Getenv("TOTP_SECRET")
				otp, err := totp.GenerateCode(secret, time.Now())
				if err == nil {
					clipboard.WriteAll(otp)
				}
			}
		}()

		go func() {
			<-exitMenuItem.clicked
			// In a real application, this would call systray.Quit()
		}()
	}

	// Call onReady to set up the mock systray
	onReady()

	// Test TOTP code generation
	secret := os.Getenv("TOTP_SECRET")
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Errorf("Error generating TOTP code: %v", err)
	}
	if len(code) != 6 {
		t.Errorf("Generated TOTP code has incorrect length: %v", code)
	}

	// Test clipboard copy
	mockSystray.items[0].Click()       // Simulate "Copy" menu item click
	time.Sleep(100 * time.Millisecond) // Wait for copy operation to complete

	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		t.Errorf("Error reading clipboard: %v", err)
	}
	if len(clipboardContent) != 6 {
		t.Errorf("Clipboard content is incorrect: %v", clipboardContent)
	}

	// Test exit menu (just for coverage, as it doesn't do much in this mock setup)
	mockSystray.items[1].Click() // Simulate "Exit" menu item click
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
func createTestIconFile(t *testing.T) {
	content := []byte{0, 0, 1, 0, 1, 0, 16, 16, 0, 0, 1, 0, 32, 0, 68, 4, 0, 0, 22, 0, 0, 0}
	err := os.WriteFile("testicon.ico", content, 0644)
	if err != nil {
		t.Fatalf("テストアイコンファイルの作成に失敗しました: %v", err)
	}
}
