package service

import (
	"RushOrder/config"
	"RushOrder/models"
	"errors"

	"gorm.io/gorm"
)

func CreateProduk(produk models.Produk) error {
	if err := config.DB.Create(&produk).Error; err != nil {
		return errors.New("gagal membuat produk: " + err.Error())
	}
	return nil
}

func GetProdukByID(id string) (*models.Produk, error) {
	var produk models.Produk
	err := config.DB.First(&produk, "id_produk = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("produk tidak ditemukan")
	}

	if err != nil {
		return nil, err
	}

	return &produk, nil
}

func GetProduk() ([]models.Produk, error) {
	var produk []models.Produk
	err := config.DB.Find(&produk).Error

	if err != nil {
		return nil, errors.New("gagal mengambil produk: " + err.Error())
	}

	return produk, nil
}

func UpdateProduk(id string, produk models.Produk) error {
	var existingProduk models.Produk
	err := config.DB.Model(&existingProduk).Where("id_produk = ?", id).Updates(map[string]interface{}{
		"nama_produk":  produk.NamaProduk,
		"deskripsi":    produk.Deskripsi,
		"harga_produk": produk.HargaProduk,
		"image_url":    produk.ImageURL,
	}).Error
	if err != nil {
		return errors.New("gagal memperbarui produk: " + err.Error())
	}
	return nil
}

func DeleteProduk(id string) error {
	var produk models.Produk
	err := config.DB.First(&produk, "id_produk = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("produk tidak ditemukan")
	}
	if err != nil {
		return errors.New("gagal menemukan produk: " + err.Error())
	}
	err = config.DB.Delete(&produk, "id_produk = ?", id).Error
	if err != nil {
		return errors.New("gagal menghapus produk: " + err.Error())
	}
	return nil
}

func SearchProduk(keyword string) ([]models.Produk, error) {
	var produk []models.Produk
	err := config.DB.Where("nama_produk LIKE ?", "%"+keyword+"%").Find(&produk).Error
	if err != nil {
		return nil, errors.New("gagal mencari produk: " + err.Error())
	}
	return produk, nil
}
