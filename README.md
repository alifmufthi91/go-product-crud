
# go-product-crud
Ini adalah aplikasi backend golang sederhana dimana user dapat mendaftar dan melakukan login. Setelah user melakukan login, User dapat mendaftarkan produk, melakukan perubahan pada produk, menampilkan data produk dan menghapus produk dari database menggunakan API.

## Tech
- Go 1.17.7
- PostgresQL

## Getting Started
``` sh
#  Clone this repo to your local machine using
git clone git@github.com:alifmufthi91/product-crud.git

#  Get into the directory
cd product-crud

#  Copy .env-example and modify the .env file
cp .env-example .env

#  Install dependencies
go build

#  Run App
go run .

## API List
|API|Routes|Method|
|----------------|-------------------------------|-----------------------------|
|User Registration|/api/users  |POST   | 
|User Login   |/api/users/login      |POST     |
|Get All User|/api/users  |GET  | 
|File Upload |/api/files/upload|POST|
|File Download|/api/files/:filename|GET|
|Register Product |/api/products|POST|
|Update Product |/api/products/:id|PUT|
|Delete Product |/api/products/:id|DELETE|
|Get All Products |/api/products|GET|
|Get Product |/api/products/:id|GET|

