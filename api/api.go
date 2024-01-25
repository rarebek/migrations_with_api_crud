package main

import (
	"fmt"
	"net/http"
	"strconv"
	"test_migration/models"
	"test_migration/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/create", CreateUser)
	router.PUT("/update/:id", UpdateUser)
	router.GET("/getone/:id", GetOne)
	router.GET("/getall", GetAll)
	router.DELETE("/delete/:id", DeleteUser)

	router.Run(":7777")
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	respAns, err := storage.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, respAns)
}

func UpdateUser(c *gin.Context) {
	uuid := c.Param("id")
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	respAns, err := storage.UpdateUser(uuid, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, respAns)
}

func GetOne(c *gin.Context) {
	uuid := c.Param("id")
	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UUID is required"})
		return
	}

	respAns, err := storage.GetOne(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, respAns)
}

func GetAll(c *gin.Context) {
	tlimit := c.Query("limit")
	limit, err := strconv.Atoi(tlimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	tpage := c.Query("page")
	page, err := strconv.Atoi(tpage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	respAns, err := storage.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, respAns)
}

func DeleteUser(c *gin.Context) {
	tempId := c.Param("id")
	id, err := strconv.Atoi(tempId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	respAns, err := storage.DeleteUser(id)
	if err != nil {
		fmt.Printf("Error deleting user: %s\n", err)
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, respAns)
}
