package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Running Unit Tests...")
	
	// Test individual service packages
	testPackages := []string{
		"./tests/service",
	}
	
	for _, pkg := range testPackages {
		fmt.Printf("Testing package: %s\n", pkg)
		cmd := exec.Command("go", "test", "-v", pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Tests failed for package %s: %v\n", pkg, err)
		} else {
			fmt.Printf("Tests passed for package %s\n", pkg)
		}
		fmt.Println("---")
	}
}