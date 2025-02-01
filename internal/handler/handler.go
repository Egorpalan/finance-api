package handler

import (
	"github.com/Egorpalan/finance-api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) TopUpBalance(c *gin.Context) {
	var request struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.TopUpBalance(request.UserID, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance topped up successfully"})
}

func (h *Handler) TransferMoney(c *gin.Context) {
	var request struct {
		SenderID   int64   `json:"sender_id"`
		ReceiverID int64   `json:"receiver_id"`
		Amount     float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.TransferMoney(request.SenderID, request.ReceiverID, request.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transfer successful"})
}

func (h *Handler) GetTransactions(c *gin.Context) {
	userID := c.Param("user_id")

	userIDInt64, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	transactions, err := h.service.GetTransactions(userIDInt64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *Handler) CreateUserHandler(c *gin.Context) {
	var request struct {
		Balance float64 `json:"balance"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.CreateUser(request.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      user.ID,
		"balance": user.Balance,
	})
}
