package global_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/global"
)

const input = `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>admin</title><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="mobile-web-app-capable" content="yes"><meta name="apple-mobile-web-app-capable" content="yes"><link rel="manifest" href="/manifest.json"><meta name="theme-color" content="#673ab8"><link rel="shortcut icon" href="/favicon.ico"><link href="/style.c10b5.css" rel="stylesheet"></head><body><div id="app"><header class="header__3QGkI"><h1>Mouthful Admin Panel</h1></header><div class="mouthful_container__1eO2q"><div class="mouthful_login__fadGC">No comments yet!</div></div></div><script defer="defer" src="/bundle.fb3a6.js"></script><script>window.fetch||document.write('<script src="/polyfills.44897.js"><\/script>')</script></body></html>`
const prefix = `/test`
const expectedOutput = `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>admin</title><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="mobile-web-app-capable" content="yes"><meta name="apple-mobile-web-app-capable" content="yes"><link rel="manifest" href="/test/manifest.json"><meta name="theme-color" content="#673ab8"><link rel="shortcut icon" href="/test/favicon.ico"><link href="/test/style.c10b5.css" rel="stylesheet"></head><body><div id="app"><header class="header__3QGkI"><h1>Mouthful Admin Panel</h1></header><div class="mouthful_container__1eO2q"><div class="mouthful_login__fadGC">No comments yet!</div></div></div><script defer="defer" src="/test/bundle.fb3a6.js"></script><script>window.fetch||document.write('<script src="/test/polyfills.44897.js"><\/script>')</script></body></html>`

func TestOverrideScriptRootInAdminHtml(t *testing.T) {
	filepath := "./t.html"

	var _, err = os.Stat(filepath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(filepath)
		defer file.Close()
		assert.Nil(t, err)

		_, err = file.WriteString(input)
		assert.Nil(t, err)

		file.Close()
	}

	err = global.OverrideScriptRootInAdminHTML("/test", filepath)
	assert.Nil(t, err)

	b, err := ioutil.ReadFile(filepath)
	assert.Nil(t, err)

	newHtml := string(b)
	assert.Equal(t, expectedOutput, newHtml)

	err = global.OverrideScriptRootInAdminHTML("/test", filepath)
	assert.Nil(t, err)

	b, err = ioutil.ReadFile(filepath)
	assert.Nil(t, err)

	newHtml = string(b)
	assert.Equal(t, expectedOutput, newHtml)

	err = os.Remove(filepath)
	assert.Nil(t, err)
}
