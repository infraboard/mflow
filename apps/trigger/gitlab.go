package trigger

import "strings"

const INJECT_PROVIDER = "@GitlabEvent"

var InjectParams = map[string]string{
	"GIT_SSH_URL":   INJECT_PROVIDER,
	"GIT_BRANCH":    INJECT_PROVIDER,
	"GIT_COMMIT_ID": INJECT_PROVIDER,
	"APP_VERSION":   INJECT_PROVIDER,
}

func (m *Commit) Short() string {
	if len(m.Id) > 8 {
		return m.Id[:8]
	}

	return m.Id
}

// "GitLab/15.5.0-pre"
func ParseGitLabServerVersion(ua string) string {
	if ua == "" {
		return ""
	}

	kv := strings.Split(ua, "/")
	if kv[0] != "GitLab" {
		return ua
	}

	if len(kv) > 1 {
		return kv[1]
	}

	return kv[0]
}
