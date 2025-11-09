package http

import (
	"time"

	"github.com/wnmay/horo/services/course-service/internal/app"
	"github.com/wnmay/horo/services/course-service/internal/domain"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service app.CourseService
}

// NewHandler — constructor
func NewHandler(s app.CourseService) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(router fiber.Router) {
	router.Post("/courses", h.CreateCourse)
	router.Get("/courses/:id", h.GetCourseByID)
	router.Get("/prophets/:prophet_id/courses", h.ListCoursesByProphet)
	router.Patch("/courses/:id", h.UpdateCourse)
	router.Patch("/courses/delete/:id", h.DeleteCourse)
	router.Get("/courses", h.FindCoursesByFilter)
}

// CreateCourse — POST /courses
func (h *Handler) CreateCourse(c *fiber.Ctx) error {
	var req struct {
		ProphetID   string  `json:"prophet_id"`
		ProphetName string  `json:"prophetname"`
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
		ProphetName: req.ProphetName,
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

// UpdateCourse — PATCH /courses/:id
func (h *Handler) UpdateCourse(c *fiber.Ctx) error {
	id := c.Params("id")
	var in domain.UpdateCourseInput
	if err := c.BodyParser(&in); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	out, err := h.service.UpdateCourse(id, &in)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"message": "updated", "data": out})
}

// DeleteCourse — PATCH /courses/delete/:id
func (h *Handler) DeleteCourse(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteCourse(id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"message": "deleted_at"})
}

// FindCoursesbyFilter — GET /courses?coursename=&prophet_name=&duration=
func (h *Handler) FindCoursesByFilter(c *fiber.Ctx) error {
	courseName := c.Query("coursename")
	prophetName := c.Query("prophetname")
	duration := c.Query("duration")

	filter := map[string]interface{}{}
	if courseName != "" {
		filter["coursename"] = courseName
	}
	if prophetName != "" {
		filter["prophetname"] = prophetName
	}
	if duration != "" {
		filter["duration"] = duration
	}

	courses, err := h.service.FindCoursesByFilter(filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(courses) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No courses found matching the filter",
			"data":    []interface{}{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": len(courses),
		"data":  courses,
	})
}
