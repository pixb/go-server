package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

// ApiBaseURL 测试服务器基础 URL
var ApiBaseURL = getApiBaseURL()

// getApiBaseURL 从环境变量获取基础 URL，若未设置则使用默认值
func getApiBaseURL() string {
	if url := os.Getenv("TEST_BASE_URL"); url != "" {
		return url
	}
	return "http://localhost:8081"
}

// ApiUnifiedResponse represents the expected unified response format for API tests
type ApiUnifiedResponse struct {
	State   int         `json:"state"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get(ApiBaseURL + "/healthz")
	if err != nil {
		t.Fatalf("Failed to connect to health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestRegisterUser tests the user registration endpoint
func TestRegisterUser(t *testing.T) {
	uniqueUsername := fmt.Sprintf("httpuser_%d", time.Now().UnixNano())
	uniqueEmail := uniqueUsername + "@example.com"

	user := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "HTTP Test User",
		"password": "testpassword123",
		"phone":    "13800138000",
		"email":    uniqueEmail,
	}
	body, _ := json.Marshal(user)

	resp, err := http.Post(ApiBaseURL+"/api/v1/users", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to call RegisterUser: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status 200 or 409, got %d", resp.StatusCode)
	}

	// Verify response format
	var unifiedResp ApiUnifiedResponse
	if err := json.NewDecoder(resp.Body).Decode(&unifiedResp); err != nil {
		t.Fatalf("Failed to decode unified response: %v", err)
	}

	// Check if it's a success or error response
	if resp.StatusCode == http.StatusOK {
		// Success response
		if unifiedResp.State != 0 {
			t.Errorf("Expected state 0 for success, got %d", unifiedResp.State)
		}
		if unifiedResp.Message != "success" {
			t.Errorf("Expected message 'success', got '%s'", unifiedResp.Message)
		}
		if unifiedResp.Data == nil {
			t.Error("Expected data field to be non-nil for success response")
		} else {
			// Verify data structure contains expected fields
			data, ok := unifiedResp.Data.(map[string]interface{})
			if !ok {
				t.Error("Expected data field to be an object")
				return
			}

			// Check if user exists in data
			user, ok := data["user"].(map[string]interface{})
			if !ok {
				t.Error("Expected data field to contain 'user' object")
				return
			}

			// Check if user ID exists
			if _, ok := user["id"]; !ok {
				t.Error("Expected user field to contain 'id'")
			}

			// Check if username matches
			if username, ok := user["username"].(string); !ok || username != uniqueUsername {
				t.Errorf("Expected username '%s', got '%v'", uniqueUsername, user["username"])
			}
		}
	} else {
		// Error response
		if unifiedResp.State == 0 {
			t.Errorf("Expected non-zero state for error, got 0")
		}
		if unifiedResp.Message == "" {
			t.Error("Expected message field to be non-empty for error response")
		}
		if unifiedResp.Data != nil {
			t.Error("Expected data field to be nil for error response")
		}
	}

	t.Logf("RegisterUser response: %+v", unifiedResp)
}

// TestLoginUser tests the user login endpoint
func TestLoginUser(t *testing.T) {
	// First register a user
	uniqueUsername := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	registerUser := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Login Test User",
		"password": "testpassword123",
		"phone":    "13800138001",
		"email":    uniqueUsername + "@example.com",
	}
	registerBody, _ := json.Marshal(registerUser)

	registerResp, err := http.Post(ApiBaseURL+"/api/v1/users", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		t.Fatalf("Failed to call RegisterUser: %v", err)
	}
	registerResp.Body.Close()

	// Then login with the registered user
	loginReq := map[string]string{
		"username": uniqueUsername,
		"password": "testpassword123",
	}
	loginBody, _ := json.Marshal(loginReq)

	resp, err := http.Post(ApiBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}

	// Verify response format
	var unifiedResp ApiUnifiedResponse
	if err := json.NewDecoder(resp.Body).Decode(&unifiedResp); err != nil {
		t.Fatalf("Failed to decode unified response: %v", err)
	}

	// Check if it's a success or error response
	if resp.StatusCode == http.StatusOK {
		// Success response
		if unifiedResp.State != 0 {
			t.Errorf("Expected state 0 for success, got %d", unifiedResp.State)
		}
		if unifiedResp.Message != "success" {
			t.Errorf("Expected message 'success', got '%s'", unifiedResp.Message)
		}
		if unifiedResp.Data == nil {
			t.Error("Expected data field to be non-nil for success response")
			return
		}

		// Check if accessToken exists in the response data
		data, ok := unifiedResp.Data.(map[string]interface{})
		if !ok {
			t.Error("Expected data field to be an object")
			return
		}

		// Check if access_token exists
		if _, ok := data["access_token"]; !ok {
			// Try accessToken (camelCase) if access_token (snake_case) not found
			if _, ok := data["accessToken"]; !ok {
				t.Errorf("Expected access_token or accessToken in response data, got %+v", data)
			}
		}

		// Check if refresh_token exists
		if _, ok := data["refresh_token"]; !ok {
			// Try refreshToken (camelCase) if refresh_token (snake_case) not found
			if _, ok := data["refreshToken"]; !ok {
				t.Errorf("Expected refresh_token or refreshToken in response data, got %+v", data)
			}
		}
	} else {
		// Error response
		if unifiedResp.State == 0 {
			t.Errorf("Expected non-zero state for error, got 0")
		}
		if unifiedResp.Message == "" {
			t.Error("Expected message field to be non-empty for error response")
		}
		if unifiedResp.Data != nil {
			t.Error("Expected data field to be nil for error response")
		}
	}

	t.Logf("LoginUser response: %+v", unifiedResp)
}

// TestGetUserProfileWithoutAuth tests getting user profile without authentication
func TestGetUserProfileWithoutAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", ApiBaseURL+"/api/v1/users/me", nil)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call GetUserProfile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}

	// Verify response format
	var unifiedResp ApiUnifiedResponse
	if err := json.NewDecoder(resp.Body).Decode(&unifiedResp); err != nil {
		t.Fatalf("Failed to decode unified response: %v", err)
	}

	// Check unified response format for error
	if unifiedResp.State == 0 {
		t.Errorf("Expected non-zero state for error, got 0")
	}
	if unifiedResp.Message == "" {
		t.Error("Expected message field to be non-empty for error response")
	}
	if unifiedResp.Data != nil {
		t.Error("Expected data field to be nil for error response")
	}

	t.Logf("GetUserProfileWithoutAuth response: %+v", unifiedResp)
}

// TestRegisterUserWithDuplicateEmail tests registering user with duplicate email
func TestRegisterUserWithDuplicateEmail(t *testing.T) {
	uniqueUsername1 := fmt.Sprintf("duplicateuser1_%d", time.Now().UnixNano())
	uniqueUsername2 := fmt.Sprintf("duplicateuser2_%d", time.Now().UnixNano())
	email := fmt.Sprintf("duplicate_%d@example.com", time.Now().UnixNano())

	// First register a user with the email
	user1 := map[string]interface{}{
		"username": uniqueUsername1,
		"nickname": "First User",
		"password": "testpassword123",
		"phone":    "13800138000",
		"email":    email,
	}
	body1, _ := json.Marshal(user1)

	resp1, err := http.Post(ApiBaseURL+"/api/v1/users", "application/json", bytes.NewReader(body1))
	if err != nil {
		t.Fatalf("Failed to call RegisterUser: %v", err)
	}
	resp1.Body.Close()

	if resp1.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for first registration, got %d", resp1.StatusCode)
		return
	}

	// Then try to register another user with the same email
	user2 := map[string]interface{}{
		"username": uniqueUsername2,
		"nickname": "Second User",
		"password": "testpassword123",
		"phone":    "13800138001",
		"email":    email,
	}
	body2, _ := json.Marshal(user2)

	resp2, err := http.Post(ApiBaseURL+"/api/v1/users", "application/json", bytes.NewReader(body2))
	if err != nil {
		t.Fatalf("Failed to call RegisterUser: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate email, got %d", resp2.StatusCode)
	}

	// Verify response format
	var unifiedResp ApiUnifiedResponse
	if err := json.NewDecoder(resp2.Body).Decode(&unifiedResp); err != nil {
		t.Fatalf("Failed to decode unified response: %v", err)
	}

	// Check error response format
	if unifiedResp.State == 0 {
		t.Errorf("Expected non-zero state for error, got 0")
	}
	if unifiedResp.Message == "" {
		t.Error("Expected message field to be non-empty for error response")
	}
	if unifiedResp.Data != nil {
		t.Error("Expected data field to be nil for error response")
	}

	t.Logf("RegisterUserWithDuplicateEmail response: %+v", unifiedResp)
}

// TestLoginWithInvalidCredentials tests login with invalid credentials
func TestLoginWithInvalidCredentials(t *testing.T) {
	// Try to login with non-existent username
	loginReq := map[string]string{
		"username": "nonexistentuser",
		"password": "wrongpassword",
	}
	loginBody, _ := json.Marshal(loginReq)

	resp, err := http.Post(ApiBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid credentials, got %d", resp.StatusCode)
	}

	// Verify response format
	var unifiedResp ApiUnifiedResponse
	if err := json.NewDecoder(resp.Body).Decode(&unifiedResp); err != nil {
		t.Fatalf("Failed to decode unified response: %v", err)
	}

	// Check error response format
	if unifiedResp.State == 0 {
		t.Errorf("Expected non-zero state for error, got 0")
	}
	if unifiedResp.Message == "" {
		t.Error("Expected message field to be non-empty for error response")
	}
	if unifiedResp.Data != nil {
		t.Error("Expected data field to be nil for error response")
	}

	t.Logf("LoginWithInvalidCredentials response: %+v", unifiedResp)
}

// TestGetUserProfileWithAuth tests getting user profile with authentication
func TestGetUserProfileWithAuth(t *testing.T) {
	// First register and login to get tokens
	uniqueUsername := fmt.Sprintf("authuser_%d", time.Now().UnixNano())
	testPassword := "testpassword123"

	// Register user
	registerUser := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Auth Test User",
		"password": testPassword,
		"phone":    "13800138000",
		"email":    uniqueUsername + "@example.com",
	}
	registerBody, _ := json.Marshal(registerUser)

	registerResp, err := http.Post(ApiBaseURL+"/api/v1/users", "application/json", bytes.NewReader(registerBody))
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

	loginResp, err := http.Post(ApiBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for login, got %d", loginResp.StatusCode)
		return
	}

	// Parse login response to get access token
	var loginUnifiedResp ApiUnifiedResponse
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

	// Now get user profile with authentication
	req, _ := http.NewRequest("GET", ApiBaseURL+"/api/v1/users/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call GetUserProfile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read the response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Try to decode as unified response first
	var unifiedResp ApiUnifiedResponse
	if err := json.Unmarshal(body, &unifiedResp); err == nil {
		// Check unified response format for success
		if unifiedResp.State != 0 {
			t.Errorf("Expected state 0 for success, got %d", unifiedResp.State)
		}

		// Check if user data exists in the response
		if unifiedResp.Data != nil {
			profileData, ok := unifiedResp.Data.(map[string]interface{})
			if ok {
				if _, ok := profileData["user"]; !ok {
					t.Errorf("Expected user in response data, got %+v", profileData)
				}
			}
		}

		t.Logf("GetUserProfileWithAuth response (unified): %+v", unifiedResp)
	} else {
		// Try to decode as raw gRPC response
		var userResponse map[string]interface{}
		if err := json.Unmarshal(body, &userResponse); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Check if user field exists directly
		if _, ok := userResponse["user"]; !ok {
			if _, ok := userResponse["User"]; !ok {
				t.Errorf("Expected user in response data, got %+v", userResponse)
			}
		}

		t.Logf("GetUserProfileWithAuth response (raw): %+v", userResponse)
	}
}

// TestUpdateUserProfile tests updating user profile
func TestUpdateUserProfile(t *testing.T) {
	// First register and login to get tokens
	uniqueUsername := fmt.Sprintf("updateuser_%d", time.Now().UnixNano())
	testPassword := "testpassword123"

	// Register user
	registerUser := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Original Nickname",
		"password": testPassword,
		"phone":    "13800138000",
		"email":    uniqueUsername + "@example.com",
	}
	registerBody, _ := json.Marshal(registerUser)

	registerResp, err := http.Post(ApiBaseURL+"/api/v1/users", "application/json", bytes.NewReader(registerBody))
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

	loginResp, err := http.Post(ApiBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for login, got %d", loginResp.StatusCode)
		return
	}

	// Parse login response to get access token
	var loginUnifiedResp ApiUnifiedResponse
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

	// Now update user profile
	updateUser := map[string]interface{}{
		"nickname": "Updated Nickname",
		"phone":    "13900139000",
		"email":    "updated-" + uniqueUsername + "@example.com",
	}
	updateBody, _ := json.Marshal(updateUser)

	req, _ := http.NewRequest("PATCH", ApiBaseURL+"/api/v1/users/me", bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call UpdateUserProfile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read the response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Try to decode as unified response first
	var unifiedResp ApiUnifiedResponse
	if err := json.Unmarshal(body, &unifiedResp); err == nil {
		// Check unified response format for success
		if unifiedResp.State != 0 {
			t.Errorf("Expected state 0 for success, got %d", unifiedResp.State)
		}

		// Check if user data exists in the response
		if unifiedResp.Data != nil {
			profileData, ok := unifiedResp.Data.(map[string]interface{})
			if ok {
				if _, ok := profileData["user"]; !ok {
					t.Errorf("Expected user in response data, got %+v", profileData)
				}
			}
		}

		t.Logf("UpdateUserProfile response (unified): %+v", unifiedResp)
	} else {
		// Try to decode as raw gRPC response
		var userResponse map[string]interface{}
		if err := json.Unmarshal(body, &userResponse); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Check if user field exists directly
		if _, ok := userResponse["user"]; !ok {
			if _, ok := userResponse["User"]; !ok {
				t.Errorf("Expected user in response data, got %+v", userResponse)
			}
		}

		t.Logf("UpdateUserProfile response (raw): %+v", userResponse)
	}
}

// TestChangePassword tests changing user password
func TestChangePassword(t *testing.T) {
	// First register and login to get tokens
	uniqueUsername := fmt.Sprintf("changepassuser_%d", time.Now().UnixNano())
	oldPassword := "oldpassword123"
	newPassword := "newpassword123"

	// Register user
	registerUser := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Change Password Test User",
		"password": oldPassword,
		"phone":    "13800138000",
		"email":    uniqueUsername + "@example.com",
	}
	registerBody, _ := json.Marshal(registerUser)

	registerResp, err := http.Post(ApiBaseURL+"/api/v1/users", "application/json", bytes.NewReader(registerBody))
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

	loginResp, err := http.Post(ApiBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for login, got %d", loginResp.StatusCode)
		return
	}

	// Parse login response to get access token
	var loginUnifiedResp ApiUnifiedResponse
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

	// Now change password
	changePassReq := map[string]string{
		"old_password": oldPassword,
		"new_password": newPassword,
	}
	changePassBody, _ := json.Marshal(changePassReq)

	req, _ := http.NewRequest("POST", ApiBaseURL+"/api/v1/users/me/password", bytes.NewReader(changePassBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to call ChangePassword: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read the response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Try to decode as unified response first
	var unifiedResp ApiUnifiedResponse
	if err := json.Unmarshal(body, &unifiedResp); err == nil {
		// Check unified response format for success
		if unifiedResp.State != 0 {
			t.Errorf("Expected state 0 for success, got %d", unifiedResp.State)
		}

		t.Logf("ChangePassword response (unified): %+v", unifiedResp)
	} else {
		// Try to decode as raw gRPC response
		var changePassResponse map[string]interface{}
		if err := json.Unmarshal(body, &changePassResponse); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		t.Logf("ChangePassword response (raw): %+v", changePassResponse)
	}

	// Verify the password change by logging in with the new password
	newLoginReq := map[string]string{
		"username": uniqueUsername,
		"password": newPassword,
	}
	newLoginBody, _ := json.Marshal(newLoginReq)

	newLoginResp, err := http.Post(ApiBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(newLoginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser with new password: %v", err)
	}
	defer newLoginResp.Body.Close()

	if newLoginResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for login with new password, got %d", newLoginResp.StatusCode)
	}

	t.Logf("Login with new password successful")
}

// TestRefreshToken tests refreshing access token
func TestRefreshToken(t *testing.T) {
	// First register and login to get tokens
	uniqueUsername := fmt.Sprintf("refreshtokenuser_%d", time.Now().UnixNano())
	testPassword := "testpassword123"

	// Register user
	registerUser := map[string]interface{}{
		"username": uniqueUsername,
		"nickname": "Refresh Token Test User",
		"password": testPassword,
		"phone":    "13800138000",
		"email":    uniqueUsername + "@example.com",
	}
	registerBody, _ := json.Marshal(registerUser)

	registerResp, err := http.Post(ApiBaseURL+"/api/v1/users", "application/json", bytes.NewReader(registerBody))
	if err != nil {
		t.Fatalf("Failed to call RegisterUser: %v", err)
	}
	registerResp.Body.Close()

	// Login to get tokens
	loginReq := map[string]string{
		"username": uniqueUsername,
		"password": testPassword,
	}
	loginBody, _ := json.Marshal(loginReq)

	loginResp, err := http.Post(ApiBaseURL+"/api/v1/auth/login", "application/json", bytes.NewReader(loginBody))
	if err != nil {
		t.Fatalf("Failed to call LoginUser: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for login, got %d", loginResp.StatusCode)
		return
	}

	// Parse login response to get refresh token
	var loginUnifiedResp ApiUnifiedResponse
	if err := json.NewDecoder(loginResp.Body).Decode(&loginUnifiedResp); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	data, ok := loginUnifiedResp.Data.(map[string]interface{})
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

	// Now test refresh token endpoint
	refreshReq := map[string]string{
		"refresh_token": refreshToken,
	}
	refreshBody, _ := json.Marshal(refreshReq)

	resp, err := http.Post(ApiBaseURL+"/api/v1/auth/refresh", "application/json", bytes.NewReader(refreshBody))
	if err != nil {
		t.Fatalf("Failed to call RefreshToken: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read the response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Try to decode as unified response first
	var unifiedResp ApiUnifiedResponse
	if err := json.Unmarshal(body, &unifiedResp); err == nil {
		// Check unified response format for success
		if unifiedResp.State != 0 {
			t.Errorf("Expected state 0 for success, got %d", unifiedResp.State)
		}

		// Check if access_token exists in the response
		if unifiedResp.Data != nil {
			refreshData, ok := unifiedResp.Data.(map[string]interface{})
			if ok {
				if _, ok := refreshData["access_token"]; !ok {
					// Try accessToken (camelCase) if access_token (snake_case) not found
					if _, ok := refreshData["accessToken"]; !ok {
						t.Errorf("Expected access_token or accessToken in refresh response data, got %+v", refreshData)
					}
				}
			}
		}

		t.Logf("RefreshToken response (unified): %+v", unifiedResp)
	} else {
		// Try to decode as raw gRPC response
		var refreshResponse map[string]interface{}
		if err := json.Unmarshal(body, &refreshResponse); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Check if access_token exists directly
		if _, ok := refreshResponse["access_token"]; !ok {
			if _, ok := refreshResponse["AccessToken"]; !ok {
				t.Errorf("Expected access_token in refresh response data, got %+v", refreshResponse)
			}
		}

		t.Logf("RefreshToken response (raw): %+v", refreshResponse)
	}
}

// TestGetInstanceProfile tests getting instance profile
func TestGetInstanceProfile(t *testing.T) {
	resp, err := http.Get(ApiBaseURL + "/api/v1/instance/profile")
	if err != nil {
		t.Fatalf("Failed to call GetInstanceProfile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Read the response body once
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Try to decode as unified response first
	var unifiedResp ApiUnifiedResponse
	if err := json.Unmarshal(body, &unifiedResp); err == nil {
		// Check unified response format for success
		if unifiedResp.State != 0 {
			t.Errorf("Expected state 0 for success, got %d", unifiedResp.State)
		}

		// Check if data exists in the response
		if unifiedResp.Data == nil {
			t.Error("Expected data field to be non-nil for success response")
			return
		}

		// Check instance profile data
		profileData, ok := unifiedResp.Data.(map[string]interface{})
		if !ok {
			t.Error("Expected data field to be an object")
			return
		}

		// Check if version exists
		if _, ok := profileData["version"]; !ok {
			t.Error("Expected version field in instance profile")
		}

		// Check if demo exists
		if _, ok := profileData["demo"]; !ok {
			t.Error("Expected demo field in instance profile")
		}

		t.Logf("GetInstanceProfile response (unified): %+v", unifiedResp)
	} else {
		// Try to decode as raw gRPC response
		var instanceResponse map[string]interface{}
		if err := json.Unmarshal(body, &instanceResponse); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		// Check if version exists directly
		if _, ok := instanceResponse["version"]; !ok {
			if _, ok := instanceResponse["Version"]; !ok {
				t.Error("Expected version field in instance profile response")
			}
		}

		// Check if demo exists directly
		if _, ok := instanceResponse["demo"]; !ok {
			if _, ok := instanceResponse["Demo"]; !ok {
				t.Error("Expected demo field in instance profile response")
			}
		}

		t.Logf("GetInstanceProfile response (raw): %+v", instanceResponse)
	}
}
