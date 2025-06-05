package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Schieck/packs-calculator/internal/domain/errs"
	"github.com/Schieck/packs-calculator/internal/dto"
	packConfigurationUseCase "github.com/Schieck/packs-calculator/internal/usecase/pack_configuration"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PackConfigurationHandler struct {
	getAllConfigurationsUseCase    *packConfigurationUseCase.GetAllConfigurationsUseCase
	getConfigurationByIDUseCase    *packConfigurationUseCase.GetConfigurationByIDUseCase
	getDefaultConfigurationUseCase *packConfigurationUseCase.GetDefaultConfigurationUseCase
	createConfigurationUseCase     *packConfigurationUseCase.CreateConfigurationUseCase
	updateConfigurationUseCase     *packConfigurationUseCase.UpdateConfigurationUseCase
	deleteConfigurationUseCase     *packConfigurationUseCase.DeleteConfigurationUseCase
	setDefaultConfigurationUseCase *packConfigurationUseCase.SetDefaultConfigurationUseCase
	logger                         *slog.Logger
	validator                      *validator.Validate
}

func NewPackConfigurationHandler(
	getAllConfigurationsUseCase *packConfigurationUseCase.GetAllConfigurationsUseCase,
	getConfigurationByIDUseCase *packConfigurationUseCase.GetConfigurationByIDUseCase,
	getDefaultConfigurationUseCase *packConfigurationUseCase.GetDefaultConfigurationUseCase,
	createConfigurationUseCase *packConfigurationUseCase.CreateConfigurationUseCase,
	updateConfigurationUseCase *packConfigurationUseCase.UpdateConfigurationUseCase,
	deleteConfigurationUseCase *packConfigurationUseCase.DeleteConfigurationUseCase,
	setDefaultConfigurationUseCase *packConfigurationUseCase.SetDefaultConfigurationUseCase,
	logger *slog.Logger,
) *PackConfigurationHandler {
	return &PackConfigurationHandler{
		getAllConfigurationsUseCase:    getAllConfigurationsUseCase,
		getConfigurationByIDUseCase:    getConfigurationByIDUseCase,
		getDefaultConfigurationUseCase: getDefaultConfigurationUseCase,
		createConfigurationUseCase:     createConfigurationUseCase,
		updateConfigurationUseCase:     updateConfigurationUseCase,
		deleteConfigurationUseCase:     deleteConfigurationUseCase,
		setDefaultConfigurationUseCase: setDefaultConfigurationUseCase,
		logger:                         logger,
		validator:                      validator.New(),
	}
}

// GetAllConfigurations handles GET /pack-configurations
// @Summary Get All Pack Configurations
// @Description Retrieve all active pack configurations
// @Tags pack-configurations
// @Accept json
// @Produce json
// @Success 200 {object} dto.PackConfigurationListResponse
// @Failure 500 {object} errs.ErrorResponse
// @Security BearerAuth
// @Router /pack-configurations [get]
func (h *PackConfigurationHandler) GetAllConfigurations(c *gin.Context) {
	configurations, err := h.getAllConfigurationsUseCase.Execute()
	if err != nil {
		h.logger.Error("Get all pack configurations use case failed", "error", err)
		c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Error: "Failed to retrieve pack configurations",
		})
		return
	}

	response := dto.ToPackConfigurationListResponse(configurations)
	c.JSON(http.StatusOK, response)
}

// GetConfigurationByID handles GET /pack-configurations/:id
// @Summary Get Pack Configuration by ID
// @Description Retrieve a specific pack configuration by its ID
// @Tags pack-configurations
// @Accept json
// @Produce json
// @Param id path int true "Pack Configuration ID"
// @Success 200 {object} dto.PackConfigurationResponse
// @Failure 400 {object} errs.ErrorResponse
// @Failure 404 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Security BearerAuth
// @Router /pack-configurations/{id} [get]
func (h *PackConfigurationHandler) GetConfigurationByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Warn("Invalid pack configuration ID", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Error: "Invalid pack configuration ID",
		})
		return
	}

	configuration, err := h.getConfigurationByIDUseCase.Execute(id)
	if err != nil {
		h.logger.Error("Get pack configuration by ID use case failed", "id", id, "error", err)
		c.JSON(http.StatusNotFound, errs.ErrorResponse{
			Error: "Pack configuration not found",
		})
		return
	}

	response := dto.ToPackConfigurationResponse(configuration)
	c.JSON(http.StatusOK, response)
}

// GetDefaultConfiguration handles GET /pack-configurations/default
// @Summary Get Default Pack Configuration
// @Description Retrieve the default pack configuration
// @Tags pack-configurations
// @Accept json
// @Produce json
// @Success 200 {object} dto.PackConfigurationResponse
// @Failure 404 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Security BearerAuth
// @Router /pack-configurations/default [get]
func (h *PackConfigurationHandler) GetDefaultConfiguration(c *gin.Context) {
	configuration, err := h.getDefaultConfigurationUseCase.Execute()
	if err != nil {
		h.logger.Error("Get default pack configuration use case failed", "error", err)
		c.JSON(http.StatusNotFound, errs.ErrorResponse{
			Error: "No default pack configuration found",
		})
		return
	}

	response := dto.ToPackConfigurationResponse(configuration)
	c.JSON(http.StatusOK, response)
}

