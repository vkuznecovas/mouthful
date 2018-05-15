package global

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// OverrideScriptRootInAdminHTML changes the src elements in the index.html of the admin panel inside the static dir.
// This is needed because preact by default bakes in an absolute src, like src="/". This results in issues with
// hosting the mouthful instance under a path that's different than /, such as "/mouthful". To fix it, on running mouthful
// this function is run and overrides the src with prefix + src.
func OverrideScriptRootInAdminHTML(prefix, filepath string) error {
	scriptOverridePattern := `.*?\"(\/.*?\")`
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	r, err := regexp.Compile(scriptOverridePattern)
	if err != nil {
		return err
	}

	newHTML := string(b)
	res := r.FindAllSubmatch(b, -1)
	for _, v := range res {
		m2 := string(v[1])
		if strings.HasPrefix(m2, prefix) {
			return nil
		}
		newHTML = strings.Replace(newHTML, m2, prefix+m2, 1)
	}

	err = ioutil.WriteFile(filepath, []byte(newHTML), 0644)
	return nil
}
