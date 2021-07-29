package handler

import (
	"backendEkost/auth"
	"backendEkost/helper"
	"backendEkost/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	//tangkap input dari user
	//map input dari user ke struct RegisterUserInput
	//struct di atas kita passing sebagai parameter service

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors} //gin.h adalah map dimana key nya adalah string dan valuenya adalah interface atau bisa apa saja

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	//user memasukan input {email dan password}
	//input ditangkap handler -> perlu melakukan mapping
	//mapping dari input user ke input struct
	//input struct passing service
	//di service mencari dengan bantuan repository user dengan email x
	//mencocokan password

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, token)

	response := helper.APIResponse("Successfuly Loggedin", http.StatusOK, "success", formatter)
	c.JSON(http.StatusUnprocessableEntity, response)
	return
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	// ada input email dari user
	// input email di-mapping ke struct input
	// struct input di-parsing ke service
	// service akan manggil repository - untuk menentukan apakah email sudah ada atau belum
	// repostory akan melakukan query ke db

	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email checking Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		response := helper.APIResponse("Email checking Failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is Available"
	}

	// if isEmailAvailable {
	// 	metaMessage = "Email is Available"
	// } else {
	// 	metaMessage = "Email has been Registered"
	// }

	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	//tangkap input dari user
	//simpan gambarnya di folder "images/"
	//diservice kita panggil repo
	//JWT (sementara hardcode, seakan-akan user yang login ID = 1)
	//Repo ambil data user yang ID = 1
	//Repo update data user simpan lokasi file
	//menangkap http form data
	//tanpa shouldbyjson

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to Upload Avatar Image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	//old : images/namaFile.png
	//new : images/userid-namafile.png -> images/1-namafile.png
	// path := "images/" + file.Filename
	//tanda %d menggantikan user id , dan %s menggantikan filename

	//harus nya dapat dari JWT, Tapi sabarr
	//Cara mendapatkan data user
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to Upload Avatar Image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to Upload Avatar Image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar Successfuly Uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	formatter := user.FormatUser(currentUser, "")

	response := helper.APIResponse("Successfuly Fetch user data", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}
