package main

import (
	"challenge-goapi/config"
	"challenge-goapi/models"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var db = config.ConnectDB()

func main() {
	defer db.Close()
	router := gin.Default()

	apiG := router.Group("/api")

	customersG := apiG.Group("/customers")
	customersG.GET("/", getCustomers)
	// BODY SHOULD SPECIFY MANUALLY A UNIQUE ID
	customersG.POST("/", addCustomer)
	customersG.DELETE("/:id", deleteCustomer)
	// ATRIBUT KOSONG SUDAH DI HANDLE
	customersG.PUT("/:id", updateCustomer)

	servicesG := apiG.Group("/services")
	servicesG.GET("/", getServices)
	// BODY SHOULD SPECIFY MANUALLY A UNIQUE ID
	servicesG.POST("/", addService)
	servicesG.DELETE("/:id", deleteService)
	// EMPTY ATTRIBUTE NOT HANDLED YET, SO WHEN ONE OR MORE ATTRIBUTE IS EMPTY IT WILL STILL UPDATE THE ATTRIBUTE TO ZERO VALUE
	servicesG.PUT("/:id", updateService)

	trxs := apiG.Group("/transactions")
	trxs.GET("/", getTrxsLaundry)
	// BODY SHOULD SPECIFY MANUALLY A UNIQUE ID
	trxs.POST("/", enrollLaundry)

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func enrollLaundry(c *gin.Context) {
	e := models.LaundryEnrollment{}
	err := c.ShouldBind(&e)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid data"})
		return
	}

	ld := models.LaundryDetail{Id: e.Id, Service_Id: e.Service_Id, Quantity: e.Quantity}
	tl := models.TrxLaundry{Id: e.Id, Customer_Id: e.Customer_Id, Entry_Date: time.Now(), Finish_Date: time.Now()}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	insertLaundryDetail(ld, tx, c)
	insertTrxLaundry(tl, ld.Id, tx, c)
	// taken_credit := getSumCreditOfStudentEnrollment(se.Student_Id, tx)
	// updateStudentCredit(taken_credit, se.Student_Id, tx)

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Laundry succcessfully enrolled!"})
}

func insertLaundryDetail(ld models.LaundryDetail, tx *sql.Tx, c *gin.Context) {
	sqlStatement := "INSERT INTO laundry_detail (id, service_id, quantity, total_price) VALUES ($1, $2, $3, $4);"

	s := getServiceById(ld.Service_Id)
	total_price := s.Price * ld.Quantity

	_, err := tx.Exec(sqlStatement, ld.Id, ld.Service_Id, ld.Quantity, total_price)
	rollbackValidate(err, tx, c)
}

func getServiceById(id int) models.Service {
	sqlStatement := "SELECT * FROM mst_service WHERE id = $1;"
	c := models.Service{}
	err := db.QueryRow(sqlStatement, id).Scan(&c.Id, &c.Name, &c.Price, &c.Unit_type_id)
	if err != nil {
		panic(err)
	}

	return c
}

func insertTrxLaundry(tl models.TrxLaundry, id int, tx *sql.Tx, c *gin.Context) {
	sqlStatement := "INSERT INTO trx_laundry (id, customer_Id, entry_date, finish_date, laundry_detail_id) VALUES ($1, $2, $3, $4, $5);"

	_, err := tx.Exec(sqlStatement, tl.Id, tl.Customer_Id, tl.Entry_Date, tl.Finish_Date, id)
	rollbackValidate(err, tx, c)
}

func rollbackValidate(err error, tx *sql.Tx, c *gin.Context) {
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}
	}
}

// END TRX

func getTrxsLaundry(c *gin.Context) {
	sqlStatement := "SELECT id, customer_id, entry_date, finish_date, laundry_detail_id FROM trx_laundry;"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	defer rows.Close()

	laundrys := []models.Laundry{}

	for rows.Next() {
		t := models.TrxLaundry{}
		err := rows.Scan(&t.Id, &t.Customer_Id, &t.Entry_Date, &t.Finish_Date, &t.Laundry_Detail_Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}

		l := models.Laundry{Id: t.Id, Entry_Date: t.Entry_Date, Finish_Date: t.Finish_Date}

		cust := getCustomerById(t.Customer_Id)
		l.Customer_Name = cust.Name
		l.Customer_Phone = cust.Phone

		ld := getLaundryDetailById(t.Laundry_Detail_Id)
		svc := getServiceById(ld.Service_Id)
		l.Service_Name = svc.Name
		l.Quantity = ld.Quantity
		l.Total_Price = ld.Total_Price

		laundrys = append(laundrys, l)
	}

	err = rows.Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, laundrys)
}

