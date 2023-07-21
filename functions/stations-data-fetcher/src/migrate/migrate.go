package migrate

import (
	"main/src/initializers"
	"main/src/models"
)

func RunMigrations() {
	initializers.DB.AutoMigrate(&models.AutoStationEntity{})

}

//func main() {
//	initializers.DB.AutoMigrate(&models.User{})
//	fmt.Println("? Migration complete")
//}
