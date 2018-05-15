package global

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

// OverrideScriptPathInBundle changes the root path in the bundle file.
// This is needed because preact by default bakes in an absolute src, like src="/". This results in issues with
// hosting the mouthful instance under a path that's different than /, such as "/mouthful". To fix it, on running mouthful
// this function is run and overrides the src with prefix + src.
func OverrideScriptPathInBundle(prefix, filepath string) error {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	newHTML := string(b)
	replWith := fmt.Sprintf(`e.p="%v"`, prefix)
	newHTML = strings.Replace(newHTML, `e.p="/"`, replWith, 1)
	err = ioutil.WriteFile(filepath, []byte(newHTML), 0644)
	return err
}

// FindAdminPanelChunkFilename returns the name of the route-panel chunk file, as it contains a hash.
func FindAdminPanelChunkFilename(root string) (string, error) {
	var files []string
	filepath.Walk(root, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if filepath.Ext(path) == ".js" && strings.HasPrefix(f.Name(), "bundle.") {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	if len(files) != 1 {
		return "", ErrCouldNotFindBundleFile
	}
	return files[0], nil
}

// RewriteAdminPanelScripts performs all the steps needed for rewriting the admin panel paths if your mouthful instance is not running under "/".
func RewriteAdminPanelScripts(path string) error {
	err := OverrideScriptRootInAdminHTML(path, StaticPath+"/index.html")
	if err != nil {
		return fmt.Errorf("Couldn't override the static admin html root")
	}
	fileName, err := FindAdminPanelChunkFilename(StaticPath)
	if err != nil {
		return fmt.Errorf("Couldn't find the admin panel chunk file")
	}
	err = OverrideScriptPathInBundle(path, StaticPath+"/"+fileName)
	if err != nil {
		return fmt.Errorf("Couldn't override the static admin script path")
	}
	return nil
}
