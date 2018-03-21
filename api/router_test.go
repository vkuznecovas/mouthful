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

	"github.com/satori/go.uuid"
	"github.com/vkuznecovas/mouthful/global"

	"github.com/vkuznecovas/mouthful/db/abstraction"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/api"
	"github.com/vkuznecovas/mouthful/api/model"
	configModel "github.com/vkuznecovas/mouthful/config/model"

	dbmodel "github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/db/sqlite"
)

const debug = false

var maxCommentLength int = 10000

var config = configModel.Config{
	Honeypot: false,
	Moderation: configModel.Moderation{
		Enabled:          true,
		SessionSecret:    "somesecret",
		AdminPassword:    "test",
		MaxCommentLength: &maxCommentLength,
	},
	API: configModel.API{
		Debug: false,
		Cache: configModel.Cache{
			Enabled:           false,
			IntervalInSeconds: 1,
			ExpiryInSeconds:   2,
		},
	},
}

func TestStatus(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.GET("/status").
		SetDebug(debug).
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
	r.GET("/v1/comments?uri="+url.PathEscape("/2017/16")).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func TestGetCommentsUnconfirmedComments(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	body := model.CreateCommentBody{
		Path:   "/2017/16",
		Body:   "body",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentBody
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
		})
	r.GET("/v1/comments?uri="+url.PathEscape("/2017/16")).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func TestGetCommentsBadQuery(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.GET("/v1/comments?uri=").
		SetDebug(debug).
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
		Path:   "/2017/16",
		Body:   "body",
		Author: "author",
		Email:  &email,
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentBody
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, body.Email, parsedBody.Email)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
		})
	r.GET("/v1/comments/"+url.QueryEscape(body.Path)).
		SetDebug(debug).
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
	r.POST("/v1/comments").
		SetBody(string("sadasdasdasd")).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func GetSessionCookie(db *abstraction.Database, r *gofight.RequestConfig) gofight.H {
	cookiePrefix := "mouthful-session"
	cookieValue := ""
	os.Setenv("ADMIN_PASSWORD", "test")
	server, _ := api.GetServer(db, &config)
	r.POST("/v1/admin/login").
		SetBody(`{"password": "test"}`).
		SetDebug(debug).
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
	r.DELETE("/v1/admin/comments").
		SetBody(string("sadasdasdasd")).
		SetCookie(cookies).
		SetDebug(debug).
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
	deleteCommentBody := model.DeleteCommentBody{
		CommentId: global.GetUUID().String(),
	}
	v, err := json.Marshal(deleteCommentBody)
	assert.Nil(t, err)
	r.DELETE("/v1/admin/comments").
		SetBody(string(v)).
		SetCookie(cookies).
		SetDebug(debug).
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
	var commentId uuid.UUID
	assert.Nil(t, err)
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentResponse
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
			uid, err := global.ParseUUIDFromString(parsedBody.Id)
			assert.Nil(t, err)
			commentId = *uid
		})
	conf := true

	bodyUpdate := model.UpdateCommentBody{
		CommentId: commentId.String(),
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	cookies := GetSessionCookie(&testDB, r)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, global.ParseAndSaniziteMarkdown("body"), comments[0].Body)
			assert.Equal(t, "author", comments[0].Author)
		})
	deleteCommentBody := model.DeleteCommentBody{
		CommentId: commentId.String(),
	}
	v, err := json.Marshal(deleteCommentBody)
	assert.Nil(t, err)
	r.DELETE("/v1/admin/comments").
		SetBody(string(v)).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func TestUpdateCommentBadRequst(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	r.PATCH("/v1/admin/comments").
		SetBody(string("sadasdasdasd")).
		SetDebug(debug).
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
	var commentId uuid.UUID
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentResponse
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
			uid, err := global.ParseUUIDFromString(parsedBody.Id)
			assert.Nil(t, err)
			commentId = *uid
		})
	conf := true
	bodyUpdate := model.UpdateCommentBody{
		CommentId: commentId.String(),
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, global.ParseAndSaniziteMarkdown("body"), comments[0].Body)
			assert.Equal(t, "author", comments[0].Author)
		})
}

