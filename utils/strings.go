package utils

import "strings"

/*
If you enter `git@github.com:Kruceo/Synco-Test.git`
You will receive `Synco-Test`
*/
func GetRepositoryNameFromGitURL(str string) string {
	//

	for i, v := range strings.Split(str, "/") {
		if b, found := strings.CutSuffix(v, ".git"); found && i != 0 {
			return b
		}
	}
	return ""
}
