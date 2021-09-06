package config

type DB struct {
	MYSQL_DSN string
}

var db_configs = Configs{
	"DB_CONNECTION_STRING": Consumer_Set(func(i interface{}) *string { return &i.(*DB).MYSQL_DSN }),
}

func (db *DB) Configs() Configs { return db_configs }
