package migrate

import (
	"main/src/initializers"
	"main/src/models"
)

func init() {
	initializers.ConnectDB()
}

func RunMigrations() {
	initializers.DB.AutoMigrate(&models.ChargerEntity{})

}

//func main() {
//	initializers.DB.AutoMigrate(&models.User{})
//	fmt.Println("? Migration complete")
//}
