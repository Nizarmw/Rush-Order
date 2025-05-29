package service

import (
	"RushOrder/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminSession struct {
	IDPegawai int    `json:"id_pegawai"`
	Nama      string `json:"username"`
}

const (
	AdminSessionName = "admin_session"
	AdminSessionKey  = "admin_id"
	SessionExpiry    = 3600
)

func InitAdminSession(secretKey string) {
	Store = sessions.NewCookieStore([]byte(secretKey))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   SessionExpiry,
		HttpOnly: true,
		Secure:   false,
	}
}

func LoginAdmin(c *gin.Context, username, password string, db *gorm.DB) (*AdminSession, error) {
	var pegawai models.Pegawai
	if err := db.Where("username = ?", username).First(&pegawai).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("username atau password salah")
		}
		return nil, fmt.Errorf("gagal mencari pegawai: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pegawai.Password), []byte(password)); err != nil {
		return nil, errors.New("username atau password salah")
	}

	adminData := &AdminSession{
		IDPegawai: pegawai.IDPegawai,
		Nama:      pegawai.Username,
	}
	sess, err := Store.Get(c.Request, AdminSessionName)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan session: %w", err)
	}

	sess.Values[AdminSessionKey] = adminData.IDPegawai
	sess.Options.MaxAge = SessionExpiry
	if err := sess.Save(c.Request, c.Writer); err != nil {
		return nil, fmt.Errorf("gagal menyimpan session: %w", err)
	}
	return adminData, nil
}

func GetAdminSession(c *gin.Context, db *gorm.DB) (*AdminSession, error) {
	sess, err := Store.Get(c.Request, AdminSessionName)
	if err != nil {
		return nil, fmt.Errorf("gagal mendapatkan session: %w", err)
	}
	idPegawai, ok := sess.Values[AdminSessionKey].(int)
	if !ok {
		return nil, nil
	}

	var pegawai models.Pegawai
	if err := db.First(&pegawai, idPegawai).Error; err != nil {
		return nil, fmt.Errorf("gagal mendapatkan data admin: %w", err)
	}

	return &AdminSession{
		IDPegawai: pegawai.IDPegawai,
		Nama:      pegawai.Username,
	}, nil
}

func LogoutAdmin(c *gin.Context) error {
	sess, err := Store.Get(c.Request, AdminSessionName)
	if err != nil {
		return fmt.Errorf("gagal mendapatkan session: %w", err)
	}
	delete(sess.Values, AdminSessionKey)
	sess.Options.MaxAge = 0
	return sess.Save(c.Request, c.Writer)
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("gagal mengenkripsi password: %w", err)
	}
	return string(hashedBytes), nil
}

func GetOrdersAdmin(db *gorm.DB) ([]models.Order, error) {
	var orders []models.Order
	if err := db.Preload("Items").
		Where("status_admin = ?", models.AdminStatusProcess).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("gagal mendapatkan order dengan status admin process: %v", err)
	}
	return orders, nil
}

// GetAdminOrders gets orders with optional status filter for admin dashboard
func GetAdminOrders(db *gorm.DB, status string) ([]models.Order, error) {
	var orders []models.Order
	query := db.Model(&models.Order{})

	if status != "" {
		if status == models.AdminStatusProcess || status == models.AdminStatusCompleted {
			query = query.Where("status_admin = ? AND status_customer = ?", status, models.CustomerStatusSuccess)
		} else {
			return nil, fmt.Errorf("invalid status filter")
		}
	} else {
		// Get all orders that have been paid (success) regardless of admin status
		query = query.Where("status_customer = ?", models.CustomerStatusSuccess)
	}

	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("gagal mendapatkan orders: %v", err)
	}
	return orders, nil
}

// GetOrderStats gets statistics for admin dashboard
func GetOrderStats(db *gorm.DB) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count pending orders (process status)
	var pendingCount int64
	if err := db.Model(&models.Order{}).
		Where("status_admin = ? AND status_customer = ?", models.AdminStatusProcess, models.CustomerStatusSuccess).
		Count(&pendingCount).Error; err != nil {
		return nil, fmt.Errorf("gagal menghitung pending orders: %v", err)
	}

	// Count completed orders
	var completedCount int64
	if err := db.Model(&models.Order{}).
		Where("status_admin = ?", models.AdminStatusCompleted).
		Count(&completedCount).Error; err != nil {
		return nil, fmt.Errorf("gagal menghitung completed orders: %v", err)
	}

	// Calculate total revenue from completed orders
	var totalRevenue int64
	if err := db.Model(&models.Order{}).
		Where("status_admin = ?", models.AdminStatusCompleted).
		Select("COALESCE(SUM(total_harga), 0)").
		Scan(&totalRevenue).Error; err != nil {
		return nil, fmt.Errorf("gagal menghitung total revenue: %v", err)
	}

	// Total orders (all paid orders)
	var totalOrders int64
	if err := db.Model(&models.Order{}).
		Where("status_customer = ?", models.CustomerStatusSuccess).
		Count(&totalOrders).Error; err != nil {
		return nil, fmt.Errorf("gagal menghitung total orders: %v", err)
	}

	stats["pendingCount"] = pendingCount
	stats["completedCount"] = completedCount
	stats["totalRevenue"] = totalRevenue
	stats["totalOrders"] = totalOrders

	return stats, nil
}
