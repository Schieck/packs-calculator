package http

import (
	"log/slog"
	"net/http"

	"github.com/Schieck/packs-calculator/internal/domain/errs"
	"github.com/Schieck/packs-calculator/internal/dto"
	calculatorUseCase "github.com/Schieck/packs-calculator/internal/usecase/pack_calculator"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CalculatorHandler struct {
	calculatePacksUseCase *calculatorUseCase.CalculatePacksUseCase
	logger                *slog.Logger
	validator             *validator.Validate
}

func NewCalculatorHandler(calculatePacksUseCase *calculatorUseCase.CalculatePacksUseCase, logger *slog.Logger) *CalculatorHandler {
	return &CalculatorHandler{
		calculatePacksUseCase: calculatePacksUseCase,
		logger:                logger,
		validator:             validator.New(),
	}
}

// Calculate handles pack calculation requests
// @Summary Calculate Optimal Packs
// @Description Calculate the optimal pack allocation for a given order quantity and available pack sizes
// @Tags calculator
// @Accept json
// @Produce json
// @Param request body dto.CalculationRequest true "Calculation parameters"
// @Success 200 {object} dto.CalculationResponse
// @Failure 400 {object} errs.ErrorResponse
// @Failure 500 {object} errs.ErrorResponse
// @Security BearerAuth
// @Router /calculate [post]
func (h *CalculatorHandler) Calculate(c *gin.Context) {
	var dtoReq dto.CalculationRequest

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

	result, err := h.calculatePacksUseCase.Execute(dtoReq.PackSizes, dtoReq.Items)
	if err != nil {
		h.logger.Error("Pack calculation use case failed", "error", err)
		c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
			Error: "Pack calculation failed",
		})
		return
	}

	dtoResponse := &dto.CalculationResponse{
		Allocation: result.Allocation.GetAllocation(),
		TotalPacks: result.Allocation.TotalPacks(),
		TotalItems: result.Allocation.TotalItems(),
		Surplus:    result.Surplus,
	}

	c.JSON(http.StatusOK, dtoResponse)
}
