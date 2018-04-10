package api

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/vkuznecovas/mouthful/global"
)

// UnmoderatedFs a file system we use to serve only the client.js
type UnmoderatedFs struct {
	http.FileSystem
}

// Exists determines if the file exists or not
func (cfs UnmoderatedFs) Exists(prefix string, filepath string) bool {
	if filepath != "/client.js" {
		return false
	}
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(global.StaticPath, p)
		_, err := os.Stat(name)
		if err != nil {
			return false
		}
		return true
	}
	return false
}
