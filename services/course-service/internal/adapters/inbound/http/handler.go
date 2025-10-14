package http

import (
	"time"

	"github.com/wnmay/horo/services/course-service/internal/app"
	"github.com/wnmay/horo/services/course-service/internal/domain"

	"github.com/gofiber/fiber/v2"
)

// Handler struct ที่รับ service มาจากชั้น app
type Handler struct {
	service app.CourseService
}

// NewHandler — constructor
func NewHandler(s app.CourseService) *Handler {
	return &Handler{service: s}
}

// Register — รวม endpoint ทั้งหมดของ course service
func (h *Handler) Register(router fiber.Router) {
	router.Post("/courses", h.CreateCourse)
	router.Get("/courses/:id", h.GetCourseByID)
	router.Get("/prophets/:prophet_id/courses", h.ListCoursesByProphet)
}

// CreateCourse — POST /courses
func (h *Handler) CreateCourse(c *fiber.Ctx) error {
	var req struct {
		ProphetID   string  `json:"prophet_id"`
		CourseName  string  `json:"coursename"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Duration    int32   `json:"duration"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	input := app.CreateCourseInput{
		ProphetID:   req.ProphetID,
		CourseName:  req.CourseName,
		Description: req.Description,
		Price:       req.Price,
		Duration:    domain.DurationEnum(req.Duration),
	}

	course, err := h.service.CreateCourse(input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(course)
}

// GetCourseByID — GET /courses/:id
func (h *Handler) GetCourseByID(c *fiber.Ctx) error {
	id := c.Params("id")

	course, err := h.service.GetCourseByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "course not found")
	}

	return c.JSON(course)
}

// ListCoursesByProphet — GET /prophets/:prophet_id/courses
func (h *Handler) ListCoursesByProphet(c *fiber.Ctx) error {
	prophetID := c.Params("prophet_id")

	courses, err := h.service.ListCoursesByProphet(prophetID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := struct {
		Timestamp time.Time        `json:"timestamp"`
		Data      []*domain.Course `json:"courses"`
	}{
		Timestamp: time.Now(),
		Data:      courses,
	}

	return c.JSON(response)
}
