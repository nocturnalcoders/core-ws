package main

import (
	"backendEkost/auth"
	"backendEkost/handler"
	"backendEkost/helper"
	"backendEkost/kost"
	"backendEkost/payment"
	"backendEkost/transaction"
	"backendEkost/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// dsn := "root:@tcp(127.0.0.1:3306)/ekostdb?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// dsn := "host=localhost user=postgres password=100299 dbname=ekostdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	dsn := "host=ec2-52-23-40-80.compute-1.amazonaws.com user=jjmonntkyghwzf password=778d68ac3710d2e35cca7456a6b1d1b59a8ff5f89bb1a59ca028ed29d8995683 dbname=de7crutjuikjlf port=5432 sslmode=require TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	kostRepository := kost.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	// kosts, err := kostRepository.FindAll()
	// kosts, err := kostRepository.FindByUserID(1)

	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println(len(kosts))
	// for _, kost := range kosts {
	// 	fmt.Println(kost.Name)
	// 	if len(kost.KostImages) > 0 {
	// 		fmt.Println("jumlah gambar")
	// 		fmt.Println(len(kost.KostImages))
	// 		fmt.Println(kost.KostImages[0].FileName)
	// 	}
	// }

	userService := user.NewService(userRepository)
	kostService := kost.NewService(kostRepository)
	authService := auth.NewService()
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, kostRepository, paymentService)

	// user, _ := userService.GetUserByID(1)

	// input := transaction.CreateTransactionInput{
	// 	KostID: 1,
	// 	Amount: 5000000,
	// 	User:   user,
	// }
	// transactionService.CreateTransaction(input)

	// input := kost.CreateKostInput{}
	// input.Name = "Kost Bu OONG MANTEPPP BEUTT"
	// input.ShortDescription = "Short amat bos"
	// input.Description = "LONGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG"
	// input.LiverCount = 10
	// input.Perks = "3 Kamar Bobo, 2 Kamar Mandi dalam, 3 Tv Room Esklusiv"

	// //Mengambil user dari UserService
	// inputUser, _ := userService.GetUserByID(28)

	// input.User = inputUser

	// _, err = kostService.CreateKost(input)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// kosts, _ := kostService.FindKosts(2)
	// fmt.Println(len(kosts))

	userHandler := handler.NewUserHandler(userService, authService)
	kostHandler := handler.NewKostHandler(kostService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	//Test Validasi Token
	// // token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2fQ.ErbRHHI-DYqCEwjmRfuBa60a40Slygl7jnXYi0Uq3bg")
	// // if err != nil {
	// // 	fmt.Println("EROOR")
	// // 	fmt.Println("EROOR")
	// // 	fmt.Println("EROOR")
	// // }
	// // if token.Valid {
	// // 	fmt.Println("VALID")
	// // 	fmt.Println("VALID")
	// // 	fmt.Println("VALID")
	// // } else {
	// // 	fmt.Println("INVALID")
	// // 	fmt.Println("INVALID")
	// // 	fmt.Println("INVALID")
	// // }
	// fmt.Println(authService.GenerateToken(1001))

	//====================================================
	// Test Service
	// input := user.LoginInput{
	// 	Email:    "christiantolie99@gmail.com",
	// 	Password: "yaudalah",
	// }
	// user, err := userService.Login(input)
	// if err != nil {
	// 	fmt.Println("Terjadi Kesalahan")
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(user.Email)
	// fmt.Println(user.Name)
	//=======================================================
	// userByEmail, err := userRepository.FindByEmail("Christiantolie99@gmail.com")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// if userByEmail.ID == 0 {
	// 	fmt.Println("User tidak ditemukan")
	// } else {
	// 	fmt.Println(userByEmail.Name)
	// }

	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/images", "./images") // "./images" -> nama folder , "/images" -> akses folder
	api := router.Group("/api/v1")       //API Versioning

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	//Perbedaan authMiddleware dengan authMiddleware()
	//authMiddleware kita mempassing middlewarenya
	//authMiddleware() brati yang dipassing nilai kembalian dari eksekui authMiddleware
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	api.GET("/kosts", kostHandler.GetKosts)
	api.GET("/kosts/:id", kostHandler.GetKost)
	api.POST("/kosts", authMiddleware(authService, userService), kostHandler.CreateKost)
	api.PUT("/kosts/:id", authMiddleware(authService, userService), kostHandler.UpdateKost)
	api.POST("/kost-images", authMiddleware(authService, userService), kostHandler.UploadImage)

	api.GET("/kosts/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetKostTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)
	//Cara mendapatkan token dengan mengambil dari token
	//Login email dan password lalu send lewat body dengan POST

	router.Run()

	// userInput := user.RegisterUserInput{}
	// userInput.Name = "test simpan dari service"
	// userInput.Email = "contoh@gmail.com"
	// userInput.Occupation = "anak band"
	// userInput.Password = "ini password"

	// userService.RegisterUser(userInput)

	// user := user.User{
	// 	Name: "test simpen",
	// }

	// userRepository.Save(user)
	// fmt.Println("Database Connected")

	// var users []user.User
	// length := len(users)

	// fmt.Println(length)

	// db.Find(&users)

	// length = len(users)
	// fmt.Println(length)

	// for _, user := range users {
	// 	fmt.Println(user.Name)
	// 	fmt.Println(user.Email)
	// 	fmt.Println("========")
	// }

	// router := gin.Default()
	// router.GET("/handler", handler)
	// router.Run()
}

// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// func handler(c *gin.Context) {
// 	dsn := "root:@tcp(127.0.0.1:3306)/ekostdb?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	var users []user.User
// 	db.Find(&users)

// 	c.JSON(http.StatusOK, users)
// }

//Jadi Method Authmiddleware akan mengembalikan sebuah function handler func
//Handlerfunc adalah func yang parameternya gin context
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//mengambil key yang ada di token
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)

	}

}

//Tidak boleh memakai tambahan parameter untuk gin.context karena nanti sudah bukan sebuah middleware lagi
// func authMiddleware(c *gin.Context) {
// 	authHeader := c.GetHeader("Authorization")

// 	if !strings.Contains(authHeader, "Bearer") {
// 		response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
// 		return
// 		//Mengapa memakai AbortWithStatusJSON , karena dia adalah sebuah middleware supaya tidak berlanjut ke proses berikutnya , atau dihentikan oleh sistem
// 	}

// 	//Bearer tokentokentoken -> dipisah oleh spasi
// 	//Isinya sekarang adalah 2 buah array of strings
// 	// var tokenString string
// 	tokenString := ""
// 	arrayToken := strings.Split(authHeader, " ")
// 	if len(tokenString) == 2 {
// 		tokenString = arrayToken[1]
// 	}

// 	//tahap validasi token
// 	token, err :=
// }