func updateCustomer(c *gin.Context) {
	var cust models.Customer
	err := c.ShouldBind(&cust)
	id := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid data"})
		return
	}

	sqlStatement := "UPDATE mst_customer SET "
	var params []interface{}
	params = append(params, id)

	if cust.Phone != "" {
		sqlStatement += "phone = $2"
		params = append(params, cust.Phone)
	}

	if cust.Name != "" {
		if cust.Phone == "" {
			sqlStatement += "name = $2"
		} else {
			sqlStatement += ", name = $3"
		}
		params = append(params, cust.Name)
	}

	sqlStatement += " WHERE id = $1;"

	_, err = db.Exec(sqlStatement, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Customer successfully updated!"})
}

func updateService(c *gin.Context) {
	var s models.Service
	err := c.ShouldBind(&s)
	id := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid data"})
		return
	}

	sqlStatement := "UPDATE mst_service SET name = $2, price = $3, unit_type_id = $4  WHERE id = $1;"

	_, err = db.Exec(sqlStatement, id, s.Name, s.Price, s.Unit_type_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Service successfully updated!"})
}

func deleteCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	sqlStatement := "DELETE FROM mst_customer WHERE id = $1;"

	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Customer successfully deleted!"})
}

func getCustomerById(id int) models.Customer {
	sqlStatement := "SELECT * FROM mst_customer WHERE id = $1;"
	c := models.Customer{}
	err := db.QueryRow(sqlStatement, id).Scan(&c.Id, &c.Name, &c.Phone)
	if err != nil {
		panic(err)
	}

	return c
}

func getLaundryDetailById(id int) models.LaundryDetail {
	sqlStatement := "SELECT * FROM laundry_detail WHERE id = $1;"
	c := models.LaundryDetail{}
	err := db.QueryRow(sqlStatement, id).Scan(&c.Id, &c.Service_Id, &c.Quantity, &c.Total_Price)
	if err != nil {
		panic(err)
	}

	return c
}

func deleteService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	sqlStatement := "DELETE FROM mst_service WHERE id = $1;"

	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Service successfully deleted!"})
}

func addCustomer(c *gin.Context) {
	var cust models.Customer
	err := c.ShouldBind(&cust)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid data"})
		return
	}

	sqlStatement := "INSERT INTO mst_customer (id, name, phone) VALUES ($1, $2, $3);"

	_, err = db.Exec(sqlStatement, cust.Id, cust.Name, cust.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "Customer successfully created!"})
}

func addService(c *gin.Context) {
	var s models.Service
	err := c.ShouldBind(&s)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": "Invalid data"})
		return
	}

	sqlStatement := "INSERT INTO mst_service (id, name, price, unit_type_id) VALUES ($1, $2, $3, $4);"

	_, err = db.Exec(sqlStatement, s.Id, s.Name, s.Price, s.Unit_type_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"code": http.StatusCreated, "message": "Service successfully created!"})
}

func getCustomers(c *gin.Context) {
	sqlStatement := "SELECT * FROM mst_customer;"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	defer rows.Close()

	customers := []models.Customer{}

	for rows.Next() {
		cust := models.Customer{}
		err := rows.Scan(&cust.Id, &cust.Name, &cust.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}

		customers = append(customers, cust)
	}

	err = rows.Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func getServices(c *gin.Context) {
	sqlStatement := "SELECT * FROM mst_service;"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}
	defer rows.Close()

	services := []models.Service{}

	for rows.Next() {
		s := models.Service{}
		err := rows.Scan(&s.Id, &s.Name, &s.Price, &s.Unit_type_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}

		services = append(services, s)
	}

	err = rows.Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, services)
}
