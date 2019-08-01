package dao

import (
	"github.com/pefish/go-error"
	"test/src/global"
	"test/src/models"
)

type TestDaoClass struct {

}

var TestDao = TestDaoClass{}

func (this *TestDaoClass) GetNameByUid(userId string) string {
	testModel := models.TestModel{}
	if notFound := global.MysqlHelper.SelectFirst(&testModel, testModel.GetTableName(),`*`, map[string]interface{}{
		`user_id`: userId,
	}); notFound {
		go_error.Throw(`user not found`, 20000)
	}
	return testModel.Name
}
