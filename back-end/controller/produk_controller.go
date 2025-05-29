package controller

import (
	"RushOrder/models"
	"RushOrder/service"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateProdukHandler(c *gin.Context) {
	idProduk := c.PostForm("id_produk")
	namaProduk := c.PostForm("nama_produk")
	deskripsi := c.PostForm("deskripsi")
	hargaProdukStr := c.PostForm("harga_produk")
	kategoriStr := strings.ToLower(c.PostForm("kategori"))
	fmt.Printf("Kategori yang diterima: '%s'\n", kategoriStr)

	kategori := models.KategoriProduk(kategoriStr)
	fmt.Printf("Kategori setelah konversi: '%s'\n", string(kategori))
	fmt.Printf("IsValid result: %t\n", kategori.IsValid())

	if !kategori.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "Kategori tidak valid, pilih: Makanan, Minuman, atau Snack",
			"received_category": kategoriStr,
			"valid_categories":  models.GetValidKategoriProduk(),
		})
		return
	}

	if idProduk == "" || namaProduk == "" || deskripsi == "" || hargaProdukStr == "" || kategoriStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua field wajib diisi"})
		return
	}

	if !kategori.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":            "Kategori tidak valid, pilih: Makanan, Minuman, atau Snack",
			"valid_categories": models.GetValidKategoriProduk(),
		})
		return
	}

	hargaProduk, err := strconv.Atoi(hargaProdukStr)
	if err != nil || hargaProduk <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Harga produk harus berupa angka positif"})
		return
	}

	file, err := c.FormFile("image_url")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal mengambil file gambar"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format gambar harus jpg, jpeg, atau png"})
		return
	}

	supa := service.NewSupabaseStorage()
	imageURL, err := supa.Upload(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupload gambar: " + err.Error(),
		})
		return
	}
	produk := models.Produk{
		IDProduk:    idProduk,
		NamaProduk:  namaProduk,
		Deskripsi:   deskripsi,
		HargaProduk: hargaProduk,
		ImageURL:    imageURL,
		Kategori:    kategori,
	}

	if err := service.CreateProduk(produk); err != nil {
		_ = supa.Delete(imageURL)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal membuat produk: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil dibuat",
		"produk":  produk})
}

func GetProduk(c *gin.Context) {
	produk, err := service.GetProduk()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil produk: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil diambil",
		"produk":  produk,
		"count":   len(produk),
	})
}

func GetProdukByID(c *gin.Context) {
	id := c.Param("id")

	produk, err := service.GetProdukByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk ditemukan",
		"data":    produk,
	})
}

func GetProdukByKategori(c *gin.Context) {
	kategoriParam := c.Param("kategori")
	kategori := models.KategoriProduk(kategoriParam)

	if !kategori.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":            "Kategori tidak valid, pilih: Makanan, Minuman, atau Snack",
			"valid_categories": models.GetValidKategoriProduk(),
		})
		return
	}

	produk, err := service.GetProdukByKategori(kategori)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Produk berhasil diambil berdasarkan kategori",
		"kategori": kategori,
		"produk":   produk,
		"count":    len(produk),
	})
}

func UpdateProduk(c *gin.Context) {
	id := c.Param("id")

	existingProduk, err := service.GetProdukByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	namaProduk := c.PostForm("nama_produk")
	deskripsi := c.PostForm("deskripsi")
	hargaProdukStr := c.PostForm("harga_produk")

	hargaProduk, err := strconv.Atoi(hargaProdukStr)
	if err != nil || hargaProduk <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Harga produk harus berupa angka positif",
		})
		return
	}

	updatedProduk := models.Produk{
		NamaProduk:  namaProduk,
		Deskripsi:   deskripsi,
		HargaProduk: hargaProduk,
		ImageURL:    existingProduk.ImageURL,
	}

	file, err := c.FormFile("image")
	if err != nil {
		log.Println("Tidak ada file gambar baru, menggunakan gambar lama")
	} else {

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "File harus berupa gambar (.jpg, .jpeg, .png, .webp)",
			})
			return
		}

		supa := service.NewSupabaseStorage()

		if existingProduk.ImageURL != "" {
			_ = supa.Delete(existingProduk.ImageURL)
		}

		imageURL, err := supa.Upload(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal mengupload gambar: " + err.Error(),
			})
			return
		}

		updatedProduk.ImageURL = imageURL
	}

	if err := service.UpdateProduk(id, updatedProduk); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	finalProduk, _ := service.GetProdukByID(id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil diupdate",
		"data":    finalProduk,
	})
}

func DeleteProduk(c *gin.Context) {
	id := c.Param("id")

	produk, err := service.GetProdukByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if produk.ImageURL != "" {
		supa := service.NewSupabaseStorage()
		if err := supa.Delete(produk.ImageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal menghapus gambar: " + err.Error(),
			})
			return
		}
	}

	if err := service.DeleteProduk(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus produk: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil dihapus",
		"data":    produk,
	})
}

func SearchProduk(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query tidak boleh kosong"})
		return
	}

	produk, err := service.SearchProduk(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencari produk: " + err.Error()})
		return
	}

	if len(produk) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Produk tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk ditemukan",
		"produk":  produk,
		"count":   len(produk),
		"query":   query,
	})
}