func TestCreateCommentBadReplyTo(t *testing.T) {
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
	var commentId uuid.UUID
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentResponse
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
			uid, err := global.ParseUUIDFromString(parsedBody.Id)
			assert.Nil(t, err)
			commentId = *uid
		})

	conf := true

	bodyUpdate := model.UpdateCommentBody{
		CommentId: commentId.String(),
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	replyTo := commentId.String()
	body2 := model.CreateCommentBody{
		Path:    "/1027/test/tttttttttt",
		Body:    "body",
		Author:  "author",
		ReplyTo: &replyTo,
	}
	bodyBytes, err = json.Marshal(body2)
	assert.Nil(t, err)
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "\"Can't reply to this comment\"", r.Body.String())
			assert.Equal(t, 400, r.Code)
		})
}

func TestCreateCommentReplyTo(t *testing.T) {
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
	var commentId uuid.UUID
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentResponse
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
			uid, err := global.ParseUUIDFromString(parsedBody.Id)
			assert.Nil(t, err)
			commentId = *uid
		})

	conf := true
	bodyUpdate := model.UpdateCommentBody{
		CommentId: commentId.String(),
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	replyTo := commentId.String()
	body2 := model.CreateCommentBody{
		Path:    "/1027/test/",
		Body:    "body",
		Author:  "author",
		ReplyTo: &replyTo,
	}
	bodyBytes, err = json.Marshal(body2)
	var commentId2 uuid.UUID
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentResponse
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body2.Path, parsedBody.Path)
			assert.Equal(t, body2.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body2.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
			uid, err := global.ParseUUIDFromString(parsedBody.Id)
			assert.Nil(t, err)
			commentId2 = *uid
		})
	bodyUpdate = model.UpdateCommentBody{
		CommentId: commentId2.String(),
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 2)
			assert.Equal(t, global.ParseAndSaniziteMarkdown("body"), comments[0].Body)
			assert.Nil(t, comments[0].ReplyTo)
			assert.True(t, uuid.Equal(commentId, *comments[1].ReplyTo))
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
	r.POST("/v1/admin/login").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
	r.POST("/v1/admin/login").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
	r.POST("/v1/admin/login").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
	r.GET("/v1/admin/comments/all").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
	r.DELETE("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
	r.GET("/v1/admin/threads").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
	r.GET("/v1/admin/threads").
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			assert.Equal(t, "[]", r.Body.String())
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
	r.GET("/v1/admin/comments/all").
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			assert.Equal(t, "[]", r.Body.String())
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
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentBody
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
		})
	cookies := GetSessionCookie(&testDB, r)
	assert.Nil(t, err)
	r.GET("/v1/admin/threads").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentBody
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, commentBody.Path, parsedBody.Path)
			assert.Equal(t, commentBody.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(commentBody.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
		})
	cookies := GetSessionCookie(&testDB, r)
	assert.Nil(t, err)
	r.GET("/v1/admin/comments/all").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
			assert.Equal(t, global.ParseAndSaniziteMarkdown(commentBody.Body), comments[0].Body)
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
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentBody
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
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
		CommentId: global.GetUUID().String(),
		Confirmed: &conf,
	}
	bodyBytes, err := json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
	var commentId uuid.UUID
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentResponse
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
			uid, err := global.ParseUUIDFromString(parsedBody.Id)
			assert.Nil(t, err)
			commentId = *uid
		})

	bodyUpdate := model.UpdateCommentBody{
		CommentId: commentId.String(),
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
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
	var commentId uuid.UUID
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentResponse
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
			uid, err := global.ParseUUIDFromString(parsedBody.Id)
			assert.Nil(t, err)
			commentId = *uid
		})

	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
	conf := true
	bodyUpdate := model.UpdateCommentBody{
		CommentId: commentId.String(),
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, global.ParseAndSaniziteMarkdown("body"), comments[0].Body)
			assert.Equal(t, "author", comments[0].Author)
			assert.Equal(t, true, comments[0].Confirmed)
		})
	conf = true
	b := "test"
	bodyUpdate = model.UpdateCommentBody{
		CommentId: commentId.String(),
		Body:      &b,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
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
		CommentId: commentId.String(),
		Author:    &a,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
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

func TestRestoreDeletedCommentBadRequst(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	r.POST("/v1/admin/comments/restore").
		SetBody(string("sadasdasdasd")).
		SetCookie(cookies).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			fmt.Println(r.Code)
			assert.Equal(t, 400, r.Code)
		})
}

