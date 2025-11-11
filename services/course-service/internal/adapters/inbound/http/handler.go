package http

import (
	"log"
	"strconv"
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
	// Course
	group := router.Group("/api")
	group.Post("/courses", h.CreateCourse)
	group.Get("/courses/popular", h.ListPopularCourses)
	group.Get("/courses/:id", h.GetCourseByID)
	group.Get("/prophets/:prophetID/courses", h.ListCoursesByProphet)
	group.Patch("/courses/:id", h.UpdateCourse)
	group.Patch("/courses/delete/:id", h.DeleteCourse)
	group.Get("/courses", h.FindCoursesByFilter)
	group.Get("/courses/prophet/courses", h.ListCurrentProphetCourses)
	// Review
	group.Post("/courses/:courseID/review", h.CreateReview)
	group.Get("/courses/review/:id", h.GetReviewByID)
	group.Get("/courses/:courseID/reviews", h.GetReviewByCourseID)
}

func (h *Handler) ListPopularCourses(c *fiber.Ctx) error {
	limit := c.Query("limit")
	if limit == "" {
		limit = "10"
	}
	limitInt, err := strconv.Atoi(limit)
	log.Println("limitInt", limitInt)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid limit")
	}
	courses, err := h.service.ListPopularCourses(c.Context(), limitInt)
	if err != nil {
		log.Println("error", err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(courses)
}

// CreateCourse — POST /courses
func (h *Handler) CreateCourse(c *fiber.Ctx) error {
	userRole := c.Get("X-User-Role")
	if userRole != "prophet" {
		return fiber.NewError(fiber.StatusForbidden, "only prophets can create courses")
	}

	prophetID := c.Get("X-User-ID")
	var req struct {
		CourseName  string  `json:"coursename"`
		CourseType  string  `json:"coursetype"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Duration    int32   `json:"duration"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	input := app.CreateCourseInput{
		ProphetID:   prophetID,
		CourseName:  req.CourseName,
		CourseType:  domain.CourseType(req.CourseType),
		Description: req.Description,
		Price:       req.Price,
		Duration:    domain.DurationEnum(req.Duration),
	}

	course, err := h.service.CreateCourse(c.Context(), input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(course)
}

// GetCourseByID — GET /courses/:id
func (h *Handler) GetCourseByID(c *fiber.Ctx) error {
	id := c.Params("id")

	course, err := h.service.GetCourseDetailByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "course not found")
	}

	return c.JSON(course)
}

// ListCoursesByProphet — GET /prophets/:prophetID/courses
func (h *Handler) ListCoursesByProphet(c *fiber.Ctx) error {
	prophetID := c.Params("prophetID")

	courses, err := h.service.ListCoursesByProphet(c.Context(), prophetID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := struct {
		Timestamp time.Time                       `json:"timestamp"`
		Data      []*domain.CourseWithProphetName `json:"courses"`
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
	out, err := h.service.UpdateCourse(c.Context(), id, &in)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"message": "updated", "data": out})
}

// DeleteCourse — PATCH /courses/delete/:id
func (h *Handler) DeleteCourse(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteCourse(c.Context(), id); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"message": "deleted_at"})
}

// FindCoursesbyFilter — GET /courses?coursename=&prophet_name=&duration=
// GetAllCourses — GET /courses
func (h *Handler) FindCoursesByFilter(c *fiber.Ctx) error {
	searchTerm := c.Query("searchterm")
	duration := c.Query("duration")
	courseType := c.Query("coursetype")
	sortBy := c.Query("sortby")
	order := c.Query("order")

	filter := app.CourseFilter{
		SearchTerm: searchTerm,
		Duration:   duration,
		CourseType: app.ParseCourseType(courseType),
	}

	sort := app.CourseSort{
		SortBy: app.ParseSortType(sortBy),
		Order:  order,
	}

	courses, err := h.service.FindCoursesByFilter(c.Context(), filter, sort)
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

// CreateReview — POST /courses/:courseID/review
func (h *Handler) CreateReview(c *fiber.Ctx) error {
	var req struct {
		CustomerID   string  `json:"customer_id"`
		CustomerName string  `json:"customername"`
		Score        float64 `json:"score"`
		Title        string  `json:"title"`
		Description  string  `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Validation
	if req.CustomerID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "customer_id is required")
	}
	if req.CustomerName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "customername is required")
	}
	if req.Score < 0 || req.Score > 5 {
		return fiber.NewError(fiber.StatusBadRequest, "score must be between 0 and 5")
	}

	input := app.CreateReviewInput{
		CourseID:     c.Params("courseId"),
		CustomerID:   req.CustomerID,
		CustomerName: req.CustomerName,
		Score:        req.Score,
		Title:        req.Title,
		Description:  req.Description,
		DeletedAt:    false,
	}

	review, err := h.service.CreateReview(c.Context(), input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(review)
}

// GetReviewByID — GET /courses/review/:id
func (h *Handler) GetReviewByID(c *fiber.Ctx) error {
	id := c.Params("id")

	review, err := h.service.GetReviewByID(c.Context(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "review not found")
	}

	return c.JSON(review)
}

// GetReviewByCourseID — GET /courses/:courseId/reviews
func (h *Handler) GetReviewByCourseID(c *fiber.Ctx) error {
	courseId := c.Params("courseId")

	reviews, err := h.service.ListReviewsByCourse(c.Context(), courseId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"timestamp": time.Now(),
		"course_id": courseId,
		"count":     len(reviews),
		"reviews":   reviews,
	}

	return c.JSON(response)
}

// ListCurrentProphetCourses — GET /courses/prophet/courses
func (h *Handler) ListCurrentProphetCourses(c *fiber.Ctx) error {
	prophetID := c.Get("X-User-Id")

	courses, err := h.service.ListCoursesByProphet(c.Context(), prophetID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(courses)
}
