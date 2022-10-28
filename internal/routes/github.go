package routes

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"

	"github.com/eleven-sh/api/internal/envvars"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

const (
	GithubOAuthAPIToCLIURLScheme = "http"
	GithubOAuthAPIToCLIURLHost   = "127.0.0.1"
	GitHubOAuthAPIToCLIURLPath   = "/github/oauth/callback"
)

// Reference: https://docs.github.com/en/developers/apps/managing-oauth-apps/troubleshooting-authorization-request-errors
func getDescriptionForGitHubErrorCode(errorCode string) string {

	if errorCode == "application_suspended" {
		return "The Eleven application has been suspended. " +
			"Please open a new issue at: https://github.com/eleven-sh/cli/issues/new"
	}

	if errorCode == "redirect_uri_mismatch" {
		return "The Eleven application has been misconfigured. " +
			"Please open a new issue at: https://github.com/eleven-sh/cli/issues/new"
	}

	if errorCode == "access_denied" {
		return "You have chosen to not authorize the Eleven application."
	}

	return fmt.Sprintf(
		"An unknown error has occured (\"%s\"). "+
			"Please open a new issue at: https://github.com/eleven-sh/cli/issues/new",
		errorCode,
	)
}

func GitHubOAuthCallback(c *gin.Context) {
	// Listen port is passed through OAuth
	// state because GitHub doesn't support
	// dynamic redirect URIs
	cliListenPort := c.Query("state")
	onlyNumbersRegex := regexp.MustCompile(`^[0-9]+$`)

	if !onlyNumbersRegex.MatchString(cliListenPort) {
		c.String(
			http.StatusBadRequest,
			"Bad state. Please retry the GitHub authorization process.",
		)
		return
	}

	githubOAuthAPIToCLIURLObj := url.URL{
		Scheme: GithubOAuthAPIToCLIURLScheme,
		Host: net.JoinHostPort(
			GithubOAuthAPIToCLIURLHost,
			cliListenPort,
		),
		Path: GitHubOAuthAPIToCLIURLPath,
	}
	githubOAuthAPIToCLIURL := githubOAuthAPIToCLIURLObj.String()

	errorCodeInQuery := c.Query("error")

	if len(errorCodeInQuery) > 0 {
		errorDescription := getDescriptionForGitHubErrorCode(errorCodeInQuery)

		c.Redirect(
			http.StatusTemporaryRedirect,
			githubOAuthAPIToCLIURL+"?error="+url.QueryEscape(errorDescription),
		)

		return
	}

	oauthCodeInQuery := c.Query("code")

	if len(oauthCodeInQuery) == 0 {
		c.String(
			http.StatusBadRequest,
			"Missing OAuth code. Please retry the GitHub authorization process.",
		)
		return
	}

	githubOAuthClient := &oauth2.Config{
		ClientID:     envvars.Get(envvars.EnvVarNameGitHubOAuthClientID),
		ClientSecret: envvars.Get(envvars.EnvVarNameGitHubOAuthClientSecret),
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	githubOAuthTokens, err := githubOAuthClient.Exchange(
		context.TODO(),
		oauthCodeInQuery,
	)

	if err != nil {
		c.Redirect(
			http.StatusTemporaryRedirect,
			githubOAuthAPIToCLIURL+"?error="+url.QueryEscape(err.Error()),
		)
		return
	}

	c.Redirect(
		http.StatusTemporaryRedirect,
		githubOAuthAPIToCLIURL+"?access_token="+url.QueryEscape(githubOAuthTokens.AccessToken),
	)
}
