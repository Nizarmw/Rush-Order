package service

import (
	"RushOrder/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminSession struct {
	IDPegawai int    `json:"id_pegawai"`
	Nama      string `json:"nama"`
}

const (
	AdminSessionName = "admin_session"
	AdminSessionKey  = "admin_id"
	SessionExpiry    = 3600
)

func InitAdminSession() {
	Store.Options.MaxAge = SessionExpiry
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
