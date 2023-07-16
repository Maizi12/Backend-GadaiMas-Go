package config

import(
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"log"
	"assignment/models"
)

var DB *gorm.DB

func ConnectDB(){

	db,err:=gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/leasing?parseTime=true"),&gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Invoice{})
	db.AutoMigrate(&models.Car{})
	db.AutoMigrate(&models.Leasing{})
	db.AutoMigrate(&models.Payment{})

	DB=db
	log.Println("Database connected")
}