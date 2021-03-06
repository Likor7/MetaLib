package models

import (
	"MetaLib/utils"
	"errors"
)

type Author struct {
	Id   uint
	Name string
}

func GetAuthorById(id uint) (*Author, error) {
	var author Author
	notExist := utils.DB.First(&author, id).RecordNotFound()
	if notExist {
		return nil, errors.New("author not exist")
	} else {
		return &author, nil
	}
}
