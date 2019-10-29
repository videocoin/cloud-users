package datastore

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Datastore struct {
	User     *UserDatastore
	Token    *TokenDatastore
	Transfer *TransferDatastore
}

func NewDatastore(uri string) (*Datastore, error) {
	ds := new(Datastore)

	db, err := gorm.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	userDs, err := NewUserDatastore(db)
	if err != nil {
		return nil, err
	}

	ds.User = userDs

	tokenDs, err := NewTokenDatastore(db)
	if err != nil {
		return nil, err
	}

	ds.Token = tokenDs

	transferDs, err := NewTransferDatastore(db)
	if err != nil {
		return nil, err
	}

	ds.Transfer = transferDs

	return ds, nil
}
