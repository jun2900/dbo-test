package controllers

import (
	"context"
	"dbo-test/internal/dal"
	"dbo-test/internal/model"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetSingleCustomer godoc
//
//	@Summary		Get a single customer
//	@Description	Get details of a specific customer by ID
//	@Tags			customers
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Customer ID"
//	@Security		Bearer
//	@Success		200	{object}	model.Customer
//	@Failure		400	{object}	errorResponse
//	@Failure		401	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/customer/{id} [get]
func GetSingleCustomer(c *gin.Context) {
	customerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid id",
		})
		return
	}

	customer, err := dal.Customer.Where(dal.Customer.ID.Eq(int32(customerID))).First()
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, errorResponse{
				Status:  errorStatus,
				Message: err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, customer)
}

// GetMultipleCustomer godoc
//
//	@Summary		Get multiple customers
//	@Description	Get a list of customers with pagination and filtering options
//	@Tags			customers
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int		false	"Page number"					default(1)
//	@Param			pagesize	query	int		false	"Number of items per page"		default(10)
//	@Param			order		query	string	false	"Order by field (asc or desc)"	default("asc")
//	@Param			name		query	string	false	"Filter by name"
//	@Param			email		query	string	false	"Filter by email"
//	@Param			phone		query	string	false	"Filter by phone"
//	@Security		Bearer
//	@Success		200	{object}	successResponse{data=PagedResults{data=[]model.Customer}}
//	@Failure		400	{object}	errorResponse
//	@Failure		401	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/customer [get]
func GetMultipleCustomer(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid page",
		})
		return
	}
	pagesize, err := strconv.Atoi(c.DefaultQuery("pagesize", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid pagesize",
		})
		return
	}
	order := c.DefaultQuery("order", "asc")
	name := c.DefaultQuery("name", "")
	email := c.DefaultQuery("email", "")
	phone := c.DefaultQuery("phone", "")

	orderParts := strings.Split(order, " ")
	descbBool := false
	if len(orderParts) > 1 {
		descbBool = strings.EqualFold(orderParts[1], "desc")
	}

	resp, totalRecords, err := queryMultipleCustomer(page, pagesize, orderBy{
		Field: orderParts[0],
		Desc:  descbBool,
	}, name, email, phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Status: "success",
		Data: PagedResults{
			Page:         int64(page),
			PageSize:     int64(pagesize),
			Data:         resp,
			TotalRecords: int(totalRecords),
		},
	})
}
func queryMultipleCustomer(
	page, pagesize int,
	order orderBy,
	name, email, phone string,
) ([]*model.Customer, int64, error) {

	customerQuery := dal.Customer
	resultOrm := customerQuery.WithContext(context.Background())

	if name != "" {
		resultOrm = resultOrm.Where(dal.Customer.Name.Like("%" + name + "%"))
	}
	if email != "" {
		resultOrm = resultOrm.Where(dal.Customer.Email.Like("%" + email + "%"))
	}
	if phone != "" {
		resultOrm = resultOrm.Where(dal.Customer.Phone.Like("%" + phone + "%"))
	}

	totalRecords, err := resultOrm.Count()
	if err != nil {
		return nil, 0, err
	}

	if page > 0 {
		offset := (page - 1) * pagesize
		resultOrm = resultOrm.Offset(offset).Limit(pagesize)
	} else {
		resultOrm = resultOrm.Limit(pagesize)
	}

	orderCol, ok := customerQuery.GetFieldByName(order.Field)
	if ok {
		if order.Desc {
			resultOrm = resultOrm.Order(orderCol.Desc())
		} else {
			resultOrm = resultOrm.Order(orderCol)
		}
	}

	resp, err := resultOrm.Find()
	if err != nil {
		return nil, 0, err
	}

	return resp, totalRecords, nil
}

type createCustomerReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// CreateCustomer godoc
//
//	@Summary		Create a new customer
//	@Description	Create a new customer with the provided details
//	@Tags			customers
//	@Accept			json
//	@Produce		json
//	@Param			customer	body	createCustomerReq	true	"Customer details"
//	@Security		Bearer
//	@Success		200	{object}	successResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		401	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/customer [post]
func CreateCustomer(c *gin.Context) {
	var input createCustomerReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	if err := dal.Customer.Create(&model.Customer{
		Name:  input.Name,
		Email: input.Email,
		Phone: input.Phone,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Status: "success",
	})
}

type updateCustomerReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// UpdateCustomer godoc
//
//	@Summary		Update an existing customer
//	@Description	Update the details of an existing customer
//	@Tags			customers
//	@Accept			json
//	@Produce		json
//	@Param			id			path	int					true	"Customer ID"
//	@Param			customer	body	updateCustomerReq	true	"Updated customer details"
//	@Security		Bearer
//	@Success		200	{object}	successResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		401	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/customer/{id} [put]
func UpdateCustomer(c *gin.Context) {
	customerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid id",
		})
		return
	}

	var input updateCustomerReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	_, err = dal.Customer.Where(dal.Customer.ID.Eq(int32(customerID))).Updates(&model.Customer{
		Name:  input.Name,
		Email: input.Email,
		Phone: input.Phone,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Status: "success",
	})
}

// DeleteCustomer godoc
//
//	@Summary		Delete a customer
//	@Description	Delete a customer by ID
//	@Tags			customers
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Customer ID"
//	@Security		Bearer
//	@Success		200	{object}	successResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		401	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/customer/{id} [delete]
func DeleteCustomer(c *gin.Context) {
	customerID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid id",
		})
		return
	}

	info, err := dal.Customer.Where(dal.Customer.ID.Eq(int32(customerID))).Delete()
	if info.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, errorResponse{
			Status:  errorStatus,
			Message: "order not found",
		})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, successResponse{
		Status: "success",
	})
}
