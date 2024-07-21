package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Dog struct {
	gorm.Model
	Name   string  `json:"name" gorm:"not null"`
	Breed  string  `json:"breed" gorm:"not null"`
	Age    int8    `json:"age" gorm:"not null"`
	Weight float32 `json:"weight" gorm:"default:0.0;not null"`
}

type FormattedDog struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Breed  string  `json:"breed"`
	Age    int8    `json:"age"`
	Weight float32 `json:"weight"`
}

var db *gorm.DB

func connectToDatabase() {
	var err error
	db, err = gorm.Open(
		"mysql",
		"root:12345@/dogs?charset=utf8&parseTime=True&loc=Local",
	)
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Dog{})
}

func createDog(c *gin.Context) {
	age, err := strconv.Atoi(c.PostForm("age"))
	if err != nil {
		log.Println(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Age is required!",
			},
		)
		return
	}
	weight, err := strconv.Atoi(c.PostForm("weight"))
	if err != nil {
		log.Println(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Weight is required!",
			},
		)
		return
	}
	if c.PostForm("name") == "" || c.PostForm("breed") == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Dog's data is required!",
			},
		)
		return
	}
	dog := Dog{
		Name:   c.PostForm("name"),
		Age:    int8(age),
		Breed:  c.PostForm("breed"),
		Weight: float32(weight),
	}
	db.Save(&dog)
	c.JSON(
		http.StatusCreated,
		gin.H{
			"status":     http.StatusCreated,
			"message":    "Dog item created successfully!",
			"resourceId": dog.ID,
		},
	)
}

func getDogs(c *gin.Context) {
	var dogs []Dog
	var _dogs []FormattedDog
	db.Find(&dogs)
	if len(dogs) <= 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No dog found!",
			},
		)
		return
	}
	for _, dog := range dogs {
		_dogs = append(
			_dogs,
			FormattedDog{
				ID:     dog.ID,
				Name:   dog.Name,
				Breed:  dog.Breed,
				Age:    dog.Age,
				Weight: dog.Weight,
			},
		)
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   _dogs,
		},
	)
}

func getPaginatedDogs(c *gin.Context) {
	var dogs []Dog
	var _dogs []FormattedDog
	pageSize := 2
	page, err := strconv.Atoi(c.DefaultQuery("pag", "1"))
	if err != nil {
		log.Println(err)
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid page number!",
			},
		)
		return
	}
	offset := (page - 1) * pageSize
	db.Limit(pageSize).Offset(offset).Find(&dogs)
	if len(dogs) <= 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No dog found!",
			},
		)
		return
	}
	for _, dog := range dogs {
		_dogs = append(
			_dogs,
			FormattedDog{
				ID:     dog.ID,
				Name:   dog.Name,
				Breed:  dog.Breed,
				Age:    dog.Age,
				Weight: dog.Weight,
			},
		)
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   _dogs,
		},
	)
}

func getPuppies(c *gin.Context) {
	var dogs []Dog
	var _dogs []FormattedDog
	db.Where("age < ?", 2).Find(&dogs)
	if len(dogs) <= 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No dog found!",
			},
		)
		return
	}
	for _, dog := range dogs {
		_dogs = append(
			_dogs,
			FormattedDog{
				ID:     dog.ID,
				Name:   dog.Name,
				Breed:  dog.Breed,
				Age:    dog.Age,
				Weight: dog.Weight,
			},
		)
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   _dogs,
		},
	)
}

func getDog(c *gin.Context) {
	var dog Dog
	dogID := c.Param("id")
	db.First(&dog, dogID)
	if dog.ID == 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No dog found!",
			},
		)
		return
	}
	_dog := FormattedDog{
		ID:     dog.ID,
		Name:   dog.Name,
		Breed:  dog.Breed,
		Age:    dog.Age,
		Weight: dog.Weight,
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status": http.StatusOK,
			"data":   _dog,
		},
	)
}

func updateDog(c *gin.Context) {
	var dog Dog
	dogID := c.Param("id")
	db.First(&dog, dogID)
	if dog.ID == 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No dog found!",
			},
		)
		return
	}
	age, err := strconv.Atoi(c.PostForm("age"))
	if err == nil {
		db.Model(&dog).Update("age", int8(age))
	}
	weight, err := strconv.Atoi(c.PostForm("weight"))
	if err == nil {
		db.Model(&dog).Update("weight", int8(weight))
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Dog updated successfully!",
		},
	)
}

func deleteDog(c *gin.Context) {
	var dog Dog
	dogID := c.Param("id")
	db.First(&dog, dogID)
	if dog.ID == 0 {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  http.StatusNotFound,
				"message": "No dog found!",
			},
		)
		return
	}
	db.Delete(&dog)
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Dog deleted successfully!",
		},
	)
}

func main() {
	connectToDatabase()

	router := gin.Default()

	api := router.Group("/api/v1/dogs")
	{
		api.GET("/", getDogs)
		api.GET("/:id", getDog)
		api.GET("/puppies", getPuppies)
		api.GET("/pages", getPaginatedDogs)
		api.POST("/", createDog)
		api.PUT("/:id", updateDog)
		api.DELETE("/:id", deleteDog)
	}
	router.Run()
}
