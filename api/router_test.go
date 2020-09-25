package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/vkuznecovas/mouthful/global"

	"github.com/vkuznecovas/mouthful/db/abstraction"
	"github.com/vkuznecovas/mouthful/db/dynamodb"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"

	"github.com/vkuznecovas/mouthful/api"
	"github.com/vkuznecovas/mouthful/api/model"
	configModel "github.com/vkuznecovas/mouthful/config/model"

	dbmodel "github.com/vkuznecovas/mouthful/db/model"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/mysql"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/postgres"
	"github.com/vkuznecovas/mouthful/db/sqlxDriver/sqlite"
)

const debug = false
const adminPassword = "test"

var maxCommentLength int = 10000

var config = configModel.Config{
	Honeypot: false,
	Moderation: configModel.Moderation{
		Enabled:          true,
		AdminPassword:    adminPassword,
		MaxCommentLength: &maxCommentLength,
	},
	API: configModel.API{
		Debug:   false,
		Logging: false,
		Cache: configModel.Cache{
			Enabled:           false,
			IntervalInSeconds: 1,
			ExpiryInSeconds:   2,
		},
		RateLimiting: configModel.RateLimiting{
			Enabled:   false,
			PostsHour: 2,
		},
	},
	Client: configModel.Client{
		UseDefaultStyle: true,
		PageSize:        10,
	},
}

var testFunctions = [...]interface{}{
	Status,
	GetCommentsNoComments,
	GetCommentsUnconfirmedComments,
	GetCommentsBadQuery,
	CreateCommentSpamTrap,
	CreateCommentBadRequst,
	DeleteCommentBadRequst,
	DeleteCommentNonExistant,
	DeleteComment,
	UpdateCommentBadRequst,
	CreateComment,
	CreateCommentBadReplyTo,
	CreateCommentReplyTo,
	LoginBadPassword,
	LoginGoodPassword,
	LoginInvalidRequest,
	GetAllCommentsUnauthorized,
	UpdateCommentUnauthorized,
	DeleteCommentUnauthorized,
	GetThreadsUnauthorized,
	GetThreadsEmpty,
	GetCommentsEmpty,
	GetThreads,
	GetComments,
	GetCommentsUnconfirmed,
	UpdateCommentNonExistant,
	UpdateCommentInvalidBody,
	UpdateComment,
	RestoreDeletedCommentBadRequst,
	RestoreDeletedCommentNonExistant,
	RestoreDeletedComment,
	CreateCommentBodyTooLong,
	GetCommentsCache,
	CreateCommentNoModeration,
	DeleteCommentBadUUID,
	UpdateCommentBadUUID,
	DeleteCommentDeletesReplyToComments,
	RateLimitingLoginCreation,
	RateLimitingDisabled,
	RateLimitingCommentCreation,
	GetCommentsWithPathNormalization,
	GetClientConfigReturnsConfig,
	CheckNoCorsSetting,
	CheckCorsSettingAllowsCorrectOrigin,
	CheckCorsSettingDoesNotAllowIncorrectOrigin,
	CreateCommentReplyToAReply,
	CreateCommentWithAuthorTooLongResultsInDefaultTruncation,
	CreateCommentWithAuthorTooLongTrucatesAccodringToConfig,
	CreateCommentEmptyBody,
	CreateCommentEmptyBodyAfterSanitization,
	CreateCommentEmptyAuthor,
	GetAdminConfig,
	OauthPathsExist,
	DeleteCommentHard,
}

func GetSessionCookie(db *abstraction.Database, r *gofight.RequestConfig) gofight.H {
	cookiePrefix := "mouthful-session"
	cookieValue := ""
	os.Setenv("ADMIN_PASSWORD", adminPassword)
	server, _ := api.GetServer(db, &config)
	r.POST("/v1/admin/login").
		SetBody(fmt.Sprintf(`{"password": "%v"}`, adminPassword)).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			cookieValue = strings.TrimSuffix(strings.Split(strings.TrimLeft(r.HeaderMap["Set-Cookie"][0], cookiePrefix+"="), " ")[0], ";")
		})
	return gofight.H{cookiePrefix: cookieValue}
}

func setupDynamoTestDb() abstraction.Database {
	database := dynamodb.CreateTestDatabase()
	err := database.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	return database
}

func setupSqliteTestDb() abstraction.Database {
	database := sqlite.CreateTestDatabase()
	return database
}

func TestRouterWithSqlite(t *testing.T) {
	for _, f := range testFunctions {
		f.(func(*testing.T, abstraction.Database))(t, setupSqliteTestDb())
	}
}

