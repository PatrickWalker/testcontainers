package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type appContainer struct {
	URI string
	err tc.ExecError
}

var compose *tc.LocalDockerCompose

func setupAppUnderTest(ctx context.Context) (*appContainer, error) {

	composeFilePaths := []string{"docker-compose.yml"}
	identifier := strings.ToLower(uuid.New().String())

	compose = tc.NewLocalDockerCompose(composeFilePaths, identifier)
	execError := compose.
		WithCommand([]string{"up", "--build", "-d", "--force-recreate"}).
		WithExposedService("users-api_1", 8082, wait.NewHTTPStrategy("/").WithPort("8080/tcp").WithStartupTimeout(30*time.Second)).
		Invoke()
	err := execError.Error

	if err != nil {
		return nil, fmt.Errorf("Could not run compose file: %v - %v", composeFilePaths, err)
	}

	uri := fmt.Sprintf("http://localhost:8082")

	return &appContainer{URI: uri}, nil
}

func teardownAppUnderTest(ctx context.Context) {
	execError := compose.Down()
	err := execError.Error
	if err != nil {
		fmt.Printf("Could not stop compose file: %v \n", err)
	}
}

func TestIntegrationAppLatestReturn(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	//if you dont do this the containers hang around and you cant run again
	defer teardownAppUnderTest(ctx)

	appC, err := setupAppUnderTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	fmt.Printf("appC.URI %v\n", appC.URI)
	resp, err := http.Get(appC.URI)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
	fmt.Println("Hello world passed")

	//check redis endpoint
	resp, err = http.Get(fmt.Sprintf("%v/redis", appC.URI))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected redis status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
	fmt.Println("Redis passed")

	//check microservice dep endpoint
	resp, err = http.Get(fmt.Sprintf("%v/dep", appC.URI))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected dep status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
	fmt.Println("DEP passed")

	//check db endpoint
	resp, err = http.Get(fmt.Sprintf("%v/db", appC.URI))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected db status code %d. Got %d.", http.StatusOK, resp.StatusCode)
	}
	fmt.Println("DB passed")

}
