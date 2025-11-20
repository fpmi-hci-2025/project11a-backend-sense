package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"sense-backend/tests/manual/client"
	"sense-backend/tests/manual/testdata"
	"sense-backend/tests/manual/testers"
)

func main() {
	var (
		baseURL     = flag.String("url", "http://localhost:8080", "API base URL")
		dataFile    = flag.String("data", "test_data.json", "Test data file path")
		skipAuth    = flag.Bool("skip-auth", false, "Skip auth tests")
		skipPub     = flag.Bool("skip-publications", false, "Skip publication tests")
		skipComment = flag.Bool("skip-comments", false, "Skip comment tests")
		skipFeed    = flag.Bool("skip-feed", false, "Skip feed tests")
		skipProfile = flag.Bool("skip-profile", false, "Skip profile tests")
		skipMedia   = flag.Bool("skip-media", false, "Skip media tests")
		skipAI      = flag.Bool("skip-ai", false, "Skip AI tests")
	)
	flag.Parse()

	// Load test data
	if err := testdata.LoadTestData(*dataFile); err != nil {
		log.Printf("Warning: Failed to load test data: %v (will create new)", err)
		// Initialize with default values
		testdata.InitializeTestData(*baseURL)
	}

	// Create API client
	apiClient := client.NewClient(*baseURL)

	// Track errors
	var errors []error

	fmt.Println("========================================")
	fmt.Println("Sense API Manual Testing Suite")
	fmt.Println("========================================")
	fmt.Printf("Base URL: %s\n", *baseURL)
	fmt.Printf("Data file: %s\n", *dataFile)
	fmt.Println()

	// Run tests
	if !*skipAuth {
		if err := testers.TestAuthEndpoints(apiClient); err != nil {
			log.Printf("Auth tests failed: %v", err)
			errors = append(errors, fmt.Errorf("auth: %w", err))
		}
		// Save tokens after auth tests
		if err := testdata.SaveTestData(*dataFile); err != nil {
			log.Printf("Warning: Failed to save test data: %v", err)
		}
	}

	if !*skipPub {
		if err := testers.TestPublicationEndpoints(apiClient); err != nil {
			log.Printf("Publication tests failed: %v", err)
			errors = append(errors, fmt.Errorf("publications: %w", err))
		}
		if err := testdata.SaveTestData(*dataFile); err != nil {
			log.Printf("Warning: Failed to save test data: %v", err)
		}
	}

	if !*skipComment {
		if err := testers.TestCommentEndpoints(apiClient); err != nil {
			log.Printf("Comment tests failed: %v", err)
			errors = append(errors, fmt.Errorf("comments: %w", err))
		}
		if err := testdata.SaveTestData(*dataFile); err != nil {
			log.Printf("Warning: Failed to save test data: %v", err)
		}
	}

	if !*skipFeed {
		if err := testers.TestFeedEndpoints(apiClient); err != nil {
			log.Printf("Feed tests failed: %v", err)
			errors = append(errors, fmt.Errorf("feed: %w", err))
		}
	}

	if !*skipProfile {
		if err := testers.TestProfileEndpoints(apiClient); err != nil {
			log.Printf("Profile tests failed: %v", err)
			errors = append(errors, fmt.Errorf("profile: %w", err))
		}
	}

	if !*skipMedia {
		if err := testers.TestMediaEndpoints(apiClient); err != nil {
			log.Printf("Media tests failed: %v", err)
			errors = append(errors, fmt.Errorf("media: %w", err))
		}
		if err := testdata.SaveTestData(*dataFile); err != nil {
			log.Printf("Warning: Failed to save test data: %v", err)
		}
	}

	if !*skipAI {
		if err := testers.TestAIEndpoints(apiClient); err != nil {
			log.Printf("AI tests failed: %v", err)
			errors = append(errors, fmt.Errorf("ai: %w", err))
		}
		if err := testdata.SaveTestData(*dataFile); err != nil {
			log.Printf("Warning: Failed to save test data: %v", err)
		}
	}

	// Final save
	if err := testdata.SaveTestData(*dataFile); err != nil {
		log.Printf("Warning: Failed to save test data: %v", err)
	}

	// Summary
	fmt.Println("\n========================================")
	fmt.Println("Testing Summary")
	fmt.Println("========================================")
	if len(errors) == 0 {
		fmt.Println("✓ All tests passed!")
		os.Exit(0)
	} else {
		fmt.Printf("✗ %d test suite(s) failed:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  - %v\n", err)
		}
		os.Exit(1)
	}
}
