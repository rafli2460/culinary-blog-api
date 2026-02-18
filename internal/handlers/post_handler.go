package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/rafli2460/culinary-blog-api/internal/service"
	"github.com/rafli2460/culinary-blog-api/pkg/response"
	"github.com/rs/zerolog/log"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostService(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(c fiber.Ctx) error {
	userIDVal := c.Locals("user_id")
	var userID int
	if idFloat, ok := userIDVal.(float64); ok {
		userID = int(idFloat)
	}

	title := c.FormValue("title")
	content := c.FormValue("content")

	file, err := c.FormFile("image")
	if err != nil {
		file = nil
	}

	err = h.postService.CreatePost(c.Context(), userID, title, content, file)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	log.Info().Int("user_id", userID).Str("title", title).Msg("New post successfully created")

	return response.Success(c, fiber.StatusCreated, "Post successfully published!", nil, nil)
}

func (h *PostHandler) DeletePost(c fiber.Ctx) error {
	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid post ID")
	}

	userIDVal := c.Locals("user_id")
	var currentUserID int
	if idFloat, ok := userIDVal.(float64); ok {
		currentUserID = int(idFloat)
	} else if idInt, ok := userIDVal.(int); ok {
		currentUserID = idInt
	}

	roleVal := c.Locals("role")
	var currentUserRole string
	if roleStr, ok := roleVal.(string); ok {
		currentUserRole = roleStr
	}

	err = h.postService.DeletePost(c.Context(), postID, currentUserID, currentUserRole)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return response.Error(c, fiber.StatusForbidden, err.Error())
		}

		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	log.Info().Int("post_id", postID).Int("deleted_by", currentUserID).Msg("post delete successfully")
	return response.Success(c, fiber.StatusOK, "post and image successfully deleted", nil, nil)
}

func (h *PostHandler) UpdatePost(c fiber.Ctx) error {
	postID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid post ID")
	}

	userIDVal := c.Locals("user_id")
	var currentUserID int
	if idFloat, ok := userIDVal.(float64); ok {
		currentUserID = int(idFloat)
	} else if idInt, ok := userIDVal.(int); ok {
		currentUserID = idInt
	}

	roleVal := c.Locals("role")
	var currentUserRole string
	if roleStr, ok := roleVal.(string); ok {
		currentUserRole = roleStr
	}

	title := c.FormValue("title")
	content := c.FormValue("content")

	file, _ := c.FormFile("image")

	err = h.postService.UpdatePost(c.Context(), postID, currentUserID, currentUserRole, title, content, file)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			return response.Error(c, fiber.StatusForbidden, err.Error())
		}
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	log.Info().Int("post_id", postID).Int("updated_by", currentUserID).Msg("Post successfully updated")

	return response.Success(c, fiber.StatusOK, "Post successfully updated!", nil, nil)
}

func (h *PostHandler) GetPost(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid post ID format")
	}

	post, err := h.postService.GetPost(c.Context(), id)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, err.Error())
	}

	return response.Success(c, fiber.StatusOK, "Post detail successfully retrieved", post, nil)
}

func (h *PostHandler) GetAllPosts(c fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	convPage, err := strconv.Atoi(page)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "please enter a number")
	}
	convLimit, err := strconv.Atoi(limit)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Page parameter must be a number")
	}

	posts, err := h.postService.GetAllPosts(c.Context(), convPage, convLimit)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Failed to retrieve posts")
	}

	return response.Success(c, fiber.StatusOK, "Posts successfully retrieved", posts, fiber.Map{
		"page":  page,
		"limit": limit,
		"count": len(posts),
	})
}
