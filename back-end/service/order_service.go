package service

import (
	"RushOrder/config"
	"RushOrder/models"
)

func GetOrderItems(orderID string) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := config.DB.Where("id_order = ?", orderID).Find(&items).Error
	return items, err
}
