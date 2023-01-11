package db

import (
	"github.com/olivere/elastic/v7"
	"github.com/ribincao/ribin-game-server/config"
	"github.com/ribincao/ribin-game-server/logger"
)

var Client *elastic.Client

const IdxEngineMatch = "engine_match"

func InitES() {
	var err error
	Client, err = elastic.NewClient(
		elastic.SetURL(config.GlobalConfig.DbConfig.EsUrl),
		elastic.SetBasicAuth(config.GlobalConfig.DbConfig.EsUsername, config.GlobalConfig.DbConfig.EsPasswd),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	logger.Info("[Engine-Tool] ElasticSearch Client Initialized!")
}
