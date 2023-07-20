package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dbUrl := os.Getenv("POSTGRES_URL")
	//tablePrefix := fmt.Sprintf("%s.", config.DBSchema)
	gormConfig := &gorm.Config{
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix:   tablePrefix, // schema name
		//	SingularTable: false,
		//},
	}

	DB, err = gorm.Open(postgres.Open(dbUrl), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("? Connected Successfully to the Database")
	//DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
}
