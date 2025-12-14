package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APITester struct {
	BaseURL string
	Token   string
}

func NewAPITester(baseURL string) *APITester {
	return &APITester{
		BaseURL: baseURL,
	}
}

func main() {
	fmt.Println("🧪 Testing API v1 Endpoints...")
	
	tester := NewAPITester("http://localhost:8080")
	
	// Test authentication
	if err := tester.TestAuthentication(); err != nil {
		fmt.Printf("❌ Authentication test failed: %v\n", err)
		return
	}
	
	// Test users endpoints
	if err := tester.TestUsers(); err != nil {
		fmt.Printf("❌ Users test failed: %v\n", err)
		return
	}
	
	// Test achievements endpoints
	if err := tester.TestAchievements(); err != nil {
		fmt.Printf("❌ Achievements test failed: %v\n", err)
		return
	}
	
	// Test students endpoints
	if err := tester.TestStudents(); err != nil {
		fmt.Printf("❌ Students test failed: %v\n", err)
		return
	}
	
	// Test reports endpoints
	if err := tester.TestReports(); err != nil {
		fmt.Printf("❌ Reports test failed: %v\n", err)
		return
	}
	
	fmt.Println("✅ All API v1 tests passed!")
}

func (t *APITester) TestAuthentication() error {
	fmt.Println("🔐 Testing Authentication endpoints...")
	
	// Test login
	loginData := map[string]string{
		"identifier": "admin@example.com",
		"password":   "password123",
	}
	
	resp, err := t.makeRequest("POST", "/api/v1/auth/login", loginData, false)
	if err != nil {
		return fmt.Errorf("login request failed: %v", err)
	}
	
	var loginResp map[string]interface{}
	if err := json.Unmarshal(resp, &loginResp); err != nil {
		return fmt.Errorf("failed to parse login response: %v", err)
	}
	
	if !loginResp["success"].(bool) {
		return fmt.Errorf("login failed: %v", loginResp["error"])
	}
	
	// Extract token
	data := loginResp["data"].(map[string]interface{})
	t.Token = data["token"].(string)
	fmt.Printf("  ✅ Login successful, token: %s...\n", t.Token[:20])
	
	// Test profile
	resp, err = t.makeRequest("GET", "/api/v1/auth/profile", nil, true)
	if err != nil {
		return fmt.Errorf("profile request failed: %v", err)
	}
	
	var profileResp map[string]interface{}
	if err := json.Unmarshal(resp, &profileResp); err != nil {
		return fmt.Errorf("failed to parse profile response: %v", err)
	}
	
	if !profileResp["success"].(bool) {
		return fmt.Errorf("profile request failed: %v", profileResp["error"])
	}
	
	fmt.Println("  ✅ Profile retrieved successfully")
	return nil
}

func (t *APITester) TestUsers() error {
	fmt.Println("👥 Testing Users endpoints...")
	
	// Test get all users
	resp, err := t.makeRequest("GET", "/api/v1/users?page=1&limit=5", nil, true)
	if err != nil {
		return fmt.Errorf("get users request failed: %v", err)
	}
	
	var usersResp map[string]interface{}
	if err := json.Unmarshal(resp, &usersResp); err != nil {
		return fmt.Errorf("failed to parse users response: %v", err)
	}
	
	if !usersResp["success"].(bool) {
		return fmt.Errorf("get users failed: %v", usersResp["error"])
	}
	
	fmt.Println("  ✅ Users list retrieved successfully")
	return nil
}

func (t *APITester) TestAchievements() error {
	fmt.Println("🏆 Testing Achievements endpoints...")
	
	// Test get achievements
	resp, err := t.makeRequest("GET", "/api/v1/achievements?page=1&limit=5", nil, true)
	if err != nil {
		return fmt.Errorf("get achievements request failed: %v", err)
	}
	
	var achievementsResp map[string]interface{}
	if err := json.Unmarshal(resp, &achievementsResp); err != nil {
		return fmt.Errorf("failed to parse achievements response: %v", err)
	}
	
	if !achievementsResp["success"].(bool) {
		return fmt.Errorf("get achievements failed: %v", achievementsResp["error"])
	}
	
	fmt.Println("  ✅ Achievements list retrieved successfully")
	return nil
}

func (t *APITester) TestStudents() error {
	fmt.Println("🎓 Testing Students endpoints...")
	
	// Test get students
	resp, err := t.makeRequest("GET", "/api/v1/students?page=1&limit=5", nil, true)
	if err != nil {
		return fmt.Errorf("get students request failed: %v", err)
	}
	
	var studentsResp map[string]interface{}
	if err := json.Unmarshal(resp, &studentsResp); err != nil {
		return fmt.Errorf("failed to parse students response: %v", err)
	}
	
	if !studentsResp["success"].(bool) {
		return fmt.Errorf("get students failed: %v", studentsResp["error"])
	}
	
	fmt.Println("  ✅ Students list retrieved successfully")
	return nil
}

func (t *APITester) TestReports() error {
	fmt.Println("📊 Testing Reports endpoints...")
	
	// Test get statistics
	resp, err := t.makeRequest("GET", "/api/v1/reports/statistics?type=overview", nil, true)
	if err != nil {
		return fmt.Errorf("get statistics request failed: %v", err)
	}
	
	var statsResp map[string]interface{}
	if err := json.Unmarshal(resp, &statsResp); err != nil {
		return fmt.Errorf("failed to parse statistics response: %v", err)
	}
	
	if !statsResp["success"].(bool) {
		return fmt.Errorf("get statistics failed: %v", statsResp["error"])
	}
	
	fmt.Println("  ✅ Statistics retrieved successfully")
	return nil
}

func (t *APITester) makeRequest(method, endpoint string, data interface{}, useAuth bool) ([]byte, error) {
	var body io.Reader
	
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}
	
	req, err := http.NewRequest(method, t.BaseURL+endpoint, body)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	if useAuth && t.Token != "" {
		req.Header.Set("Authorization", "Bearer "+t.Token)
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	return io.ReadAll(resp.Body)
}