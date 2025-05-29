package controller

import (
	"RushOrder/service"
	"RushOrder/session"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddCartRequest struct {
	IDProduk   string `json:"id_produk" binding:"required"`
	NamaProduk string `json:"nama_produk" binding:"required"`
	Jumlah     int    `json:"jumlah" binding:"required,gt=0"`
	Harga      int    `json:"harga" binding:"required,gt=0"`
}

func AddToCartHandler(c *gin.Context) {
	var req AddCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	item := session.CartItem{
		IDProduk:   req.IDProduk,
		NamaProduk: req.NamaProduk,
		Jumlah:     req.Jumlah,
		Harga:      req.Harga,
	}

	err := service.AddToCart(c.Writer, c.Request, item)
	if err != nil {
		// Tambahkan log untuk membantu debugging
		log.Println("AddToCart error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal tambah ke cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item ditambahkan ke cart"})
}

func GetCartHandler(c *gin.Context) {
	cart, total, err := service.GetCart(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal ambil cart"})
		return
	}

	response := struct {
		Items map[string]session.CartItem `json:"items"`
		Total int                         `json:"total"`
	}{
		Items: cart,
		Total: total,
	}

	c.JSON(http.StatusOK, response)
}

func RemoveFromCartHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_produk wajib diisi"})
		return
	}

	err := service.RemoveFromCart(c.Writer, c.Request, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal hapus item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item dihapus dari cart"})
}

func ClearCartHandler(c *gin.Context) {
	err := service.ClearCart(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal kosongkan cart"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cart dikosongkan"})
}

type UpdateCartItemRequest struct {
	IDProduk string `json:"id_produk" binding:"required"`
	Jumlah   int    `json:"jumlah" binding:"required,gt=0"`
}

func UpdateCartItemHandler(c *gin.Context) {
	var req UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := service.UpdateCartItemHandler(c.Writer, c.Request, req.IDProduk, req.Jumlah)
	if err != nil {
		if err.Error() == "session tidak ditemukan" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session tidak ditemukan"})
			return
		}
		if err.Error() == "produk tidak ditemukan di cart" {
			c.JSON(http.StatusNotFound, gin.H{"error": "produk tidak ditemukan di cart"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal update item"})
		return
	}
	if req.Jumlah <= 0 {
		c.JSON(http.StatusOK, gin.H{"message": "item dihapus dari cart"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "item diperbarui di cart"})
	}
}

// CheckoutCartHandler handles the checkout process to convert cart to order
func CheckoutCartHandler(c *gin.Context) {
	orderID, err := service.CheckoutCart(c.Writer, c.Request)
	if err != nil {
		switch err.Error() {
		case "session tidak ditemukan":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session tidak ditemukan"})
		case "keranjang kosong":
			c.JSON(http.StatusBadRequest, gin.H{"error": "keranjang kosong"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal checkout"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "checkout berhasil",
		"order_id": orderID,
	})
}
