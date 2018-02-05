package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/vkuznecovas/mouthful/db/abstraction"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/api"
	"github.com/vkuznecovas/mouthful/api/model"
	dbmodel "github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/db/sqlite"
)

func TestStatus(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()

	r.GET("/status").
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "{\"message\":\"OK\"}", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGetCommentsNoComments(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	r.GET("/comments?uri="+url.PathEscape("/2017/16")).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func TestGetCommentsBadQuery(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	r.GET("/comments?uri=").
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestCreateCommentSpamTrap(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	email := "email"
	body := model.CreateCommentBody{
		Path:   "2017/16",
		Body:   "body",
		Author: "author",
		Email:  &email,
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments/"+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "404 page not found", r.Body.String())
			assert.Equal(t, 404, r.Code)
		})
}

func TestCreateCommentBadRequst(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	r.POST("/comments").
		SetBody(string("sadasdasdasd")).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func GetSessionCookie(db *abstraction.Database, r *gofight.RequestConfig) gofight.H {
	cookiePrefix := "mouthful-session"
	cookieValue := ""
	os.Setenv("ADMIN_PASSWORD", "test")
	server := api.GetServer(db)
	r.POST("/admin/login").
		SetBody(`{"password": "test"}`).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			cookieValue = strings.Split(strings.TrimLeft(r.HeaderMap["Set-Cookie"][0], cookiePrefix+"="), " ")[0]
		})
	return gofight.H{cookiePrefix: cookieValue}
}

func TestDeleteCommentBadRequst(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server := api.GetServer(&testDB)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	r.DELETE("/comments").
		SetBody(string("sadasdasdasd")).
		SetCookie(cookies).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			fmt.Println(r.Code)
			assert.Equal(t, 400, r.Code)
		})
}

func TestDeleteCommentNonExistant(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server := api.GetServer(&testDB)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	r.DELETE("/comments").
		SetBody(string("{\"commentId\": 1}")).
		SetCookie(cookies).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func TestDeleteComment(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	server := api.GetServer(&testDB)

	body := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	conf := true
	bodyUpdate := model.UpdateCommentBody{
		CommentId: 1,
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	cookies := GetSessionCookie(&testDB, r)
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, "body", comments[0].Body)
			assert.Equal(t, "author", comments[0].Author)
		})
	r.DELETE("/comments").
		SetBody(string("{\"commentId\": 1}")).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 0)

		})
}

func TestUpdateCommentBadRequst(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server := api.GetServer(&testDB)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	r.PATCH("/comments").
		SetBody(string("sadasdasdasd")).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestCreateComment(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server := api.GetServer(&testDB)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	body := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	conf := true
	bodyUpdate := model.UpdateCommentBody{
		CommentId: 1,
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, "body", comments[0].Body)
			assert.Equal(t, "author", comments[0].Author)
		})
}

func TestLoginBadPassword(t *testing.T) {
	os.Setenv("ADMIN_PASSWORD", "test")
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := model.LoginBody{
		Password: "t",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/admin/login").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 401, r.Code)
		})
}

func TestLoginGoodPassword(t *testing.T) {
	os.Setenv("ADMIN_PASSWORD", "test")
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := model.LoginBody{
		Password: "test",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/admin/login").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.NotNil(t, r.HeaderMap["Set-Cookie"])
			assert.NotEqual(t, "", r.HeaderMap["Set-Cookie"])
			assert.True(t, strings.Contains(r.HeaderMap["Set-Cookie"][0], "mouthful-session="))
			assert.Equal(t, 204, r.Code)
		})
}

func TestLoginInvalidRequest(t *testing.T) {
	os.Setenv("ADMIN_PASSWORD", "test")
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := "asdasdasdasdasdasdasd"
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/admin/login").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestUpdateCommentUnauthorized(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := "Doesnt matter, really"
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 401, r.Code)
		})
}

func TestDeleteCommentUnauthorized(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := "Doesnt matter, really"
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.DELETE("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 401, r.Code)
		})
}

func TestGetThreadsUnauthorized(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := "Doesnt matter, really"
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.GET("/threads").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 401, r.Code)
		})
}

func TestGetThreadsEmpty(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server := api.GetServer(&testDB)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	body := "Doesnt matter, really"
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.GET("/threads").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var threads []dbmodel.Thread
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &threads)
			assert.Nil(t, err)
			assert.Len(t, threads, 0)
		})
}

func TestGetThreads(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server := api.GetServer(&testDB)
	r := gofight.New()
	body := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	cookies := GetSessionCookie(&testDB, r)
	assert.Nil(t, err)
	r.GET("/threads").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var threads []dbmodel.Thread
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &threads)
			assert.Nil(t, err)
			assert.Len(t, threads, 1)
			assert.Equal(t, "/1027/test/", threads[0].Path)
		})
}

func TestGetCommentsUnconfirmed(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	conf := true
	bodyUpdate := model.UpdateCommentBody{
		CommentId: 1,
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 0)
		})
}

func TestUpdateCommentNonExistant(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server := api.GetServer(&testDB)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)

	conf := true
	bodyUpdate := model.UpdateCommentBody{
		CommentId: 1,
		Confirmed: &conf,
	}
	bodyBytes, err := json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func TestUpdateCommentInvalidBody(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server := api.GetServer(&testDB)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)

	body := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})

	bodyUpdate := model.UpdateCommentBody{
		CommentId: 1,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}
func TestUpdateComment(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server := api.GetServer(&testDB)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	body := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})

	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 0)
		})
	conf := true
	bodyUpdate := model.UpdateCommentBody{
		CommentId: 1,
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, "body", comments[0].Body)
			assert.Equal(t, "author", comments[0].Author)
			assert.Equal(t, true, comments[0].Confirmed)
		})
	conf = true
	b := "test"
	bodyUpdate = model.UpdateCommentBody{
		CommentId: 1,
		Body:      &b,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, "test", comments[0].Body)
			assert.Equal(t, "author", comments[0].Author)
			assert.Equal(t, true, comments[0].Confirmed)
		})
	conf = true
	a := "test"
	bodyUpdate = model.UpdateCommentBody{
		CommentId: 1,
		Author:    &a,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, "test", comments[0].Body)
			assert.Equal(t, "test", comments[0].Author)
			assert.Equal(t, true, comments[0].Confirmed)
		})
}
