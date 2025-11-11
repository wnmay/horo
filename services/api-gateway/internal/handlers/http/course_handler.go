package http_handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wnmay/horo/shared/env"
)

type CourseHandler struct {
	courseServiceURL string
	client           *http.Client
}

func NewCourseHandler() *CourseHandler {
	return &CourseHandler{
		courseServiceURL: env.GetString("COURSE_SERVICE_URL", "http://localhost:3005"),
		client:           &http.Client{},
	}
}

func (h *CourseHandler) CreateCourse(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "POST", h.courseServiceURL, "/api/courses")
}

func (h *CourseHandler) GetCourseByID(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "GET", h.courseServiceURL, fmt.Sprintf("/api/courses/%s", c.Params("id")))
}

func (h *CourseHandler) ListCoursesByProphet(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "GET", h.courseServiceURL, fmt.Sprintf("/api/prophets/%s/courses", c.Params("prophetId")))
}

func (h *CourseHandler) UpdateCourse(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "PATCH", h.courseServiceURL, fmt.Sprintf("/api/courses/%s", c.Params("id")))
}

func (h *CourseHandler) DeleteCourse(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "PATCH", h.courseServiceURL, fmt.Sprintf("/api/courses/delete/%s", c.Params("id")))
}

func (h *CourseHandler) FindCoursesByFilter(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "GET", h.courseServiceURL, "/api/courses")
}

func (h *CourseHandler) CreateReview(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "POST", h.courseServiceURL, fmt.Sprintf("/api/courses/%s/review", c.Params("courseId")))
}

func (h *CourseHandler) GetReviewByID(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "GET", h.courseServiceURL, fmt.Sprintf("/api/courses/review/%s", c.Params("id")))
}

func (h *CourseHandler) ListReviewsByCourse(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "GET", h.courseServiceURL, fmt.Sprintf("/api/courses/%s/reviews", c.Params("courseId")))
}

func (h *CourseHandler) ListCurrentProphetCourses(c *fiber.Ctx) error {
	return ProxyRequest(c, h.client, "GET", h.courseServiceURL, "/api/courses/prophet/courses")
}