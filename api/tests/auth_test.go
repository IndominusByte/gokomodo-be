package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/IndominusByte/gokomodo-be/api/internal/config"
	"github.com/creent-production/cdk-go/auth"
	"github.com/go-chi/jwtauth"
	"github.com/stretchr/testify/assert"
)

const (
	prefix  = "/auth"
	email   = "user@example.com"
	email_2 = "user2@example.com"
)

func TestValidationLogin(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name    string
		payload map[string]string
	}{

		{
			name:    "required",
			payload: map[string]string{"email": "", "password": ""},
		},
		{
			name:    "minimum",
			payload: map[string]string{"email": "a", "password": "a"},
		},
		{
			name:    "maximum",
			payload: map[string]string{"email": createMaximum(200), "password": createMaximum(200)},
		},
		{
			name:    "invalid email",
			payload: map[string]string{"email": "test@asdcom"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, prefix+"/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "required":
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["email"].(string))
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["password"].(string))
			case "minimum":
				assert.Equal(t, "Shorter than minimum length 6.", data["detail_message"].(map[string]interface{})["password"].(string))
			case "maximum":
				assert.Equal(t, "Longer than maximum length 100.", data["detail_message"].(map[string]interface{})["password"].(string))
			case "invalid email":
				assert.Equal(t, "Not a valid email address.", data["detail_message"].(map[string]interface{})["email"].(string))
			}
			assert.Equal(t, 422, response.Result().StatusCode)
		})
	}
}

func TestLogin(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name       string
		expected   string
		payload    map[string]string
		statusCode int
	}{
		{
			name:       "email not found",
			expected:   "Invalid credential.",
			payload:    map[string]string{"email": "test@test.com", "password": "asdasd"},
			statusCode: 422,
		},
		{
			name:       "password wrong",
			expected:   "Invalid credential.",
			payload:    map[string]string{"email": email, "password": "asdasd2"},
			statusCode: 422,
		},
		{
			name:       "success",
			expected:   "",
			payload:    map[string]string{"email": email, "password": "string"},
			statusCode: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, prefix+"/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			if test.name != "success" {
				assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_app"].(string))
				assert.Equal(t, test.statusCode, response.Result().StatusCode)
			} else {
				assert.NotNil(t, data["results"].(map[string]interface{})["access_token"].(string))
				assert.NotNil(t, data["results"].(map[string]interface{})["refresh_token"].(string))
				assert.Equal(t, 200, response.Result().StatusCode)
			}
		})
	}
}

func TestValidationFreshToken(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	// login
	body, err := json.Marshal(map[string]string{"email": email, "password": "string"})
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest(http.MethodPost, prefix+"/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, s)

	body, _ = io.ReadAll(response.Result().Body)
	json.Unmarshal(body, &data)

	accessToken := data["results"].(map[string]interface{})["access_token"].(string)

	tests := [...]struct {
		name    string
		payload map[string]string
	}{
		{
			name:    "required",
			payload: map[string]string{"password": ""},
		},
		{
			name:    "minimum",
			payload: map[string]string{"password": "a"},
		},
		{
			name:    "maximum",
			payload: map[string]string{"password": createMaximum(200)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, prefix+"/fresh-token", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", "Bearer "+accessToken)

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "required":
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["password"].(string))
			case "minimum":
				assert.Equal(t, "Shorter than minimum length 6.", data["detail_message"].(map[string]interface{})["password"].(string))
			case "maximum":
				assert.Equal(t, "Longer than maximum length 100.", data["detail_message"].(map[string]interface{})["password"].(string))
			}
			assert.Equal(t, 422, response.Result().StatusCode)
		})
	}
}

