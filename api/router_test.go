package api_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

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

func TestDeleteCommentBadRequst(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	r.DELETE("/comments").
		SetBody(string("sadasdasdasd")).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestDeleteCommentNonExistant(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	r.DELETE("/comments").
		SetBody(string("{\"commentId\": 1}")).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func TestDeleteComment(t *testing.T) {
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
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
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

func TestUpdateCommentBadRequst(t *testing.T) {
	testDB := sqlite.CreateTestDatabase()
	r := gofight.New()
	r.PATCH("/comments").
		SetBody(string("sadasdasdasd")).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}

func TestCreateComment(t *testing.T) {
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
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "", r.Body.String())
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
	r := gofight.New()

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
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}

func TestUpdateCommentInvalidBody(t *testing.T) {
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

	bodyUpdate := model.UpdateCommentBody{
		CommentId: 1,
	}
	bodyBytes, err = json.Marshal(bodyUpdate)
	assert.Nil(t, err)
	r.PATCH("/comments").
		SetBody(string(bodyBytes[:])).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 400, r.Code)
		})
}
func TestUpdateComment(t *testing.T) {
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
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 204, r.Code)
		})
	r.GET("/comments?uri="+url.QueryEscape(body.Path)).
		SetDebug(true).
		Run(api.GetServer(&testDB), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
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
