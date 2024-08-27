package cmd

import (
	"fmt"

	"godoai/knowledge/chromemDB"
	"godoai/knowledge/dbtype"
)

func NewDbFactory(config dbtype.DbConfig) (*dbtype.DbFactory, error) {
	var db dbtype.Database
	var err error
	switch config.Type {
	case "chromem":
		db, err = chromemDB.NewDB(config)
	default:
		return nil, fmt.Errorf("unknown database type: %s", config.Type)
	}

	if err != nil {
		return nil, err
	}

	return &dbtype.DbFactory{DB: db}, nil
}
