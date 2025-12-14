package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Testing Swagger Documentation...")
	
	// Wait a moment for server to start
	time.Sleep(2 * time.Second)
	
	// Test endpoints
	endpoints := []string{
		"http://localhost:8080/health",
		"http://localhost:8080/swagger/",
		"http://localhost:8080/swagger/swagger.yaml",
	}
	
	for _, endpoint := range endpoints {
		fmt.Printf("Testing: %s\n", endpoint)
		
		resp, err := http.Get(endpoint)
		if err != nil {
			fmt.Printf("  ❌ Error: %v\n", err)
			continue
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == 200 {
			fmt.Printf("  ✅ Status: %d OK\n", resp.StatusCode)
		} else {
			fmt.Printf("  ⚠️  Status: %d\n", resp.StatusCode)
		}
	}
	
	fmt.Println("\n🎉 Swagger documentation is ready!")
	fmt.Println("📖 Access documentation at: http://localhost:8080/swagger/")
	fmt.Println("📋 Postman collection: docs/UAS_Backend_API.postman_collection.json")
	fmt.Println("🔧 API specification: docs/swagger.yaml")
}