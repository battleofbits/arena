package arena

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
)

func checkError(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func GetConnection() *sql.DB {
	// XXX
	//config, err := toml.LoadFile("config.toml")
	//checkError(err)
	//configTree := config.Get("postgres").(*toml.TomlTree)
	//user := configTree.Get("user").(string)
	//password := configTree.Get("password").(string)
	//host := configTree.Get("host").(string)
	//database := configTree.Get("database").(string)
	//port := configTree.Get("port").(int64)

	rawUrl := "postgres://postgres_arena@localhost:5432/arena?sslmode=disable"
	url, err := pq.ParseURL(rawUrl)
	checkError(err)
	db, err := sql.Open("postgres", url)
	checkError(err)
	return db
}