// CreateConfiguration handles POST /pack-configurations
// @Summary Create Pack Configuration
// @Description Create a new pack configuration
// @Tags pack-configurations
// @Accept json
// @Produce json
// @Param request body dto.CreatePackConfigurationRequest true "Pack configuration data"
// @Success 201 {object} dto.PackConfigurationResponse
// @Failure 400 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Security BearerAuth
// @Router /pack-configurations [post]
func (h *PackConfigurationHandler) CreateConfiguration(c *gin.Context) {
	var dtoReq dto.CreatePackConfigurationRequest

	if err := c.ShouldBindJSON(&dtoReq); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	if err := h.validator.Struct(&dtoReq); err != nil {
		h.logger.Warn("Request validation failed", "error", err)
		c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Error:   "Validation failed",
			Details: errs.FormatValidationErrors(err),
		})
		return
	}

	configuration, err := h.createConfigurationUseCase.Execute(dtoReq.Name, dtoReq.PackSizes)
	if err != nil {
		h.logger.Error("Create pack configuration use case failed", "error", err)
		c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Error: "Failed to create pack configuration",
		})
		return
	}

	response := dto.ToPackConfigurationResponse(configuration)
	c.JSON(http.StatusCreated, response)
}

// UpdateConfiguration handles PUT /pack-configurations/:id
// @Summary Update Pack Configuration
// @Description Update an existing pack configuration
// @Tags pack-configurations
// @Accept json
// @Produce json
// @Param id path int true "Pack Configuration ID"
// @Param request body dto.UpdatePackConfigurationRequest true "Updated pack configuration data"
// @Success 200 {object} dto.PackConfigurationResponse
// @Failure 400 {object} errs.ErrorResponse
// @Failure 404 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Security BearerAuth
// @Router /pack-configurations/{id} [put]
func (h *PackConfigurationHandler) UpdateConfiguration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Warn("Invalid pack configuration ID", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Error: "Invalid pack configuration ID",
		})
		return
	}

	var dtoReq dto.UpdatePackConfigurationRequest

	if err := c.ShouldBindJSON(&dtoReq); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	if err := h.validator.Struct(&dtoReq); err != nil {
		h.logger.Warn("Request validation failed", "error", err)
		c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Error:   "Validation failed",
			Details: errs.FormatValidationErrors(err),
		})
		return
	}

	configuration, err := h.updateConfigurationUseCase.Execute(id, dtoReq.Name, dtoReq.PackSizes, dtoReq.IsDefault)
	if err != nil {
		h.logger.Error("Update pack configuration use case failed", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Error: "Failed to update pack configuration",
		})
		return
	}

	response := dto.ToPackConfigurationResponse(configuration)
	c.JSON(http.StatusOK, response)
}

// DeleteConfiguration handles DELETE /pack-configurations/:id
// @Summary Delete Pack Configuration
// @Description Delete a pack configuration (soft delete)
// @Tags pack-configurations
// @Accept json
// @Produce json
// @Param id path int true "Pack Configuration ID"
// @Success 204
// @Failure 400 {object} errs.ErrorResponse
// @Failure 404 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Security BearerAuth
// @Router /pack-configurations/{id} [delete]
func (h *PackConfigurationHandler) DeleteConfiguration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Warn("Invalid pack configuration ID", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Error: "Invalid pack configuration ID",
		})
		return
	}

	err = h.deleteConfigurationUseCase.Execute(id)
	if err != nil {
		h.logger.Error("Delete pack configuration use case failed", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Error: "Failed to delete pack configuration",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// SetDefaultConfiguration handles PATCH /pack-configurations/:id/default
// @Summary Set Default Pack Configuration
// @Description Set a specific pack configuration as the default
// @Tags pack-configurations
// @Accept json
// @Produce json
// @Param id path int true "Pack Configuration ID"
// @Success 200 {object} dto.PackConfigurationResponse
// @Failure 400 {object} errs.ErrorResponse
// @Failure 404 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Security BearerAuth
// @Router /pack-configurations/{id}/default [patch]
func (h *PackConfigurationHandler) SetDefaultConfiguration(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Warn("Invalid pack configuration ID", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, errs.ErrorResponse{
			Error: "Invalid pack configuration ID",
		})
		return
	}

	err = h.setDefaultConfigurationUseCase.Execute(id)
	if err != nil {
		h.logger.Error("Set default pack configuration use case failed", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Error: "Failed to set default pack configuration",
		})
		return
	}

	configuration, err := h.getConfigurationByIDUseCase.Execute(id)
	if err != nil {
		h.logger.Error("Get pack configuration by ID use case failed after setting default", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Error: "Failed to retrieve updated pack configuration",
		})
		return
	}

	response := dto.ToPackConfigurationResponse(configuration)
	c.JSON(http.StatusOK, response)
}
