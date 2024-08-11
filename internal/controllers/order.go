package controllers

import (
	"context"
	"dbo-test/internal/dal"
	"dbo-test/internal/model"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//	@Summary		Get Single Order
//	@Description	get single order by ID
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Order ID"
//	@Security		Bearer
//	@Success		200	{object}	model.Order
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/order/{id} [get]
func GetSingleOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid id",
		})
		return
	}

	order, err := dal.Order.Where(dal.Order.ID.Eq(int32(orderID))).First()
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, errorResponse{
				Status:  errorStatus,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"message": "order not found",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetMultipleOrder godoc
//	@Summary		Get Multiple Order
//	@Description	get multiple order with pagination and filtering options
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int		false	"Page number"					default(1)
//	@Param			pagesize	query	int		false	"Number of items per page"		default(10)
//	@Param			order		query	string	false	"Order by field (asc or desc)"	default("asc")
//	@Param			dateFrom	query	string	false	"Filter by order date from"		Format(date)
//	@Param			dateTo		query	string	false	"Filter by order date to"		Format(date)
//	@Param			amountFrom	query	number	false	"Filter by order amount from"	default(0)
//	@Param			amountTo	query	number	false	"Filter by order amount to"		default(0)
//	@Security		Bearer
//	@Success		200	{object}	PagedResults{data=[]model.Order}
//	@Failure		400	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/order [get]
func GetMultipleOrder(c *gin.Context) {
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

	amountFromStr := c.DefaultQuery("amountFrom", "0")
	amountToStr := c.DefaultQuery("amountTo", "0")

	amountFrom, err := strconv.ParseFloat(amountFromStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid amountFrom",
		})
		return
	}
	amountTo, err := strconv.ParseFloat(amountToStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid amountTo",
		})
		return
	}

	createdFromTime, err := parseTimeParam(c, "dateFrom")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: fmt.Sprintf("invalid dateFrom: %s", err.Error()),
		})
		return
	}

	createdToTime, err := parseTimeParam(c, "dateTo")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: fmt.Sprintf("invalid dateTo: %s", err.Error()),
		})
		return
	}

	orderParts := strings.Split(order, " ")
	descbBool := false
	if len(orderParts) > 1 {
		descbBool = strings.EqualFold(orderParts[1], "desc")
	}
	resp, totalRecords, err := queryMultipleOrder(page, pagesize, orderBy{
		Field: orderParts[0],
		Desc:  descbBool,
	}, createdFromTime, createdToTime, amountFrom, amountTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PagedResults{
		Page:         int64(page),
		PageSize:     int64(pagesize),
		Data:         resp,
		TotalRecords: int(totalRecords),
	})
}

// parseTimeParam parses a time parameter, setting it to midnight UTC if only a date is provided
func parseTimeParam(c *gin.Context, paramName string) (time.Time, error) {
	timeStr := c.Query(paramName)
	if timeStr == "" {
		return time.Time{}, nil
	}

	// Try parsing with different formats
	formats := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02T15:04:05Z",      // RFC3339 without timezone
		"2006-01-02",                // Date only
	}

	var parsedTime time.Time
	var err error
	for _, format := range formats {
		parsedTime, err = time.Parse(format, timeStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, err
	}

	// If only date was provided, set time to midnight UTC
	if len(timeStr) == 10 { // "YYYY-MM-DD" is 10 characters
		parsedTime = time.Date(parsedTime.Year(), parsedTime.Month(), parsedTime.Day(), 0, 0, 0, 0, time.UTC)
	}

	return parsedTime, nil
}

