package main

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

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

	c.Request.URL.Path = "/todo"
}

func getTodo(c *gin.Context) {

	var todos []todoModel
	config.GetDB().Find(&todos)

	html := ""

	for _, item := range todos {

		if item.Completed == 1 {
			html += `<li><div class="view">
                <input class="toggle" type="checkbox" id="check" onclick="goCheck(` + cast.ToString(item.ID) + `)"  checked >
                <label>` + cast.ToString(item.Title) + `</label>
                    <a class="destroy" href="http://localhost:8080/todo/` + cast.ToString(item.ID) + `"></a>
                    <!-- <button class="destroy" onclick="deleted()"></button> -->
            </div></li>
       `
		} else {
			html += `<li><div class="view">
			<input class="toggle" type="checkbox" id="check" onclick="goCheck(` + cast.ToString(item.ID) + `)" >
			<label>` + cast.ToString(item.Title) + `</label>
				<a class="destroy" href="http://localhost:8080/todo/` + cast.ToString(item.ID) + `"></a>
				<!-- <button class="destroy" onclick="deleted()"></button> -->
		</div></li>`

		}

	}
	gin.Default().LoadHTMLGlob("template/*")

	c.HTML(http.StatusOK, "Index", gin.H{
		"title": "Todo List",
		"data":  template.HTML(html),
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
	c.Redirect(http.StatusMovedPermanently, "http://localhost:8080/todo")

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
func updateCheck(c *gin.Context) {
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
	completed := 0
	if todo.Completed == 0 {
		completed = 1
	} else {
		completed = 0
	}
	config.GetDB().Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "todo Updated!",
	})
}

func main() {

	router := gin.Default()
	router.LoadHTMLGlob("template/*")
	router.GET("/todo", getTodo)
	router.POST("/todo", createdTodo)
	router.GET("/todo/:id", deleteTodo)
	router.PATCH("/todo/:id", updateTodo)
	router.PATCH("/todos/:id", updateCheck)

	router.StaticFS("/css", http.Dir("assets"))
	router.Run()

}
