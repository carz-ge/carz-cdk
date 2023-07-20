package migrate

import (
	"log"
	"main/src/initializers"
	"main/src/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func RunMigrations() {
	initializers.DB.AutoMigrate(&models.ChargerEntity{})

}

//func main() {
//	initializers.DB.AutoMigrate(&models.User{})
//	fmt.Println("? Migration complete")
//}
