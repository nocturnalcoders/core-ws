package handler

import (
	"backendEkost/helper"
	"backendEkost/kost"
	"backendEkost/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 1. Tangkap Parameter di handler
// 2. Handler ke service (Diservice yang menentukan apakah repository yang di call)
// 3. Ke Repository : FindAll , FindByUserID
// 4. Ke DB

//GET : api/v1/kosts?user_id=10
//Akan mengambil user id yang ber ID 10
//GET : api/v1/kosts
//Akan mengambil semua data kosts

type kostHandler struct {
	service kost.Service
}

func NewKostHandler(service kost.Service) *kostHandler {
	return &kostHandler{service}
}

//api/v1/kosts
func (h *kostHandler) GetKosts(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	//Kost lebih dari 1 , kumpulan dari object kost
	kosts, err := h.service.GetKosts(userID)
	if err != nil {
		response := helper.APIResponse("Error to Get Kosts", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of Kosts", http.StatusOK, "success", kost.FormatKosts(kosts))
	c.JSON(http.StatusOK, response)

}

func (h *kostHandler) GetKost(c *gin.Context) {
	// bentuknya api/v1/kosts/id -> id bisa 1 / 2 / 3 brapapun
	//Handler -> Mapping id yang di URL ke Struct Input utk dimasukan ke Service , Call formatter
	//Service -> Inputan Struct untuk menangkap ID di URL , pakai shouldbindJSOn
	//Repository  -> untuk get kost by id

	var input kost.GetKostDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to Get Detail of Kost", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// kost, err := h.service.GetKostByID(input)
	// if err != nil {
	// 	response := helper.APIResponse("Failed to Get Detail of Kost", http.StatusBadRequest, "error", nil)
	// 	c.JSON(http.StatusBadRequest, response)
	// 	return
	// }

	kostDetail, err := h.service.GetKostByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to Get Detail of Kost", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Kost Detail", http.StatusOK, "success", kost.FormatKostDetail(kostDetail))
	c.JSON(http.StatusOK, response)
}

func (h *kostHandler) CreateKost(c *gin.Context) {
	var input kost.CreateKostInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed To Create Kost", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	input.User = currentUser

	newKost, err := h.service.CreateKost(input)
	if err != nil {
		response := helper.APIResponse("Failed To Create Kost", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success To Create Kost", http.StatusOK, "success", kost.FormatKost(newKost))
	c.JSON(http.StatusOK, response)
}

func (h *kostHandler) UpdateKost(c *gin.Context) {
	var inputID kost.GetKostDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to Update Kost", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData kost.CreateKostInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed To Update Kost", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	inputData.User = currentUser

	updatedKost, err := h.service.UpdateKost(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to Update Kost", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success To Update Kost", http.StatusOK, "success", kost.FormatKost(updatedKost))
	c.JSON(http.StatusOK, response)
}

func (h *kostHandler) UploadImage(c *gin.Context) {
	var input kost.CreateKostImageInput

	err := c.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed To Upload Kost Image", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to Upload Kost Image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to Upload Kost Image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.SaveKostImage(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to Upload Kost Image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Kost Image Successfuly Uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}

//User mengisi form
//Tangkap parameter dari user ke input struct
//ambil current user dari JWT/Handler
//Panggil service, parameter si input struct (dan juga buat slug)
//panggil repository untuk simpan data kost baru
//field id user, saat mengirim request handler tau bahwa user x yang membuat

//Update Kost
//1. Repository update data kost
//2. Service perlu parameter (Find Kost by ID, tangkap parameter, nilainya kita kasih ke struct kost )
//3. input dari user yang ada di URI (dipassing ke service layer)
//4. Mapping dari input ke input struct (Ada 2 1. dari user, dan 1 dari uri)
//5. Handler
//6. User masukan input

//Upload Images
//Hanlder :
//Tangkap input dan ubah ke struct input
//save image kost ke suatu folder
//service  (Kondisi memanggil point 2 di repo, panggil repo point 1)
//Repository :
//1. Create images/save data image ke dalam table kost_images
//2. ubah is_primary true ke false
//nomor 2 maksudnya adalah : saat user upload gambar, dia akan bisa masukan gambar untuk primary atau tidak .. kareana 1 kost, hanya perlu 1 gambar yang primary ..
