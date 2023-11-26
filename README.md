#  Parking lot of a rental agency

The objective of this project is to create a RESTful API in `Go lang` to manage the parking lot of a rental agency. 
Each car is characterized by 
- Model
- Registration
- Mileage
- Available status

## Technologies Used

* Golang
* MySQL
* GORM (Object Relational Mapping (ORM) library for Golang)
* Gorilla mux (HTTP router)
* Testify (Test toolkit)

## Installation and Configuration

Clone this repository to your local machine

```sh
git clone https://github.com/abdeel07/backend-go-cars.git
```

Gorilla mux

```sh
go get -u github.com/gorilla/mux
```

GORM & MySQL Driver

```sh
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

Testify

```sh
go get github.com/stretchr/testify
```

Create 2 MySQL database
1. parking_lot
2. parking_lot_test (For the test)

## Test the API

```sh
cd handlers_test/
go test
```
