package controller

import (
	"RushOrder/config"
	"RushOrder/models"
	"RushOrder/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AdminLoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AdminLoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
			return
		}

		adminSession, err := service.LoginAdmin(c, req.Username, req.Password, db)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Login berhasil",
			"admin":   adminSession,
		})
	}
}

func AdminRegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AdminRegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
			return
		}

		hashedPassword, err := service.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal enkripsi password"})
			return
		}

		newAdmin := models.Pegawai{
			Username: req.Username,
			Password: hashedPassword,
		}

		if err := db.Create(&newAdmin).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan admin"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Admin berhasil didaftarkan"})
	}
}

func AdminLogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := service.LogoutAdmin(c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal logout"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
	}
}

func GetOrdersAdminHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		orders, err := service.GetOrdersAdmin(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Berhasil mendapatkan order dengan status admin process",
			"orders":  orders,
		})
	}
}

func GetAdminOrdersHandler(c *gin.Context) {

	status := c.Query("status")

	if status != "" && status != models.AdminStatusProcess && status != models.AdminStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status filter"})
		return
	}

	orders, err := service.GetAdminOrders(config.DB, status)
	if err != nil {
		log.Printf("[GetAdminOrdersHandler] Error getting admin orders: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data order"})
		return
	}

	log.Printf("[GetAdminOrdersHandler] Filter: '%s', Found %d orders.", status, len(orders))
	for _, o := range orders {
		log.Printf("[GetAdminOrdersHandler] Order ID: %s, CustStatus: %s, AdminStatus: %s, CreatedAt: %s",
			o.IDOrder, o.StatusCustomer, o.StatusAdmin, o.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	response := make([]gin.H, len(orders))
	for i, order := range orders {
		items, _ := service.GetOrderItems(order.IDOrder)
		response[i] = gin.H{
			"id_order":        order.IDOrder,
			"id_pemesan":      order.IDPemesan,
			"total_harga":     order.TotalHarga,
			"status_customer": order.StatusCustomer,
			"status_admin":    order.StatusAdmin,
			"items":           items,
			"created_at":      order.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mendapatkan daftar order",
		"orders":  response,
		"count":   len(orders),
	})
}
