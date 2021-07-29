package kost

import "backendEkost/user"

//tidak memakai json , karena dikirim lewat body biasanya POST
//Ada beberapa mengirim JSON
//Dengan bentuk body, Query Param, Menyatu dengan URL -> Memakai URI
type GetKostDetailInput struct {
	ID int `uri:"id" binding:"required"`
}

//Finding Target -> goal amount diganti liver count
//What will liver get -> what will funders get di jsona adalah Perks
type CreateKostInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	Latitude         string `json:"Latitude" binding:"required"`
	Longitude        string `json:"Longitude" binding:"required"`
	Price            string `json:"Price" binding:"required"`
	Perks            string `json:"Perks" binding:"required"`
	//Mengambil user dari JWT
	User user.User
	//User akan diambil untuk pembuatan slug
}

type CreateKostImageInput struct {
	KostID    int  `form:"kost_id" binding:"required"`
	IsPrimary bool `form:"is_primary"`
	User      user.User
}
