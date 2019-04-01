package models

type TestModel struct {
	UserId uint64 `db:"user_id" json:"user_id"`
	Name   string `db:"name" json:"name"`

	BaseModel
}

func (this *TestModel) GetTableName() string {
	return `test`
}