func TestRestoreDeletedCommentNonExistant(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	deleteCommentBody := model.DeleteCommentBody{
		CommentId: global.GetUUID().String(),
	}
	v, err := json.Marshal(deleteCommentBody)
	assert.Nil(t, err)
	r.DELETE("/v1/admin/comments/restore").
		SetBody(string(v)).
		SetCookie(cookies).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func TestRestoreDeletedComment(t *testing.T) {
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
	var commentId uuid.UUID
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentResponse
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
			uid, err := global.ParseUUIDFromString(parsedBody.Id)
			assert.Nil(t, err)
			commentId = *uid
		})
	conf := true
	bodyUpdate := model.UpdateCommentBody{
		CommentId: commentId.String(),
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	cookies := GetSessionCookie(&testDB, r)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, global.ParseAndSaniziteMarkdown("body"), comments[0].Body)
			assert.Equal(t, "author", comments[0].Author)
		})
	deleteCommentBody := model.DeleteCommentBody{
		CommentId: commentId.String(),
	}
	v, err := json.Marshal(deleteCommentBody)
	r.DELETE("/v1/admin/comments").
		SetBody(string(v)).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
			assert.Equal(t, "\"Thread not found\"", r.Body.String())
		})
	r.POST("/v1/admin/comments/restore").
		SetBody(string(v)).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)

		})

}

func TestCreateCommentBodyTooLong(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	email := "email"
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	body := model.CreateCommentBody{
		Path:   "/2017/16",
		Body:   strings.Repeat("#", maxCommentLength),
		Author: "author",
		Email:  &email,
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentBody
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, body.Path, parsedBody.Path)
			assert.Equal(t, body.Author, parsedBody.Author)
			assert.Equal(t, body.Email, parsedBody.Email)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(body.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
		})
	body = model.CreateCommentBody{
		Path:   "/2017/16",
		Body:   strings.Repeat("#", maxCommentLength+1),
		Author: "author",
		Email:  &email,
	}
	bodyBytes, err = json.Marshal(body)
	assert.Nil(t, err)
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestGetCommentsCache(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	newConfig := config
	newConfig.API.Cache.Enabled = true
	server, err := api.GetServer(&testDB, &newConfig)
	assert.Nil(t, err)
	r := gofight.New()
	commentBody := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(commentBody)
	assert.Nil(t, err)
	var commentId uuid.UUID
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			var parsedBody model.CreateCommentResponse
			err = json.Unmarshal([]byte(r.Body.String()), &parsedBody)
			assert.Nil(t, err)
			assert.Equal(t, commentBody.Path, parsedBody.Path)
			assert.Equal(t, commentBody.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(commentBody.Body), parsedBody.Body)
			assert.Equal(t, 200, r.Code)
			uid, err := global.ParseUUIDFromString(parsedBody.Id)
			assert.Nil(t, err)
			commentId = *uid
		})
	cookies := GetSessionCookie(&testDB, r)
	conf := true
	bodyUpdate := model.UpdateCommentBody{
		CommentId: commentId.String(),
		Confirmed: &conf,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(commentBody.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			val := r.HeaderMap.Get("X-Cache")
			assert.Equal(t, val, "MISS")
		})
	r.GET("/v1/comments?uri="+url.QueryEscape(commentBody.Path)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 1)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(commentBody.Body), comments[0].Body)
			val := r.HeaderMap.Get("X-Cache")
			assert.Equal(t, val, "HIT")
		})
}
