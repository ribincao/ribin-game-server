package db

var GlobalDB DB

func InitDB() {
	GlobalDB = NewRedisDB()
}