func TestFreshToken(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	// login
	body, err := json.Marshal(map[string]string{"email": email, "password": "string"})
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest(http.MethodPost, prefix+"/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, s)

	body, _ = io.ReadAll(response.Result().Body)
	json.Unmarshal(body, &data)

	accessToken := data["results"].(map[string]interface{})["access_token"].(string)

	cfg, _ := config.New()
	token := auth.GenerateAccessToken(&auth.AccessToken{Sub: strconv.Itoa(0), Exp: jwtauth.ExpireIn(cfg.JWT.AccessExpires), Fresh: true})
	tokenUserNotFound := auth.NewJwtTokenRSA(cfg.JWT.PublicKey, cfg.JWT.PrivateKey, cfg.JWT.Algorithm, token)

	tests := [...]struct {
		name       string
		payload    map[string]string
		token      string
		statusCode int
	}{
		{
			name:       "user not found",
			payload:    map[string]string{"password": "asdasd"},
			token:      tokenUserNotFound,
			statusCode: 401,
		},
		{
			name:       "password not same",
			payload:    map[string]string{"password": "asdasd2"},
			token:      accessToken,
			statusCode: 422,
		},
		{
			name:       "success",
			payload:    map[string]string{"password": "string"},
			token:      accessToken,
			statusCode: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, prefix+"/fresh-token", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Add("Authorization", "Bearer "+test.token)

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "user not found":
				assert.Equal(t, "User not found.", data["detail_message"].(map[string]interface{})["_header"].(string))
			case "password not same":
				assert.Equal(t, "Password does not match with our records.", data["detail_message"].(map[string]interface{})["_app"].(string))
			case "success":
				assert.NotNil(t, data["results"].(map[string]interface{})["access_token"].(string))
			}
			assert.Equal(t, test.statusCode, response.Result().StatusCode)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	// login
	body, err := json.Marshal(map[string]string{"email": email, "password": "string"})
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest(http.MethodPost, prefix+"/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, s)

	body, _ = io.ReadAll(response.Result().Body)
	json.Unmarshal(body, &data)

	refreshToken := data["results"].(map[string]interface{})["refresh_token"].(string)

	cfg, _ := config.New()
	token := auth.GenerateRefreshToken(&auth.RefreshToken{Sub: strconv.Itoa(0), Exp: jwtauth.ExpireIn(cfg.JWT.RefreshExpires)})
	tokenUserNotFound := auth.NewJwtTokenRSA(cfg.JWT.PublicKey, cfg.JWT.PrivateKey, cfg.JWT.Algorithm, token)

	tests := [...]struct {
		name       string
		token      string
		statusCode int
	}{
		{
			name:       "user not found",
			token:      tokenUserNotFound,
			statusCode: 401,
		},
		{
			name:       "success",
			token:      refreshToken,
			statusCode: 200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, prefix+"/refresh-token", nil)
			req.Header.Add("Authorization", "Bearer "+test.token)

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "user not found":
				assert.Equal(t, "User not found.", data["detail_message"].(map[string]interface{})["_header"].(string))
			case "success":
				assert.NotNil(t, data["results"].(map[string]interface{})["access_token"].(string))
			}
			assert.Equal(t, test.statusCode, response.Result().StatusCode)
		})
	}
}

func TestAccessRevoke(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	// login
	body, err := json.Marshal(map[string]string{"email": email, "password": "string"})
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest(http.MethodPost, prefix+"/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, s)

	body, _ = io.ReadAll(response.Result().Body)
	json.Unmarshal(body, &data)

	accessToken := data["results"].(map[string]interface{})["access_token"].(string)

	tests := [...]struct {
		name       string
		statusCode int
	}{
		{
			name:       "success",
			statusCode: 200,
		},
		{
			name:       "revoked",
			statusCode: 401,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, prefix+"/access-revoke", nil)
			req.Header.Add("Authorization", "Bearer "+accessToken)

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "success":
				assert.Equal(t, "An access token has revoked.", data["detail_message"].(map[string]interface{})["_app"].(string))
			case "revoked":
				assert.Equal(t, "token is revoked", data["detail_message"].(map[string]interface{})["_header"].(string))
			}
			assert.Equal(t, test.statusCode, response.Result().StatusCode)
		})
	}

}

func TestRefreshRevoke(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	// login
	body, err := json.Marshal(map[string]string{"email": email, "password": "string"})
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest(http.MethodPost, prefix+"/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, s)

	body, _ = io.ReadAll(response.Result().Body)
	json.Unmarshal(body, &data)

	refreshToken := data["results"].(map[string]interface{})["refresh_token"].(string)

	tests := [...]struct {
		name       string
		statusCode int
	}{
		{
			name:       "success",
			statusCode: 200,
		},
		{
			name:       "revoked",
			statusCode: 401,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, prefix+"/refresh-revoke", nil)
			req.Header.Add("Authorization", "Bearer "+refreshToken)

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "success":
				assert.Equal(t, "An refresh token has revoked.", data["detail_message"].(map[string]interface{})["_app"].(string))
			case "revoked":
				assert.Equal(t, "token is revoked", data["detail_message"].(map[string]interface{})["_header"].(string))
			}
			assert.Equal(t, test.statusCode, response.Result().StatusCode)
		})
	}

}
