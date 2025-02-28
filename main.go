package main

import (
	// "encoding/json"
	// "fmt"
	// "net/http"
	"strconv"
	// "time"

	// "github.com/gorilla/mux"
	"github.com/gofiber/fiber/v2"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var products []Product

func createProduct(c *fiber.Ctx) error {
	var newProduct Product
	if err := c.BodyParser(&newProduct); err != nil {
		c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	products = append(products, newProduct)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Product created successfully",
		"product": newProduct,
	})
}

func getProduct(c *fiber.Ctx) error {
	return c.JSON(products)
}

func getProductById(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		fiber.NewError(404, err.Error())
	}

	for _, product := range products {
		if product.Id == id {
			return c.JSON(product)
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("product not found")
}

func updateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid product ID")
	}

	found := false
	for index, product := range products {
		if product.Id == id {
			found = true
			products[index].Name = c.FormValue("name")
			products[index].Description = c.FormValue("description")
			price, err := strconv.ParseFloat(c.FormValue("price"), 64)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			products[index].Price = price
		}
	}
	if !found {
		return fiber.NewError(404, "Not Found")

	}

	return fiber.NewError(404, "Can Update Product")
}

func deleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())

	}
	for i, product := range products {
		if product.Id == id {
			products = append(products[:i], products[i+1:]...)
			return c.SendString("Product Deleted")
		}
	}

	return fiber.NewError(fiber.StatusNotFound, "Product not found")
}

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func handleError(c *fiber.Ctx) error {
	response := Response{
		Message: "An error occurred",
		Status:  fiber.StatusInternalServerError,
	}
	return c.Status(response.Status).JSON(response)
}

func main() {
	// mux := mux.NewRouter()
	app := fiber.New()

	app.Post("/products", createProduct)
	app.Get("/products", getProduct)
	app.Get("/products/:id", getProductById)
	app.Put("/products", updateProduct)
	app.Delete("/products", deleteProduct)
	app.Get("/products", handleError)

	app.Listen(":8080")

}
