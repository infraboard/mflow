package job

const INJECT_PROVIDER = "@GitlabEvent"

var InjectParams = map[string]string{
	"GIT_SSH_URL":   INJECT_PROVIDER,
	"GIT_BRANCH":    INJECT_PROVIDER,
	"GIT_COMMIT_ID": INJECT_PROVIDER,
	"APP_VERSION":   INJECT_PROVIDER,
}
