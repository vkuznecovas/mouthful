package sqlite_test

import (
	"os"
	"testing"

	"github.com/vkuznecovas/mouthful/db/sqlxDriver/sqlite"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/afero"
)

type Database struct {
	Dialect                   string  `json:"dialect"`
	Database                  *string `json:"database,omitempty"`
	Username                  *string `json:"username,omitempty"`
	Password                  *string `json:"password,omitempty"`
	Host                      *string `json:"host,omitempty"`
	Port                      *string `json:"port,omitempty"`
	TablePrefix               *string `json:"tablePrefix,omitempty"`
	DynamoDBThreadReadUnits   *int64  `json:"dynamoDBThreadReadUnits,omitempty"`
	DynamoDBCommentReadUnits  *int64  `json:"dynamoDBCommentReadUnits,omitempty"`
	DynamoDBThreadWriteUnits  *int64  `json:"dynamoDBThreadWriteUnits,omitempty"`
	DynamoDBCommentWriteUnits *int64  `json:"dynamoDBCommentWriteUnits,omitempty"`
	DynamoDBIndexWriteUnits   *int64  `json:"dynamoDBIndexWriteUnits,omitempty"`
	DynamoDBIndexReadUnits    *int64  `json:"dynamoDBIndexReadUnits,omitempty"`
	AwsAccessKeyID            *string `json:"awsAccessKeyID,omitempty"`
	AwsSecretAccessKey        *string `json:"awsSecretAccessKey,omitempty"`
	AwsRegion                 *string `json:"awsRegion,omitempty"`
	SSLEnabled                *bool   `json:"sslEnabled,omitempty"`
}

func TestCreateDirectoryIfNotExists(t *testing.T) {
	path := "./test/test.db"
	fs := afero.NewMemMapFs()
	fs.Mkdir("./test", os.ModePerm)
	d, err := fs.Stat("./test")
	assert.Nil(t, err)
	assert.NotNil(t, d)

	err = sqlite.CreateDirectoryIfNotExists(path, fs)
	assert.Nil(t, err)
	f, err := fs.Stat("./test")
	assert.Nil(t, err)
	assert.NotNil(t, f)

	path = "./test2/test.db"
	_, err = fs.Stat("./test2")
	assert.NotNil(t, err)
	assert.True(t, os.IsNotExist(err))

	err = sqlite.CreateDirectoryIfNotExists(path, fs)
	assert.Nil(t, err)
	f, err = fs.Stat("./test2")
	assert.Nil(t, err)
	assert.NotNil(t, f)
}