func TestRouterWithDynamoDb(t *testing.T) {
	db := setupDynamoTestDb()
	driver := db.GetUnderlyingStruct()
	driverCasted := driver.(*dynamodb.Database)
	for _, f := range testFunctions {
		f.(func(*testing.T, abstraction.Database))(t, db)
		driverCasted.WipeOutData()
	}
}

func TestRouterWithPostgresDb(t *testing.T) {
	db := postgres.CreateTestDatabase()
	driver := db.GetUnderlyingStruct()
	driverCasted := driver.(*sqlxDriver.Database)
	// clean out before start
	driverCasted.WipeOutData()
	for _, f := range testFunctions {
		f.(func(*testing.T, abstraction.Database))(t, db)
		err := driverCasted.WipeOutData()
		assert.Nil(t, err)
	}
}

func TestRouterWithMysqlDb(t *testing.T) {
	db := mysql.CreateTestDatabase()
	driver := db.GetUnderlyingStruct()
	driverCasted := driver.(*sqlxDriver.Database)
	// clean out before start
	driverCasted.WipeOutData()
	for _, f := range testFunctions {
		f.(func(*testing.T, abstraction.Database))(t, db)
		err := driverCasted.WipeOutData()
		assert.Nil(t, err)
	}
}

func Status(t *testing.T, testDB abstraction.Database) {
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

func GetCommentsNoComments(t *testing.T, testDB abstraction.Database) {
	r := gofight.New()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.GET("/v1/comments?uri="+url.PathEscape("/2017/16")).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func GetCommentsUnconfirmedComments(t *testing.T, testDB abstraction.Database) {
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
			assert.Equal(t, body.Path+"/", parsedBody.Path)
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

func GetCommentsBadQuery(t *testing.T, testDB abstraction.Database) {
	r := gofight.New()
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r.GET("/v1/comments?uri=").
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func CreateCommentSpamTrap(t *testing.T, testDB abstraction.Database) {
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
			assert.Equal(t, body.Path+"/", parsedBody.Path)
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

func CreateCommentBadRequst(t *testing.T, testDB abstraction.Database) {
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

func DeleteCommentBadRequst(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	r.DELETE("/v1/admin/comments").
		SetBody(string("sadasdasdasd")).
		SetCookie(cookies).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func DeleteCommentNonExistant(t *testing.T, testDB abstraction.Database) {
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

func DeleteComment(t *testing.T, testDB abstraction.Database) {
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

func DeleteCommentHard(t *testing.T, testDB abstraction.Database) {
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
		Hard:      true,
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
	r.GET("/v1/admin/comments/all").
		SetDebug(debug).
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

func UpdateCommentBadRequst(t *testing.T, testDB abstraction.Database) {
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

func CreateComment(t *testing.T, testDB abstraction.Database) {
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

func CreateCommentBadReplyTo(t *testing.T, testDB abstraction.Database) {
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

func CreateCommentReplyTo(t *testing.T, testDB abstraction.Database) {
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
			assert.True(t, bytes.Equal(commentId.Bytes(), comments[1].ReplyTo.Bytes()))
		})
}

func LoginBadPassword(t *testing.T, testDB abstraction.Database) {
	os.Setenv("ADMIN_PASSWORD", "test")
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

func LoginGoodPassword(t *testing.T, testDB abstraction.Database) {
	os.Setenv("ADMIN_PASSWORD", adminPassword)
	r := gofight.New()
	body := model.LoginBody{
		Password: adminPassword,
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

func LoginInvalidRequest(t *testing.T, testDB abstraction.Database) {
	os.Setenv("ADMIN_PASSWORD", adminPassword)
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

func GetAllCommentsUnauthorized(t *testing.T, testDB abstraction.Database) {
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

func UpdateCommentUnauthorized(t *testing.T, testDB abstraction.Database) {
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

func DeleteCommentUnauthorized(t *testing.T, testDB abstraction.Database) {
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

func GetThreadsUnauthorized(t *testing.T, testDB abstraction.Database) {
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

func GetThreadsEmpty(t *testing.T, testDB abstraction.Database) {
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

func GetCommentsEmpty(t *testing.T, testDB abstraction.Database) {
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

func GetThreads(t *testing.T, testDB abstraction.Database) {
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

func GetComments(t *testing.T, testDB abstraction.Database) {
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

func GetCommentsUnconfirmed(t *testing.T, testDB abstraction.Database) {
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

func UpdateCommentNonExistant(t *testing.T, testDB abstraction.Database) {
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

func UpdateCommentInvalidBody(t *testing.T, testDB abstraction.Database) {
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
func UpdateComment(t *testing.T, testDB abstraction.Database) {
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

func RestoreDeletedCommentBadRequst(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	r.POST("/v1/admin/comments/restore").
		SetBody(string("sadasdasdasd")).
		SetCookie(cookies).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func RestoreDeletedCommentNonExistant(t *testing.T, testDB abstraction.Database) {
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

func RestoreDeletedComment(t *testing.T, testDB abstraction.Database) {
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

func CreateCommentBodyTooLong(t *testing.T, testDB abstraction.Database) {
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
			assert.Equal(t, body.Path+"/", parsedBody.Path)
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

func GetCommentsCache(t *testing.T, testDB abstraction.Database) {
	newConfig := config
	newConfig.API.Cache.Enabled = true
	server, err := api.GetServer(&testDB, &newConfig)
	assert.Nil(t, err)
	r := gofight.New()
	commentBody := model.CreateCommentBody{
		Path:   "1027/test",
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
			assert.Equal(t, "/"+commentBody.Path+"/", parsedBody.Path)
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
	r.GET("/v1/comments?uri="+url.QueryEscape("/"+commentBody.Path+"/")).
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

func CreateCommentNoModeration(t *testing.T, testDB abstraction.Database) {
	newConfig := config
	newConfig.Moderation.Enabled = false
	server, err := api.GetServer(&testDB, &newConfig)
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
			var parsedBody model.CreateCommentResponse
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

func DeleteCommentBadUUID(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	body := model.DeleteCommentBody{
		CommentId: "535622c2-2da5-11e8-b4670ed5f89f718b",
	}
	bodyJson, err := json.Marshal(body)
	assert.Nil(t, err)
	r.DELETE("/v1/admin/comments").
		SetBody(string(bodyJson)).
		SetCookie(cookies).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func UpdateCommentBadUUID(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	newBody := "test"
	body := model.UpdateCommentBody{
		CommentId: "535622c2-2da5-11e8-b4670ed5f89f718b",
		Body:      &newBody,
	}
	bodyJson, err := json.Marshal(body)
	assert.Nil(t, err)
	r.PATCH("/v1/admin/comments").
		SetBody(string(bodyJson)).
		SetCookie(cookies).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func DeleteCommentDeletesReplyToComments(t *testing.T, testDB abstraction.Database) {
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

	cookies := GetSessionCookie(&testDB, r)

	replyTo := commentId.String()
	body = model.CreateCommentBody{
		Path:    "/1027/test/",
		Body:    "body",
		Author:  "author",
		ReplyTo: &replyTo,
	}
	bodyBytes, err = json.Marshal(body)
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
			_, err := global.ParseUUIDFromString(parsedBody.Id)
			assert.Nil(t, err)
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
	r.GET("/v1/admin/comments/all").
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var comments []dbmodel.Comment
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &comments)
			assert.Nil(t, err)
			assert.Len(t, comments, 2)
			for _, v := range comments {
				assert.NotNil(t, v.DeletedAt)
			}
		})
}

func RateLimitingCommentCreation(t *testing.T, testDB abstraction.Database) {
	newConfig := config
	newConfig.API.RateLimiting.Enabled = true
	newConfig.API.RateLimiting.PostsHour = 10
	server, err := api.GetServer(&testDB, &newConfig)
	assert.Nil(t, err)
	r := gofight.New()
	body := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	for i := 0; i <= newConfig.API.RateLimiting.PostsHour; i++ {
		r.POST("/v1/comments").
			SetBody(string(bodyBytes[:])).
			SetDebug(debug).
			Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				if i == newConfig.API.RateLimiting.PostsHour {
					assert.Equal(t, 429, r.Code)
				} else {
					assert.Equal(t, 200, r.Code)
				}
			})
	}

}

func RateLimitingLoginCreation(t *testing.T, testDB abstraction.Database) {
	newConfig := config
	newConfig.API.RateLimiting.Enabled = true
	newConfig.API.RateLimiting.PostsHour = 100
	server, err := api.GetServer(&testDB, &newConfig)
	assert.Nil(t, err)
	r := gofight.New()
	os.Setenv("ADMIN_PASSWORD", "test")
	body := model.LoginBody{
		Password: "t",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	for i := 0; i <= newConfig.API.RateLimiting.PostsHour; i++ {
		r.POST("/v1/admin/login").
			SetBody(string(bodyBytes[:])).
			SetDebug(debug).
			Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				if i == newConfig.API.RateLimiting.PostsHour {
					assert.Equal(t, 429, r.Code)
				} else {
					assert.Equal(t, 401, r.Code)
				}
			})
	}
}

func RateLimitingDisabled(t *testing.T, testDB abstraction.Database) {
	newConfig := config
	newConfig.API.RateLimiting.Enabled = false
	newConfig.API.RateLimiting.PostsHour = 1000
	server, err := api.GetServer(&testDB, &newConfig)
	assert.Nil(t, err)
	r := gofight.New()
	os.Setenv("ADMIN_PASSWORD", "test")
	body := model.LoginBody{
		Password: "t",
	}
	bodyBytes, err := json.Marshal(body)
	assert.Nil(t, err)
	for i := 0; i <= newConfig.API.RateLimiting.PostsHour; i++ {
		r.POST("/v1/admin/login").
			SetBody(string(bodyBytes[:])).
			SetDebug(debug).
			Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, 401, r.Code)
			})
	}
}

func GetCommentsWithPathNormalization(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	commentBody := model.CreateCommentBody{
		Path:   "1027/test",
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
			assert.Equal(t, "/"+commentBody.Path+"/", parsedBody.Path)
			assert.Equal(t, commentBody.Author, parsedBody.Author)
			assert.Equal(t, global.ParseAndSaniziteMarkdown(commentBody.Body), parsedBody.Body)
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
			assert.Equal(t, "/"+commentBody.Path+"/", threads[0].Path)
		})
}

func GetClientConfigReturnsConfig(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()

	r.GET("/v1/client/config").
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var cc configModel.ClientConfig
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &cc)
			assert.Nil(t, err)
			assert.Equal(t, false, cc.Honeypot)
			assert.Equal(t, true, cc.UseDefaultStyle)
			assert.Equal(t, true, cc.Moderation)
			assert.Equal(t, 10000, *cc.MaxCommentLength)
			assert.Equal(t, 10, cc.PageSize)
		})
}

func CheckNoCorsSetting(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()

	r.GET("/v1/client/config").
		SetHeader(gofight.H{"Origin": "http://google.co.uk"}).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			var cc configModel.ClientConfig
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			err = json.Unmarshal(body, &cc)
			assert.Nil(t, err)
			assert.Equal(t, false, cc.Honeypot)
			assert.Equal(t, true, cc.UseDefaultStyle)
			assert.Equal(t, true, cc.Moderation)
			assert.Equal(t, 10000, *cc.MaxCommentLength)
			assert.Equal(t, 10, cc.PageSize)
			assert.Equal(t, "*", r.HeaderMap.Get("Access-Control-Allow-Origin"))
		})
}

func CheckCorsSettingAllowsCorrectOrigin(t *testing.T, testDB abstraction.Database) {
	c := config
	origins := []string{"http://google.co.uk", "https://google.co.uk"}
	c.API.Cors = configModel.Cors{
		Enabled:        true,
		AllowedOrigins: &origins,
	}
	server, err := api.GetServer(&testDB, &c)
	assert.Nil(t, err)
	r := gofight.New()

	for _, v := range origins {
		r.GET("/v1/client/config").
			SetHeader(gofight.H{"Origin": v}).
			SetDebug(debug).
			Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, 200, r.Code)
				var cc configModel.ClientConfig
				body, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				err = json.Unmarshal(body, &cc)
				assert.Nil(t, err)
				assert.Equal(t, false, cc.Honeypot)
				assert.Equal(t, true, cc.UseDefaultStyle)
				assert.Equal(t, true, cc.Moderation)
				assert.Equal(t, 10000, *cc.MaxCommentLength)
				assert.Equal(t, 10, cc.PageSize)
				assert.Equal(t, v, r.HeaderMap.Get("Access-Control-Allow-Origin"))
			})
	}

}

func CheckCorsSettingDoesNotAllowIncorrectOrigin(t *testing.T, testDB abstraction.Database) {
	c := config
	origins := []string{"http://google.co.uk", "https://google.co.uk"}
	c.API.Cors = configModel.Cors{
		Enabled:        true,
		AllowedOrigins: &origins,
	}
	server, err := api.GetServer(&testDB, &c)
	assert.Nil(t, err)
	r := gofight.New()

	r.GET("/v1/client/config").
		SetHeader(gofight.H{"Origin": "http://example.com"}).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 403, r.Code)
		})

}

func CreateCommentReplyToAReply(t *testing.T, testDB abstraction.Database) {
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
	replyToReplyId := commentId2.String()
	body3 := model.CreateCommentBody{
		Path:    "/1027/test/",
		Body:    "body",
		Author:  "author",
		ReplyTo: &replyToReplyId,
	}
	bodyBytes, err = json.Marshal(body3)
	var commentId3 uuid.UUID
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
			commentId3 = *uid
		})
	bodyUpdate = model.UpdateCommentBody{
		CommentId: commentId3.String(),
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
			assert.Len(t, comments, 3)
			assert.Equal(t, global.ParseAndSaniziteMarkdown("body"), comments[0].Body)
			assert.Nil(t, comments[0].ReplyTo)
			assert.True(t, bytes.Equal(commentId.Bytes(), comments[1].ReplyTo.Bytes()))
			assert.True(t, bytes.Equal(commentId.Bytes(), comments[2].ReplyTo.Bytes()))
		})
}

func CreateCommentWithAuthorTooLongResultsInDefaultTruncation(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	commentBody := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "authorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthor",
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
			assert.Equal(t, "authorauthorauthorauthorauthorauthorauthorautho...", parsedBody.Author)
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
			assert.Equal(t, "authorauthorauthorauthorauthorauthorauthorautho...", comments[0].Author)
			assert.Equal(t, 50, len(comments[0].Author))
			assert.Equal(t, global.ParseAndSaniziteMarkdown(commentBody.Body), comments[0].Body)
		})
}

func CreateCommentWithAuthorTooLongTrucatesAccodringToConfig(t *testing.T, testDB abstraction.Database) {
	ml := 10
	configCopy := config
	configCopy.Moderation.MaxAuthorLength = &ml
	server, err := api.GetServer(&testDB, &configCopy)
	assert.Nil(t, err)
	r := gofight.New()
	commentBody := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "body",
		Author: "authorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthorauthor",
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
			assert.Equal(t, "authora...", parsedBody.Author)
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
			assert.Equal(t, "authora...", comments[0].Author)
			assert.Equal(t, ml, len(comments[0].Author))
			assert.Equal(t, global.ParseAndSaniziteMarkdown(commentBody.Body), comments[0].Body)
		})

}

func CreateCommentEmptyBody(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	commentBody := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "",
		Author: "author",
	}
	bodyBytes, err := json.Marshal(commentBody)
	assert.Nil(t, err)
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
			b, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "\"Bad request\"", string(b))
		})
}

func CreateCommentEmptyBodyAfterSanitization(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	commentBody := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   `<script>alert("qweqweqwe")</script>`,
		Author: "author",
	}
	bodyBytes, err := json.Marshal(commentBody)
	assert.Nil(t, err)
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
			b, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "\"Bad request\"", string(b))
		})
}

