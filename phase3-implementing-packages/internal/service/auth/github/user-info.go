package github

import "strconv"

type GitHubUserInfo struct{
	UID int  `json:"id"`
	Extra interface{}
}

func (p *GitHubUserInfo) GetProviderUID() string {
	return strconv.Itoa(p.UID)
}
