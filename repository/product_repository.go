package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	query := "SELECT id, product_name, price FROM product"

	rows, err := pr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer rows.Close()

	var productList []model.Product
	var productObj model.Product

	for rows.Next() {
		if err = rows.Scan(
			&productObj.ID,
			&productObj.Name,
			&productObj.Price,
		); err != nil {
			fmt.Println(err)
			return nil, err
		}

		productList = append(productList, productObj)

	}

	return productList, nil
}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {

	var id int
	query, err := pr.connection.Prepare("INSERT INTO product" +
		" (product_name, price)" +
		" VALUES ($1, $2) RETURNING id")

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	defer query.Close()

	err = query.QueryRow(
		product.Name,
		product.Price,
	).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}

func (pr *ProductRepository) GetProductById(id_product int) (*model.Product, error) {
	query, err := pr.connection.Prepare("SELECT * FROM product WHERE id = $1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer query.Close()

	var productObj model.Product

	err = query.QueryRow(id_product).Scan(
		&productObj.ID,
		&productObj.Name,
		&productObj.Price,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No product found with the given ID")
			return nil, nil
		}
		fmt.Println(err)
		return nil, err
	}

	return &productObj, nil
}
