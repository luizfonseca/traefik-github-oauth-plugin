package model

type RequestGenerateOAuthPageURL struct {
	RedirectURI string `json:"redirect_uri"`
	AuthURL     string `json:"auth_url"`
}

type RequestRedirect struct {
	RID  string `in:"query=rid;form=rid;required" json:"rid"`
	Code string `in:"query=code;form=code;required" json:"code"`
}

type RequestGetAuthResult struct {
	RID string `in:"query=rid;form=rid;required" json:"rid"`
}

type ResponseGenerateOAuthPageURL struct {
	OAuthPageURL string `json:"oauth_page_url"`
}

type ResponseGetAuthResult struct {
	RedirectURI     string   `json:"redirect_uri"`
	GitHubUserID    string   `json:"github_user_id"`
	GitHubUserLogin string   `json:"github_user_login"`
	GithubTeamIDs   []string `json:"github_team_ids"`
}

type ResponseError struct {
	Message string `json:"msg"`
}

type AuthRequest struct {
	RedirectURI     string   `json:"redirect_uri"`
	AuthURL         string   `json:"auth_url"`
	GitHubUserID    string   `json:"github_user_id"`
	GitHubUserLogin string   `json:"github_user_login"`
	GithubTeamIDs   []string `json:"github_team_ids"`
}
