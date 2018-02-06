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
	configModel "github.com/vkuznecovas/mouthful/config/model"

	dbmodel "github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/db/sqlite"
)

var config = configModel.Config{
	Honeypot: false,
	Moderation: configModel.Moderation{
		Enabled:       true,
		SessionSecret: "somesecret",
		AdminPassword: "test",
	},
}

func TestStatus(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.GET("/status").
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "{\"message\":\"OK\"}", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestGetCommentsNoComments(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.GET("/comments?uri="+url.PathEscape("/2017/16")).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func TestGetCommentsBadQuery(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.GET("/comments?uri=").
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestCreateCommentSpamTrap(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	email := "email"

	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
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
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments/"+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "404 page not found", r.Body.String())
			assert.Equal(t, 404, r.Code)
		})
}

func TestCreateCommentBadRequst(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.POST("/comments").
		SetBody(string("sadasdasdasd")).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func GetSessionCookie(db *abstraction.Database, r *gofight.RequestConfig) gofight.H {
	cookiePrefix := "mouthful-session"
	cookieValue := ""
	os.Setenv("ADMIN_PASSWORD", "test")
	server, _ := api.GetServer(db, &config)
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
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
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
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
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
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)

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
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
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
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
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
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.POST("/admin/login").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	assert.Nil(t, err)
	r.POST("/admin/login").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.POST("/admin/login").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestGetAllCommentsUnauthorized(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := "Doesnt matter, really"
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.GET("/comments/all").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 401, r.Code)
		})
}

func TestUpdateCommentUnauthorized(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := "Doesnt matter, really"
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 401, r.Code)
		})
}

func TestDeleteCommentUnauthorized(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := "Doesnt matter, really"
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.DELETE("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 401, r.Code)
		})
}

func TestGetThreadsUnauthorized(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	body := "Doesnt matter, really"
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.GET("/threads").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 401, r.Code)
		})
}

func TestGetThreadsEmpty(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	r.GET("/threads").
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

func TestGetCommentsEmpty(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	r.GET("/comments/all").
		SetDebug(true).
		SetCookie(cookies).
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

func TestGetThreads(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
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

func TestGetComments(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	commentBody := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(commentBody)
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
	r.GET("/comments/all").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, commentBody.Author, comments[0].Author)
			assert.Equal(t, commentBody.Body, comments[0].Body)
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
	server, err := api.GetServer(&testDB, &config)
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

func TestUpdateCommentNonExistant(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
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
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
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
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
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
