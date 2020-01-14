package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"todoapi/config"

	"github.com/jinzhu/gorm"
)

type (
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}

	transformedTodo struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

func createdTodo(c *gin.Context) {

	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
	config.GetDB().Save(&todo)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated, "message": "Todo Created!", "resourceId": todo.ID,
	})
}

func getTodo(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTodo

	config.GetDB().Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusCreated, gin.H{
			"status": http.StatusCreated, "message": "Todo not Found",
		})
		return
	}

	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": _todos,
	})

}

func deleteTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	config.GetDB().First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "No Todo Found",
		})
		return
	}

	config.GetDB().Delete(&todo)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "todo deleted Succes",
	})

}

func updateTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	config.GetDB().First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "todo not found",
		})
		return
	}

	config.GetDB().Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	config.GetDB().Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "todo Updated!",
	})
}

func main() {

	router := gin.Default()
	router.GET("/todo", getTodo)
	router.POST("/todo", createdTodo)
	router.DELETE("/todo/:id", deleteTodo)
	router.PATCH("/todo/:id", updateTodo)

	router.Run()

}
