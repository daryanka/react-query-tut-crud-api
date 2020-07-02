package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"strings"
)

type Post struct {
	ID   string `json:"id"`
	Name string `json:"name" validate:"required"`
	Body string `json:"body" validate:"required"`
}

type DefaultError struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Field string `json:"field,omitempty"`
}

var Posts = []Post{}

func main() {
	r := gin.Default()
	r.RedirectTrailingSlash = true
	r.Use(cors.Default())

	postRoutes := r.Group("/posts")
	{
		postRoutes.GET("/", Get)
		postRoutes.GET("/:id", GetIndividual)
		postRoutes.POST("/", Create)
		postRoutes.PUT("/:id", Update)
		postRoutes.DELETE("/:id", Delete)
	}

	if err := r.Run(":5050"); err != nil {
		panic("unable to run on port 5050")
	}
}

func Get(c *gin.Context) {
	c.JSON(200, Posts)
}

func GetIndividual(c *gin.Context) {
	id := c.Param("id")
	for _, el := range Posts {
		if el.ID == id {
			c.JSON(200, el)
			return
		}
	}

	c.JSON(404, gin.H{
		"error": true,
		"message": "post not found",
	})
}

func Update(c *gin.Context) {
	var reqBody Post
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, ErrorMsg("invalid JSON body"))
		return
	}

	errs := Validate(reqBody)
	if errs != nil {
		c.JSON(422,errs)
		return
	}

	id := c.Param("id")
	for i, el := range Posts {
		if el.ID == id {
			Posts[i].Body = reqBody.Body
			Posts[i].Name = reqBody.Name

			c.JSON(200, Posts[i])
			return
		}
	}

	c.JSON(422, ErrorMsg("post not found"))
}

func Create(c *gin.Context) {
	var reqBody Post
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(422, gin.H{
			"error": true,
			"message": "invalid JSON body",
		})
		return
	}

	// Validation
	errs := Validate(reqBody)
	if errs != nil {
		c.JSON(422, errs)
		return
	}
	reqBody.ID = uuid.New().String()
	Posts = append(Posts, reqBody)

	c.JSON(200, reqBody)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	for i, el := range Posts {
		if el.ID == id {
			Posts = append(Posts[:i], Posts[i+1:]...)

			c.JSON(200, gin.H{
				"error": false,
			})
			return
		}
	}

	c.JSON(404, gin.H{
		"error": true,
		"message": "post not found",
	})
}

func ErrorMsg(message string) DefaultError {
	return DefaultError{
		Error:   true,
		Message: message,
	}
}

func Validate(a interface{}) []DefaultError {
	v := validator.New()
	err := v.Struct(a)
	var errorsStruct []DefaultError
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			label := e.ActualTag() //e.g. required, email, etc
			structName := e.Field()

			switch label {
			case "required":
				errorsStruct = append(errorsStruct, DefaultError{
					Error:   true,
					Message: fmt.Sprintf("%v field is required", strings.ToLower(structName)),
					Field:   strings.ToLower(structName),
				})
				break
			default:
				errorsStruct = append(errorsStruct, DefaultError{
					Error:   true,
					Message: fmt.Sprintf("%v field is invalid", strings.ToLower(structName)),
					Field:   strings.ToLower(structName),
				})
				break
			}
		}
	}

	return errorsStruct
}