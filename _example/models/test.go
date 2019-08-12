package models

import (
	"github.com/pefish/go-error"
	"github.com/pefish/go-mysql"
)

type TestModel struct {
	UserId uint64 `db:"user_id" json:"user_id"`
	Name   string `db:"name" json:"name"`

	BaseModel
}

func (this *TestModel) GetTableName() string {
	return `test`
}

func (this *TestModel) GetNameByUid(userId string) string {
	testModel := TestModel{}
	if notFound := go_mysql.MysqlHelper.SelectFirst(&testModel, testModel.GetTableName(),`*`, map[string]interface{}{
		`user_id`: userId,
	}); notFound {
		go_error.Throw(`user not found`, 20000)
	}
	return testModel.Name
}
