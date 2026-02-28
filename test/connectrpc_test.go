package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

// ConnectrpcBaseURL 测试服务器基础 URL
var ConnectrpcBaseURL = getConnectrpcBaseURL()

// getConnectrpcBaseURL 从环境变量获取基础 URL，若未设置则使用默认值
func getConnectrpcBaseURL() string {
	if url := os.Getenv("TEST_BASE_URL"); url != "" {
		return url
	}
	return "http://localhost:8081"
}

// ConnectrpcUnifiedResponse represents the expected unified response format for Connect RPC tests
type ConnectrpcUnifiedResponse struct {
	State   int         `json:"state"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// TestConnectRPCRegisterUser tests user registration using Connect RPC endpoint
func TestConnectRPCRegisterUser(t *testing.T) {
	uniqueUsername := fmt.Sprintf("connectrpcuser_%d", time.Now().UnixNano())
	uniqueEmail := uniqueUsername + "@example.com"

	user := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Connect RPC Test User",
		"password": "testpassword123",
		"phone":    "13800138000",
		"email":    uniqueEmail,
	}
	body, _ := json.Marshal(user)

	resp, err := http.Post(ConnectrpcBaseURL+"/goserver.api.v1.UserService/RegisterUser", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to call Connect RPC RegisterUser: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status 200 or 409, got %d", resp.StatusCode)
	}

	// Connect RPC returns raw protobuf response, not unified JSON format
	// We'll check the status code only for now
	t.Logf("Connect RPC RegisterUser response status: %d", resp.StatusCode)
}

// TestConnectRPCLoginUser tests user login using Connect RPC endpoint
func TestConnectRPCLoginUser(t *testing.T) {
	// First register a user using REST API
	uniqueUsername := fmt.Sprintf("connectrpclogin_%d", time.Now().UnixNano())
	registerUser := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Connect RPC Login Test User",
		"password": "testpassword123",
		"phone":    "13800138001",
		"email":    uniqueUsername + "@example.com",
	}
	registerBody, _ := json.Marshal(registerUser)

	registerResp, err := http.Post(ConnectrpcBaseURL+"/api/v1/users", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		t.Fatalf("Failed to call RegisterUser: %v", err)
	}
	registerResp.Body.Close()

	// Then login with the registered user using Connect RPC
	loginReq := map[string]string{
		"username": uniqueUsername,
		"password": "testpassword123",
	}
	loginBody, _ := json.Marshal(loginReq)

	resp, err := http.Post(ConnectrpcBaseURL+"/goserver.api.v1.AuthService/Login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call Connect RPC LoginUser: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}

	// Connect RPC returns raw protobuf response, not unified JSON format
	// We'll check the status code only for now
	t.Logf("Connect RPC LoginUser response status: %d", resp.StatusCode)
}

// TestConnectRPCGetUserProfile tests getting user profile using Connect RPC endpoint
func TestConnectRPCGetUserProfile(t *testing.T) {
	body := []byte(`{}`)
	req, _ := http.NewRequest("POST", ConnectrpcBaseURL+"/goserver.api.v1.UserService/GetUserProfile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call Connect RPC GetUserProfile: %v", err)
	}
	defer resp.Body.Close()

	// Connect RPC should return 401 for unauthorized requests
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}

	t.Logf("Connect RPC GetUserProfile response status: %d", resp.StatusCode)
}

// TestConnectRPCRefreshToken tests token refresh using Connect RPC endpoint
func TestConnectRPCRefreshToken(t *testing.T) {
	// First register and login to get tokens
	uniqueUsername := fmt.Sprintf("connectrpcrefresh_%d", time.Now().UnixNano())

	// Register user
	registerUser := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Connect RPC Refresh Test User",
		"password": "testpassword123",
		"phone":    "13800138002",
		"email":    uniqueUsername + "@example.com",
	}
	registerBody, _ := json.Marshal(registerUser)

	registerResp, err := http.Post(ConnectrpcBaseURL+"/api/v1/users", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		t.Fatalf("Failed to call RegisterUser: %v", err)
	}
	registerResp.Body.Close()

	// Login to get tokens
	loginReq := map[string]string{
		"username": uniqueUsername,
		"password": "testpassword123",
	}
	loginBody, _ := json.Marshal(loginReq)

	loginResp, err := http.Post(ConnectrpcBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for login, got %d", loginResp.StatusCode)
		return
	}

	// Parse login response to get refresh token
	var unifiedResp ConnectrpcUnifiedResponse
	if err := json.NewDecoder(loginResp.Body).Decode(&unifiedResp); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	data, ok := unifiedResp.Data.(map[string]interface{})
	if !ok {
		t.Error("Expected data field to be an object")
		return
	}

	refreshToken, ok := data["refresh_token"].(string)
	if !ok {
		// Try refreshToken (camelCase) if refresh_token (snake_case) not found
		refreshToken, ok = data["refreshToken"].(string)
		if !ok {
			t.Error("Expected refresh_token or refreshToken in login response")
			return
		}
	}

	// Now test Connect RPC refresh token endpoint
	refreshReq := map[string]string{
		"refresh_token": refreshToken,
	}
	refreshBody, _ := json.Marshal(refreshReq)

	resp, err := http.Post(ConnectrpcBaseURL+"/goserver.api.v1.AuthService/RefreshToken", "application/json", bytes.NewReader(refreshBody))
	if err != nil {
		t.Fatalf("Failed to call Connect RPC RefreshToken: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}

	t.Logf("Connect RPC RefreshToken response status: %d", resp.StatusCode)
}

// TestConnectRPCGetUserProfileWithAuth tests getting user profile with authentication using Connect RPC endpoint
func TestConnectRPCGetUserProfileWithAuth(t *testing.T) {
	// First register and login to get tokens
	uniqueUsername := fmt.Sprintf("connectrpcauth_%d", time.Now().UnixNano())
	testPassword := "testpassword123"

	// Register user
	registerUser := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Connect RPC Auth Test User",
		"password": testPassword,
		"phone":    "13800138000",
		"email":    uniqueUsername + "@example.com",
	}
	registerBody, _ := json.Marshal(registerUser)

	registerResp, err := http.Post(ConnectrpcBaseURL+"/api/v1/users", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		t.Fatalf("Failed to call RegisterUser: %v", err)
	}
	registerResp.Body.Close()

	// Login to get access token
	loginReq := map[string]string{
		"username": uniqueUsername,
		"password": testPassword,
	}
	loginBody, _ := json.Marshal(loginReq)

	loginResp, err := http.Post(ConnectrpcBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for login, got %d", loginResp.StatusCode)
		return
	}

	// Parse login response to get access token
	var loginUnifiedResp ConnectrpcUnifiedResponse
	if err := json.NewDecoder(loginResp.Body).Decode(&loginUnifiedResp); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	data, ok := loginUnifiedResp.Data.(map[string]interface{})
	if !ok {
		t.Error("Expected data field to be an object")
		return
	}

	accessToken, ok := data["access_token"].(string)
	if !ok {
		// Try accessToken (camelCase) if access_token (snake_case) not found
		accessToken, ok = data["accessToken"].(string)
		if !ok {
			t.Error("Expected access_token or accessToken in login response")
			return
		}
	}

	// Now get user profile with authentication using Connect RPC
	body := []byte(`{}`)
	req, _ := http.NewRequest("POST", ConnectrpcBaseURL+"/goserver.api.v1.UserService/GetUserProfile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call Connect RPC GetUserProfile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	t.Logf("Connect RPC GetUserProfile with auth response status: %d", resp.StatusCode)
}

// TestConnectRPCUpdateUserProfile tests updating user profile using Connect RPC endpoint
func TestConnectRPCUpdateUserProfile(t *testing.T) {
	// First register and login to get tokens
	uniqueUsername := fmt.Sprintf("connectrpcupdate_%d", time.Now().UnixNano())
	testPassword := "testpassword123"

	// Register user
	registerUser := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Original Connect RPC User",
		"password": testPassword,
		"phone":    "13800138000",
		"email":    uniqueUsername + "@example.com",
	}
	registerBody, _ := json.Marshal(registerUser)

	registerResp, err := http.Post(ConnectrpcBaseURL+"/api/v1/users", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		t.Fatalf("Failed to call RegisterUser: %v", err)
	}
	registerResp.Body.Close()

	// Login to get access token
	loginReq := map[string]string{
		"username": uniqueUsername,
		"password": testPassword,
	}
	loginBody, _ := json.Marshal(loginReq)

	loginResp, err := http.Post(ConnectrpcBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for login, got %d", loginResp.StatusCode)
		return
	}

	// Parse login response to get access token
	var loginUnifiedResp ConnectrpcUnifiedResponse
	if err := json.NewDecoder(loginResp.Body).Decode(&loginUnifiedResp); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	data, ok := loginUnifiedResp.Data.(map[string]interface{})
	if !ok {
		t.Error("Expected data field to be an object")
		return
	}

	accessToken, ok := data["access_token"].(string)
	if !ok {
		// Try accessToken (camelCase) if access_token (snake_case) not found
		accessToken, ok = data["accessToken"].(string)
		if !ok {
			t.Error("Expected access_token or accessToken in login response")
			return
		}
	}

	// Now update user profile using Connect RPC
	updateUser := map[string]interface{}{
		"nickname": "Updated Connect RPC User",
		"phone":    "13900139000",
		"email":    "updated-" + uniqueUsername + "@example.com",
	}
	updateBody, _ := json.Marshal(updateUser)

	req, _ := http.NewRequest("POST", ConnectrpcBaseURL+"/goserver.api.v1.UserService/UpdateUserProfile", bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call Connect RPC UpdateUserProfile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	t.Logf("Connect RPC UpdateUserProfile response status: %d", resp.StatusCode)
}

// TestConnectRPCChangePassword tests changing user password using Connect RPC endpoint
func TestConnectRPCChangePassword(t *testing.T) {
	// First register and login to get tokens
	uniqueUsername := fmt.Sprintf("connectrpcchangepass_%d", time.Now().UnixNano())
	oldPassword := "oldpassword123"
	newPassword := "newpassword123"

	// Register user
	registerUser := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Connect RPC Change Password Test User",
		"password": oldPassword,
		"phone":    "13800138000",
		"email":    uniqueUsername + "@example.com",
	}
	registerBody, _ := json.Marshal(registerUser)

	registerResp, err := http.Post(ConnectrpcBaseURL+"/api/v1/users", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		t.Fatalf("Failed to call RegisterUser: %v", err)
	}
	registerResp.Body.Close()

	// Login to get access token
	loginReq := map[string]string{
		"username": uniqueUsername,
		"password": oldPassword,
	}
	loginBody, _ := json.Marshal(loginReq)

	loginResp, err := http.Post(ConnectrpcBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for login, got %d", loginResp.StatusCode)
		return
	}

	// Parse login response to get access token
	var loginUnifiedResp ConnectrpcUnifiedResponse
	if err := json.NewDecoder(loginResp.Body).Decode(&loginUnifiedResp); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	data, ok := loginUnifiedResp.Data.(map[string]interface{})
	if !ok {
		t.Error("Expected data field to be an object")
		return
	}

	accessToken, ok := data["access_token"].(string)
	if !ok {
		// Try accessToken (camelCase) if access_token (snake_case) not found
		accessToken, ok = data["accessToken"].(string)
		if !ok {
			t.Error("Expected access_token or accessToken in login response")
			return
		}
	}

	// Now change password using Connect RPC
	changePassReq := map[string]interface{}{
		"old_password": oldPassword,
		"new_password": newPassword,
	}
	changePassBody, _ := json.Marshal(changePassReq)

	req, _ := http.NewRequest("POST", ConnectrpcBaseURL+"/goserver.api.v1.UserService/ChangePassword", bytes.NewReader(changePassBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call Connect RPC ChangePassword: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	t.Logf("Connect RPC ChangePassword response status: %d", resp.StatusCode)
}

// TestConnectRPCGetInstanceProfile tests getting instance profile using Connect RPC endpoint
func TestConnectRPCGetInstanceProfile(t *testing.T) {
	body := []byte(`{}`)
	req, _ := http.NewRequest("POST", ConnectrpcBaseURL+"/goserver.api.v1.InstanceService/GetInstanceProfile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call Connect RPC GetInstanceProfile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	t.Logf("Connect RPC GetInstanceProfile response status: %d", resp.StatusCode)
}