func CreateCommentEmptyAuthor(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	commentBody := model.CreateCommentBody{
		Path:   "/1027/test/",
		Body:   "asdasds",
		Author: "",
	}
	bodyBytes, err := json.Marshal(commentBody)
	assert.Nil(t, err)
	r.POST("/v1/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
			b, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "\"Bad request\"", string(b))
		})
}

func GetAdminConfig(t *testing.T, testDB abstraction.Database) {
	server, err := api.GetServer(&testDB, &config)
	assert.Nil(t, err)
	r := gofight.New()
	cookies := GetSessionCookie(&testDB, r)
	r.GET("/v1/admin/config").
		SetHeader(gofight.H{"Origin": "http://google.co.uk"}).
		SetDebug(debug).
		SetCookie(cookies).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 200, r.Code)
			body, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			var adminConfig configModel.AdminConfig
			err = json.Unmarshal(body, &adminConfig)
			assert.Nil(t, err)
			assert.False(t, adminConfig.DisablePasswordLogin)
			assert.Len(t, *adminConfig.OauthProviders, 0)
		})
}

func OauthPathsExist(t *testing.T, testDB abstraction.Database) {
	configCopy := config
	configCopy.Moderation.OAauthProviders = &someFakeOauthProviders
	configCopy.Moderation.OAuthCallbackOrigin = &fakeOrigin
	server, err := api.GetServer(&testDB, &configCopy)
	assert.Nil(t, err)
	r := gofight.New()
	r.GET("/v1/oauth/auth/github").
		SetHeader(gofight.H{"Origin": "http://google.co.uk"}).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 307, r.Code)
		})
	r.GET("/v1/oauth/callbacks/github").
		SetHeader(gofight.H{"Origin": "http://google.co.uk"}).
		SetDebug(debug).
		Run(server, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 500, r.Code)
		})
}
