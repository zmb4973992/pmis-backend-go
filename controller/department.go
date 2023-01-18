package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"pmis-backend-go/dto"
	"pmis-backend-go/serializer/response"
	"pmis-backend-go/service"
	"pmis-backend-go/util"
	"strconv"
)

type departmentController struct{}

func (departmentController) Get(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("department-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Common{
			Data:    nil,
			Code:    util.ErrorInvalidURIParameters,
			Message: util.GetMessage(util.ErrorInvalidURIParameters),
		})
		return
	}
	res := service.DepartmentService.Get(departmentID)
	c.JSON(http.StatusOK, res)
	return
}

func (departmentController) Create(c *gin.Context) {
	var param dto.DepartmentCreateOrUpdateDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.Creator = &userID
		param.LastModifier = &userID
	}

	res := service.DepartmentService.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (departmentController) Update(c *gin.Context) {
	var param dto.DepartmentCreateOrUpdateDTO
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("department-id"))
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理creator、lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.DepartmentService.Update(&param)
	c.JSON(200, res)
}

func (departmentController) Delete(c *gin.Context) {
	departmentID, err := strconv.Atoi(c.Param("department-id"))
	if err != nil {
		c.JSON(http.StatusOK, response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.DepartmentService.Delete(departmentID)
	c.JSON(http.StatusOK, res)
}

func (departmentController) List(c *gin.Context) {
	var param dto.DepartmentListDTO
	err := c.ShouldBindJSON(&param)
	//如果json没有传参，会提示EOF错误，这里允许正常运行；如果是其他错误，就正常报错
	if err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}

	tempRoleNames, exists := c.Get("role_names")
	if exists {
		roleNames := tempRoleNames.([]string)
		if len(roleNames) > 0 {
			param.RoleNames = roleNames
		}
	}

	tempBusinessDivisionIDs, exists := c.Get("business_division_ids")
	if exists {
		businessDivisionIDs := tempBusinessDivisionIDs.([]int)
		if len(businessDivisionIDs) > 0 {
			param.BusinessDivisionIDs = businessDivisionIDs
		}
	}

	tempDepartmentIDs, exists := c.Get("department_ids")
	if exists {
		departmentIDs := tempDepartmentIDs.([]int)
		if len(departmentIDs) > 0 {
			param.DepartmentIDs = departmentIDs
		}
	}

	//生成Service,然后调用它的方法
	res := service.DepartmentService.List(param)
	c.JSON(http.StatusOK, res)
}
