package kost

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetKosts(userID int) ([]Kost, error)
	GetKostByID(input GetKostDetailInput) (Kost, error)
	CreateKost(input CreateKostInput) (Kost, error)
	UpdateKost(inputID GetKostDetailInput, inputData CreateKostInput) (Kost, error)
	SaveKostImage(input CreateKostImageInput, fileLocation string) (KostImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetKosts(userID int) ([]Kost, error) {
	//mengapa tidak memakai json karena parameter yang dikirimkan oleh user akan lsg mendapatkan integer
	//nanti di cek apakah userID ada atau tidak
	//hanya mengambil userID yang bersangkutan
	//jika kosong ,kita akan menampilkan data kosts
	if userID != 0 {
		kosts, err := s.repository.FindByUserID(userID)
		if err != nil {
			return kosts, err
		}

		return kosts, nil
	}

	kosts, err := s.repository.FindAll()
	if err != nil {
		return kosts, err
	}

	return kosts, nil

}

func (s *service) GetKostByID(input GetKostDetailInput) (Kost, error) {
	kost, err := s.repository.FindByID(input.ID)
	if err != nil {
		return kost, err
	}

	return kost, nil
}

func (s *service) CreateKost(input CreateKostInput) (Kost, error) {
	//Melakukan mapping dari inputan user ke createkostinput
	//Kemudian dari CreateCampaignInput menjadi object Kost
	kost := Kost{}
	kost.Name = input.Name
	kost.ShortDescription = input.ShortDescription
	kost.Description = input.Description
	kost.Perks = input.Perks
	kost.Latitude = input.Latitude
	kost.Longitude = input.Longitude
	kost.Price = input.Price
	kost.UserID = input.User.ID
	//Nama kost + id user ke brp -> nama-kost-10
	//Membuat variabel untuk membuat gabungan antara nama kost dan id

	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	kost.Slug = slug.Make(slugCandidate)

	//Proses pembuatan Slug
	//Memanggil Repository

	newKost, err := s.repository.Save(kost)
	if err != nil {
		return newKost, err
	}

	return newKost, nil
}

func (s *service) UpdateKost(inputID GetKostDetailInput, inputData CreateKostInput) (Kost, error) {
	kost, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return kost, err
	}

	//Check user yang memiliki kost
	//nanti tambahin role
	if kost.UserID != inputData.User.ID {
		return kost, errors.New("Not an owner of Kost")
	}

	kost.Name = inputData.Name
	kost.ShortDescription = inputData.ShortDescription
	kost.Description = inputData.Description
	kost.Perks = inputData.Perks
	kost.Latitude = inputData.Latitude
	kost.Longitude = inputData.Longitude
	kost.Price = inputData.Price

	updatedKost, err := s.repository.Update(kost)
	if err != nil {
		return updatedKost, err
	}

	return updatedKost, nil
}

func (s *service) SaveKostImage(input CreateKostImageInput, fileLocation string) (KostImage, error) {
	kost, err := s.repository.FindByID(input.KostID)
	if err != nil {
		return KostImage{}, err
	}

	if kost.UserID != input.User.ID {
		return KostImage{}, errors.New("Not an owner of Kost")
	}

	IsPrimary := 0
	if input.IsPrimary {
		IsPrimary = 1

		_, err := s.repository.MarkAllImagesAsNonPrimary(input.KostID)
		if err != nil {
			return KostImage{}, err
		}
	}

	kostImage := KostImage{}
	kostImage.KostID = input.KostID
	kostImage.IsPrimary = IsPrimary
	kostImage.FileName = fileLocation

	newKostImage, err := s.repository.CreateImage(kostImage)
	if err != nil {
		return newKostImage, err
	}

	return newKostImage, nil
}
