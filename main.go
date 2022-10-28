package main

import (
	"log"

	"github.com/eleven-sh/api/internal/envvars"
	"github.com/eleven-sh/api/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	GitHubOAuthCLIToAPIURLPath = "/github/oauth/callback"
)

func main() {
	// godotenv returns an error in prod
	// when the file ".env" is not present...
	_ = godotenv.Load()
	// ...so, we let it fail silently and ensure that
	// all env vars are still set.
	envvars.Ensure(".env.dist")

	gin.SetMode(envvars.Get(envvars.EnvVarNameGinMode))

	r := gin.Default()

	r.GET(
		GitHubOAuthCLIToAPIURLPath,
		routes.GitHubOAuthCallback,
	)

	log.Printf(
		"HTTP server listening at: http://0.0.0.0:%s/",
		envvars.Get(envvars.EnvVarNamePort),
	)

	if err := r.Run(":" + envvars.Get(envvars.EnvVarNamePort)); err != nil {
		log.Fatalf("%v", err)
	}
}
