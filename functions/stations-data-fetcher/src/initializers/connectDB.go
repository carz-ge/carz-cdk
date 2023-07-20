package initializers

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	//tablePrefix := fmt.Sprintf("%s.", config.DBSchema)
	gormConfig := &gorm.Config{
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix:   tablePrefix, // schema name
		//	SingularTable: false,
		//},
	}

	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("? Connected Successfully to the Database")
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
}
