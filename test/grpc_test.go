package test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	v1pb "github.com/pixb/go-server/proto/gen/api/v1"
)

const grpcAddr = "localhost:8081"

func TestGRPCHealthEndpoint(t *testing.T) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := v1pb.NewUserServiceClient(conn)
	_, err = client.GetUserProfile(context.Background(), &v1pb.GetUserProfileRequest{})
	if err == nil {
		t.Error("Expected error for unauthenticated request, got nil")
	}
	t.Logf("Unauthenticated error (expected): %v", err)
}

func TestGRPCRegisterUser(t *testing.T) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := v1pb.NewUserServiceClient(conn)

	uniqueUsername := fmt.Sprintf("grpcuser_%d", time.Now().UnixNano())
	resp, err := client.RegisterUser(context.Background(), &v1pb.RegisterUserRequest{
		Username: uniqueUsername,
		Nickname: "gRPC Test User",
		Password: "grpcpassword",
		Phone:    "13800138000",
		Email:    uniqueUsername + "@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	if resp.User == nil {
		t.Error("Expected user in response")
	} else if resp.User.Username == "" {
		t.Error("Expected username in response user")
	}
	t.Logf("GRPC RegisterUser response: user=%+v", resp.User)
}

func TestGRPCLoginUser(t *testing.T) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := v1pb.NewUserServiceClient(conn)

	// First register a user
	uniqueUsername := fmt.Sprintf("grpcuser_%d", time.Now().UnixNano())
	registerResp, err := client.RegisterUser(context.Background(), &v1pb.RegisterUserRequest{
		Username: uniqueUsername,
		Nickname: "gRPC Login Test User",
		Password: "testpassword123",
		Phone:    "13800138001",
		Email:    uniqueUsername + "@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}
	t.Logf("GRPC RegisterUser response: %+v", registerResp)

	// Then login with the registered user
	authClient := v1pb.NewAuthServiceClient(conn)
	resp, err := authClient.Login(context.Background(), &v1pb.LoginRequest{
		Username: uniqueUsername,
		Password: "testpassword123",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	if resp.AccessToken == "" {
		t.Error("Expected accessToken in response")
	}
	t.Logf("GRPC LoginUser response: accessToken=%s, user=%s", resp.AccessToken[:30]+"...", resp.User.Username)
}

func TestGRPCGetUserProfileWithAuth(t *testing.T) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	authClient := v1pb.NewUserServiceClient(conn)

	// First register a user
	uniqueUsername := fmt.Sprintf("profileuser_%d", time.Now().UnixNano())
	_, err = authClient.RegisterUser(context.Background(), &v1pb.RegisterUserRequest{
		Username: uniqueUsername,
		Nickname: "gRPC Profile Test User",
		Password: "testpassword123",
		Phone:    "13800138002",
		Email:    uniqueUsername + "@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Then login with the registered user
	loginClient := v1pb.NewAuthServiceClient(conn)
	loginResp, err := loginClient.Login(context.Background(), &v1pb.LoginRequest{
		Username: uniqueUsername,
		Password: "testpassword123",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+loginResp.AccessToken))

	client := v1pb.NewUserServiceClient(conn)
	resp, err := client.GetUserProfile(ctx, &v1pb.GetUserProfileRequest{})
	if err != nil {
		t.Fatalf("Failed to get user profile: %v", err)
	}

	t.Logf("GRPC GetUserProfile response type: %T", resp)
	t.Logf("GRPC GetUserProfile response: %+v", resp)
	// Use reflection to inspect the response object
	v := reflect.ValueOf(resp).Elem()
	t.Logf("Number of fields: %d", v.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)
		t.Logf("Field %d: %s = %v", i, field.Name, value)
	}
	// Check if Username field exists directly or through User field
	if v.FieldByName("Username").IsValid() {
		if v.FieldByName("Username").String() == "" {
			t.Error("Expected username in response")
		}
	} else if v.FieldByName("User").IsValid() {
		user := v.FieldByName("User").Elem()
		if user.FieldByName("Username").String() == "" {
			t.Error("Expected username in response")
		}
	} else {
		t.Error("Expected Username or User field in response")
	}
}

func TestGRPCUpdateUserProfile(t *testing.T) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	authClient := v1pb.NewUserServiceClient(conn)

	// First register a user
	uniqueUsername := fmt.Sprintf("updateuser_%d", time.Now().UnixNano())
	_, err = authClient.RegisterUser(context.Background(), &v1pb.RegisterUserRequest{
		Username: uniqueUsername,
		Nickname: "Original Nickname",
		Password: "testpassword123",
		Phone:    "13800138000",
		Email:    uniqueUsername + "@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Then login with the registered user
	loginClient := v1pb.NewAuthServiceClient(conn)
	loginResp, err := loginClient.Login(context.Background(), &v1pb.LoginRequest{
		Username: uniqueUsername,
		Password: "testpassword123",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+loginResp.AccessToken))

	client := v1pb.NewUserServiceClient(conn)

	// Update user profile
	updateResp, err := client.UpdateUserProfile(ctx, &v1pb.UpdateUserProfileRequest{
		Nickname: "Updated Nickname",
		Phone:    "13900139000",
		Email:    "updated-" + uniqueUsername + "@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to update user profile: %v", err)
	}

	// Verify the update
	t.Logf("GRPC UpdateUserProfile response type: %T", updateResp)
	t.Logf("GRPC UpdateUserProfile response: %+v", updateResp)
	// Use reflection to inspect the response object
	v := reflect.ValueOf(updateResp).Elem()
	t.Logf("Number of fields: %d", v.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)
		t.Logf("Field %d: %s = %v", i, field.Name, value)
	}
	// Check if fields exist directly or through User field
	nicknameFound := false
	phoneFound := false
	emailFound := false
	if v.FieldByName("Nickname").IsValid() {
		if v.FieldByName("Nickname").String() == "Updated Nickname" {
			nicknameFound = true
		} else {
			t.Errorf("Expected nickname to be 'Updated Nickname', got %s", v.FieldByName("Nickname").String())
		}
		if v.FieldByName("Phone").String() == "13900139000" {
			phoneFound = true
		} else {
			t.Errorf("Expected phone to be '13900139000', got %s", v.FieldByName("Phone").String())
		}
		if v.FieldByName("Email").String() == "updated-"+uniqueUsername+"@example.com" {
			emailFound = true
		} else {
			t.Errorf("Expected email to be 'updated-%s@example.com', got %s", uniqueUsername, v.FieldByName("Email").String())
		}
	} else if v.FieldByName("User").IsValid() {
		user := v.FieldByName("User").Elem()
		if user.FieldByName("Nickname").String() == "Updated Nickname" {
			nicknameFound = true
		} else {
			t.Errorf("Expected nickname to be 'Updated Nickname', got %s", user.FieldByName("Nickname").String())
		}
		if user.FieldByName("Phone").String() == "13900139000" {
			phoneFound = true
		} else {
			t.Errorf("Expected phone to be '13900139000', got %s", user.FieldByName("Phone").String())
		}
		if user.FieldByName("Email").String() == "updated-"+uniqueUsername+"@example.com" {
			emailFound = true
		} else {
			t.Errorf("Expected email to be 'updated-%s@example.com', got %s", uniqueUsername, user.FieldByName("Email").String())
		}
	} else {
		t.Error("Expected fields or User field in response")
	}
	if !nicknameFound {
		t.Error("Expected nickname in response")
	}
	if !phoneFound {
		t.Error("Expected phone in response")
	}
	if !emailFound {
		t.Error("Expected email in response")
	}
}

func TestGRPCChangePassword(t *testing.T) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	authClient := v1pb.NewUserServiceClient(conn)

	// First register a user
	uniqueUsername := fmt.Sprintf("changepassuser_%d", time.Now().UnixNano())
	_, err = authClient.RegisterUser(context.Background(), &v1pb.RegisterUserRequest{
		Username: uniqueUsername,
		Nickname: "Change Password Test User",
		Password: "oldpassword123",
		Phone:    "13800138000",
		Email:    uniqueUsername + "@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Then login with the registered user
	loginClient := v1pb.NewAuthServiceClient(conn)
	loginResp, err := loginClient.Login(context.Background(), &v1pb.LoginRequest{
		Username: uniqueUsername,
		Password: "oldpassword123",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+loginResp.AccessToken))

	client := v1pb.NewUserServiceClient(conn)

	// Change password
	changePassResp, err := client.ChangePassword(ctx, &v1pb.ChangePasswordRequest{
		OldPassword: "oldpassword123",
		NewPassword: "newpassword123",
	})
	if err != nil {
		t.Fatalf("Failed to change password: %v", err)
	}

	// Verify the password change by logging in with the new password
	newLoginResp, err := loginClient.Login(context.Background(), &v1pb.LoginRequest{
		Username: uniqueUsername,
		Password: "newpassword123",
	})
	if err != nil {
		t.Fatalf("Failed to login with new password: %v", err)
	}

	if newLoginResp.AccessToken == "" {
		t.Error("Expected access token in response")
	}

	t.Logf("GRPC ChangePassword response: %+v", changePassResp)
	t.Logf("GRPC Login with new password successful: accessToken=%s", newLoginResp.AccessToken[:30]+"...")
}

// TestGRPCRefreshToken tests refreshing access token using gRPC
func TestGRPCRefreshToken(t *testing.T) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// First register a user
	userClient := v1pb.NewUserServiceClient(conn)
	uniqueUsername := fmt.Sprintf("refreshtokenuser_%d", time.Now().UnixNano())
	_, err = userClient.RegisterUser(context.Background(), &v1pb.RegisterUserRequest{
		Username: uniqueUsername,
		Nickname: "Refresh Token Test User",
		Password: "testpassword123",
		Phone:    "13800138000",
		Email:    uniqueUsername + "@example.com",
	})
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	// Then login with the registered user to get tokens
	authClient := v1pb.NewAuthServiceClient(conn)
	loginResp, err := authClient.Login(context.Background(), &v1pb.LoginRequest{
		Username: uniqueUsername,
		Password: "testpassword123",
	})
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	// Test refresh token
	refreshResp, err := authClient.RefreshToken(context.Background(), &v1pb.RefreshTokenRequest{
		RefreshToken: loginResp.RefreshToken,
	})
	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}

	if refreshResp.AccessToken == "" {
		t.Error("Expected access token in refresh response")
	}

	t.Logf("GRPC RefreshToken response: accessToken=%s", refreshResp.AccessToken[:30]+"...")
}

// TestGRPCGetInstanceProfile tests getting instance profile using gRPC
func TestGRPCGetInstanceProfile(t *testing.T) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := v1pb.NewInstanceServiceClient(conn)
	resp, err := client.GetInstanceProfile(context.Background(), &v1pb.GetInstanceProfileRequest{})
	if err != nil {
		t.Fatalf("Failed to get instance profile: %v", err)
	}

	// Verify response
	if resp.Version == "" {
		t.Error("Expected version in instance profile")
	}

	// Check if demo field exists (should be present)
	// Note: demo is a boolean, so we just verify the response is valid

	t.Logf("GRPC GetInstanceProfile response: version=%s, demo=%v", resp.Version, resp.Demo)

	// Check if admin field exists
	if resp.Admin != nil {
		t.Logf("GRPC GetInstanceProfile admin: id=%d, username=%s", resp.Admin.Id, resp.Admin.Username)
	} else {
		t.Log("GRPC GetInstanceProfile: no admin user found (instance needs setup)")
	}
}
