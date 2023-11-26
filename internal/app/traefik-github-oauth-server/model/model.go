package model

import (
	"errors"
	"net/http"
)

type RequestGenerateOAuthPageURL struct {
	RedirectURI string `json:"redirect_uri" binding:"required"`
	AuthURL     string `json:"auth_url" binding:"required"`
}

func (m *RequestGenerateOAuthPageURL) Bind(r *http.Request) error {
	if (m.RedirectURI == "") || (m.AuthURL == "") {
		return errors.New("invalid request")
	}
	return nil
}

type ResponseGenerateOAuthPageURL struct {
	OAuthPageURL string `json:"oauth_page_url"`
}

type RequestRedirect struct {
	RID  string `form:"rid" url:"rid" binding:"required"`
	Code string `form:"code" url:"code" binding:"required"`
}

type RequestGetAuthResult struct {
	RID string `form:"rid" url:"rid" binding:"required"`
}

type ResponseGetAuthResult struct {
	RedirectURI     string `json:"redirect_uri"`
	GitHubUserID    string `json:"github_user_id"`
	GitHubUserLogin string `json:"github_user_login"`
}

type ResponseError struct {
	Message string `json:"msg"`
}

type AuthRequest struct {
	RedirectURI     string `json:"redirect_uri"`
	AuthURL         string `json:"auth_url"`
	GitHubUserID    string `json:"github_user_id"`
	GitHubUserLogin string `json:"github_user_login"`
}
