package main

import (
	"ShopAPI/internal"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {
	internal.InitDB()
	router := gin.Default()
	router.GET("/client", internal.GetClientsByFullName)

	request, _ := http.NewRequest("GET", "/client", nil)
	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Errorf("Status code is wrong. Expected %d, got %d", http.StatusOK, writer.Code)
	}
}
