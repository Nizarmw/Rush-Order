package controller

import (
	"RushOrder/service"
	"RushOrder/session"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Nama string `json:"nama" binding:"required"`
	Meja int    `json:"meja" binding:"required,gt=0"`
}

func CustomerLoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nama dan meja wajib diisi"})
		return
	}

	id := uuid.New().String()
	customerData := session.CustomerSession{
		ID:    id,
		Nama:  req.Nama,
		Meja:  req.Meja,
		Cart:  make(map[string]session.CartItem),
		Total: 0,
	}

	err := service.CreateSession(c.Writer, c.Request, customerData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat sesi"})
		return
	}

	c.JSON(http.StatusOK, customerData)
}

func GetCustomerSessionHandler(c *gin.Context) {
	customer, err := service.GetSession(c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal ambil sesi"})
		return
	}
	if customer == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func LogoutHandler(c *gin.Context) {
	err := service.ClearCustomerSession(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal logout"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logout berhasil"})
}
