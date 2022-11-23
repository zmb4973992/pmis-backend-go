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

type projectController struct{}

func (projectController) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.ProjectService.Get(id)
	c.JSON(http.StatusOK, res)
	return
}

func (projectController) Create(c *gin.Context) {
	var param dto.ProjectCreateOrUpdateDTO
	//先把json参数绑定到model
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

	res := service.ProjectService.Create(&param)
	c.JSON(http.StatusOK, res)
	return
}

func (projectController) CreateInBatches(c *gin.Context) {
	var param []dto.ProjectCreateOrUpdateDTO
	//先把json参数绑定到model
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
		for i := range param {
			param[i].Creator = &userID
			param[i].LastModifier = &userID
		}
	}

	res := service.ProjectService.CreateInBatches(param)
	c.JSON(http.StatusOK, res)
	return
}

func (projectController) Update(c *gin.Context) {
	var param dto.ProjectCreateOrUpdateDTO
	//先把json参数绑定到model
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidJSONParameters))
		return
	}
	//把uri上的id参数传递给结构体形式的入参
	param.ID, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}

	//处理lastModifier字段
	tempUserID, _ := c.Get("user_id")
	if tempUserID != nil {
		userID := tempUserID.(int)
		param.LastModifier = &userID
	}

	res := service.ProjectService.Update(&param)
	c.JSON(http.StatusOK, res)
}

func (projectController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK,
			response.Failure(util.ErrorInvalidURIParameters))
		return
	}
	res := service.ProjectService.Delete(id)
	c.JSON(http.StatusOK, res)
}

func (projectController) List(c *gin.Context) {
	var param dto.ProjectListDTO
	err := c.ShouldBindJSON(&param)

	if err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest,
			response.FailureForList(util.ErrorInvalidJSONParameters))
		return
	}

	tempTopRole, exists := c.Get("top_role")
	if exists {
		topRole := tempTopRole.(string)
		if topRole != "" {
			param.TopRole = topRole
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
	res := service.ProjectService.List(param)
	c.JSON(http.StatusOK, res)
}
