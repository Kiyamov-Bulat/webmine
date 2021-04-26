package models

type Validator interface {
	IsValid() bool
}

type MineModel struct {
	ID uint `json:"id" gorm:"primary_key"`
}

func (model *MineModel) IsValid() bool {
	return model.ID != 0
}

func setItemFromDB(tableName string, item Validator, id uint) {
	GetDB().Table(tableName).Where("id = ?", id).Find(item)
}

func getAllItems(tableName string, i interface{}) {
	GetDB().Table(tableName).Find(i)
}
