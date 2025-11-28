package main

import (
	"UAS_BACKEND/domain/model"
	"UAS_BACKEND/domain/service"
	"fmt"
	"log"
)

func main() {
	// 1. Inisialisasi Service
	authService := service.AuthService

	// 2. Data dari User (Flow 1)
	req := model.LoginRequest{
		Identifier: "admin@example.com",
		Password:   "rahasia123",
	}

	// 3. Eksekusi Login
	resp, err := authService.Login(req)

	// Error Handling
	if err != nil {
		log.Fatalf("Login Gagal: %v", err)
	}

	// Output Sukses
	fmt.Println("=== Login Berhasil ===")
	fmt.Printf("User  : %s (%s)\n", resp.User.FullName, resp.User.Email)
	fmt.Printf("Token : %s\n", resp.Token)

	// PasswordHash tidak akan muncul di JSON output karena tag `json:"-"`
	// fmt.Printf("Data: %+v\n", resp.User)
}