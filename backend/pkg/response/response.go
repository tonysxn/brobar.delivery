package response

import (
	"github.com/gofiber/fiber/v3"
)

type SuccessResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type Pagination struct {
	TotalCount int    `json:"total_count"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	OrderBy    string `json:"order_by"`
	OrderDir   string `json:"order_dir"`
}

type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type paginated interface {
	isPaginated()
	GetData() interface{}
	GetPagination() Pagination
}

func (p PaginatedResponse[T]) isPaginated() {}

func (p PaginatedResponse[T]) GetData() interface{} {
	return p.Data
}

func (p PaginatedResponse[T]) GetPagination() Pagination {
	return p.Pagination
}

func Success(c fiber.Ctx, data interface{}) error {
	resp := SuccessResponse{
		Success: true,
	}

	if p, ok := data.(paginated); ok {
		resp.Data = p.GetData()
		resp.Pagination = p.GetPagination()
	} else {
		resp.Data = data
	}

	if resp.Data == nil {
		resp.Data = make([]string, 0)
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func Error(c fiber.Ctx, status int, err error) error {
	return c.Status(status).JSON(ErrorResponse{
		Success: false,
		Error:   err.Error(),
	})
}

func NotFound(c fiber.Ctx) error {
	return Error(c, fiber.StatusNotFound, fiber.ErrNotFound)
}

func BadRequest(c fiber.Ctx, err error) error {
	return Error(c, fiber.StatusBadRequest, err)
}

func InternalServerError(c fiber.Ctx) error {
	return Error(c, fiber.StatusInternalServerError, fiber.ErrInternalServerError)
}

type FileServiceResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Filename string `json:"filename"`
		Message  string `json:"message"`
		Path     string `json:"path"`
	} `json:"data"`
}
