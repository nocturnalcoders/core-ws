package kost

import "gorm.io/gorm"

//Campaign -> Kost

type Repository interface {
	FindAll() ([]Kost, error)
	FindByUserID(userID int) ([]Kost, error)
	FindByID(ID int) (Kost, error)
	Save(kost Kost) (Kost, error)
	Update(kost Kost) (Kost, error)
	CreateImage(kostImage KostImage) (KostImage, error)
	MarkAllImagesAsNonPrimary(kostID int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Kost, error) {
	//Mengambil Semua nilai yang ada di DB
	//Slice of = untuk mengambil banyaknya Data
	var kosts []Kost

	err := r.db.Preload("KostImages", "kost_images.is_primary = 1").Find(&kosts).Error
	if err != nil {
		return kosts, err
	}

	return kosts, nil
}

func (r *repository) FindByUserID(userID int) ([]Kost, error) {
	var kosts []Kost
	//preload akan load sebuah relasi kost_images

	err := r.db.Where("user_id = ?", userID).Preload("KostImages", "kost_images.is_primary = 1").Find(&kosts).Error
	// err := r.db.Where("user_id = ?", userID).Preload("KostImages", "kost_images.is_primary = 1").Find(&kosts).Error
	//"KostImages" -> nama field
	//"kost_images.is_primary" -> nama tabelnya
	//"kost_images.is_primary = 1" -> melakukan filter bahwa kost images saat kita melakukan load kost yang dibuat user id , kita skalian mau ambil datanya images , tapi yang diambil hanya is_primary 1
	if err != nil {
		return kosts, err
	}

	return kosts, nil

}

func (r *repository) FindByID(ID int) (Kost, error) {
	var kost Kost
	err := r.db.Preload("User").Preload("KostImages").Where("id = ?", ID).Find(&kost).Error

	if err != nil {
		return kost, err
	}

	return kost, nil
}

func (r *repository) Save(kost Kost) (Kost, error) {
	err := r.db.Create(&kost).Error
	if err != nil {
		return kost, err
	}

	return kost, nil
}

func (r *repository) Update(kost Kost) (Kost, error) {
	err := r.db.Save(&kost).Error

	if err != nil {
		return kost, err
	}

	return kost, nil
}

func (r *repository) CreateImage(kostImage KostImage) (KostImage, error) {
	err := r.db.Create(&kostImage).Error
	if err != nil {
		return kostImage, err
	}

	return kostImage, nil
}

func (r *repository) MarkAllImagesAsNonPrimary(kostID int) (bool, error) {
	//querynya
	//UPDATE kost_images SET is_primary = false WHERE kostID = 1
	//Query Gorm
	err := r.db.Model(&KostImage{}).Where("kost_id = ?", kostID).Update("is_primary", false).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
