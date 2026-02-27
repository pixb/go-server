package common

// IsPublicMethod checks if a procedure is public (no authentication required)
func IsPublicMethod(procedure string) bool {
	publicMethods := map[string]bool{
		"/goserver.api.v1.AuthService/Login":         true,
		"/goserver.api.v1.AuthService/RefreshToken":  true,
		"/goserver.api.v1.AuthService/ValidateToken": true,
		"/goserver.api.v1.UserService/RegisterUser":  true,
	}

	return publicMethods[procedure]
}
