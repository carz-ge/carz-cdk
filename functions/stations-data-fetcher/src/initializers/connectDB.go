package initializers

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB(config Config) {
	var err error
	//tablePrefix := fmt.Sprintf("%s.", config.DBSchema)
	gormConfig := &gorm.Config{
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix:   tablePrefix, // schema name
		//	SingularTable: false,
		//},
	}

	url := config.DBUrl
	log.Println(url)
	DB, err = gorm.Open(postgres.Open(url), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("? Connected Successfully to the Database")
	//DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
}