func queryMultipleOrder(
	page, pagesize int,
	order orderBy,
	dateFrom, dateTo time.Time,
	amountFrom, amountTo float64,
) ([]*model.Order, int64, error) {

	orderQuery := dal.Order
	resultOrm := orderQuery.WithContext(context.Background())

	if !dateFrom.IsZero() && !dateTo.IsZero() {
		resultOrm = resultOrm.Where(orderQuery.OrderDate.Gte(dateFrom), orderQuery.OrderDate.Lte(dateTo))
	} else if !dateFrom.IsZero() && dateTo.IsZero() {
		resultOrm = resultOrm.Where(orderQuery.OrderDate.Gte(dateFrom)).Debug()
	} else if dateFrom.IsZero() && !dateTo.IsZero() {
		resultOrm = resultOrm.Where(orderQuery.OrderDate.Lte(dateTo))
	}

	if amountFrom > 0 && amountTo > 0 {
		resultOrm = resultOrm.Where(orderQuery.Amount.Gte(amountFrom), orderQuery.Amount.Lte(amountTo))
	} else if amountFrom > 0 && amountTo <= 0 {
		resultOrm = resultOrm.Where(orderQuery.Amount.Gte(amountFrom))
	} else if amountFrom <= 0 && amountTo > 0 {
		resultOrm = resultOrm.Where(orderQuery.Amount.Lte(amountTo))
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

	orderCol, ok := orderQuery.GetFieldByName(order.Field)
	if ok {
		if order.Desc {
			resultOrm = resultOrm.Order(orderCol.Desc())
		} else {
			resultOrm = resultOrm.Order(orderCol)
		}
	}

	resp, err := resultOrm.Find()
	if err != nil {
		return nil, -1, err
	}
	return resp, totalRecords, nil
}

type createOrderReq struct {
	OrderDate  time.Time `json:"order_date" format:"date-time"`
	Amount     float64   `json:"amount"`
	CustomerID int32     `json:"customer_id"`
}

// CreateOrder godoc
//	@Summary		Create a new order
//	@Description	Create a new order with the provided details
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			order	body	createOrderReq	true	"Order details"
//	@Security		Bearer
//	@Success		200	{object}	successResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		401	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/order [post]
func CreateOrder(c *gin.Context) {
	var input createOrderReq
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	if err := dal.Order.Create(&model.Order{
		OrderDate:  input.OrderDate,
		Amount:     input.Amount,
		CustomerID: input.CustomerID,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Status: successStatus,
		Data:   nil,
	})
}

type updateOrderReq struct {
	OrderDate  time.Time `json:"order_date" format:"date-time"`
	Amount     float64   `json:"amount"`
	CustomerID int32     `json:"customer_id"`
}

// UpdateOrder godoc
//
//	@Summary		Update an existing order
//	@Description	Update the details of an existing order
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int				true	"Order ID"
//	@Param			order	body	updateOrderReq	true	"Updated order details"
//	@Security		Bearer
//	@Success		200	{object}	successResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		401	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/order/{id} [put]
func UpdateOrder(c *gin.Context) {
	var input updateOrderReq

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid id",
		})
		return
	}

	info, err := dal.Order.Where(dal.Order.ID.Eq(int32(orderID))).Updates(model.Order{
		OrderDate:  input.OrderDate,
		Amount:     input.Amount,
		CustomerID: input.CustomerID,
	})
	if info.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, errorResponse{
			Status:  errorStatus,
			Message: "order not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Status: successStatus,
		Data:   nil,
	})
}

// DeleteOrder godoc
//	@Summary		Delete an order
//	@Description	Delete an order by ID
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Order ID"
//	@Security		Bearer
//	@Success		200	{object}	successResponse
//	@Failure		400	{object}	errorResponse
//	@Failure		401	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Router			/order/{id} [delete]
func DeleteOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{
			Status:  errorStatus,
			Message: "invalid id",
		})
		return
	}

	info, err := dal.Order.Where(dal.Order.ID.Eq(int32(orderID))).Delete()
	if info.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, errorResponse{
			Status:  errorStatus,
			Message: "order not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Status:  errorStatus,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, successResponse{
		Status: successStatus,
		Data:   nil,
	})

}
