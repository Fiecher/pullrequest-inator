package api

import (
	"errors"
	"net/http"
	"pullrequest-inator/internal/api/dtos"
	"pullrequest-inator/internal/infrastructure/encoding"
	"pullrequest-inator/internal/infrastructure/services"

	"github.com/labstack/echo/v4"
)

type Server struct {
	prService   *services.PullRequestService
	teamService *services.TeamService
	userService *services.UserService
}

func NewServer(prService *services.PullRequestService, teamService *services.TeamService, userService *services.UserService) (*Server, error) {
	if prService == nil {
		return nil, errors.New("prService is required")
	}
	if teamService == nil {
		return nil, errors.New("teamService is required")
	}
	if userService == nil {
		return nil, errors.New("userService is required")
	}

	return &Server{
		prService:   prService,
		teamService: teamService,
		userService: userService,
	}, nil
}

func (s *Server) PostPullRequestCreate(ctx echo.Context) error {
	var input PostPullRequestCreateJSONRequestBody
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request",
				"details": err.Error(),
			},
		})
	}

	dtoReq := &dtos.PullRequest{
		PullRequestId:   input.PullRequestId,
		PullRequestName: input.PullRequestName,
		AuthorId:        input.AuthorId,
	}

	pr, err := s.prService.CreatePullRequest(ctx.Request().Context(), dtoReq)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, map[string]any{
		"pr": ToAPIPullRequest(*pr),
	})
}

func (s *Server) PostPullRequestMerge(ctx echo.Context) error {
	var input PostPullRequestMergeJSONRequestBody
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request",
				"details": err.Error(),
			},
		})
	}

	prID := encoding.DecodeID(input.PullRequestId)
	pr, err := s.prService.MarkAsMerged(ctx.Request().Context(), prID)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	if pr == nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    "INTERNAL",
				"message": "merge failed: nil PR",
			},
		})
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"pr": ToAPIPullRequest(*pr),
	})
}

func (s *Server) PostPullRequestReassign(ctx echo.Context) error {
	var input PostPullRequestReassignJSONRequestBody
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request",
				"details": err.Error(),
			},
		})
	}

	prID := encoding.DecodeID(input.PullRequestId)
	oldID := encoding.DecodeID(input.OldUserId)

	resp, err := s.prService.ReassignReviewer(ctx.Request().Context(), oldID, prID)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	if resp == nil || resp.Pr.PullRequestId == "" {
		return ctx.JSON(http.StatusInternalServerError, map[string]any{
			"error": map[string]string{
				"code":    "INTERNAL",
				"message": "failed to reassign: empty response",
			},
		})
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"pr":          ToAPIPullRequest(resp.Pr),
		"replaced_by": resp.ReplacedBy,
	})
}

func (s *Server) PostTeamAdd(ctx echo.Context) error {
	var team Team
	if err := ctx.Bind(&team); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request",
				"details": err.Error(),
			},
		})
	}

	dtoTeam := FromAPITeam(team)
	err := s.teamService.CreateTeamWithUsers(ctx.Request().Context(), &dtoTeam)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusCreated, map[string]any{
		"team": team,
	})
}

func (s *Server) GetTeamGet(ctx echo.Context, params GetTeamGetParams) error {
	team, err := s.teamService.GetTeamByName(ctx.Request().Context(), params.TeamName)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, ToAPITeam(*team))
}

func (s *Server) PostUsersSetIsActive(ctx echo.Context) error {
	var input PostUsersSetIsActiveJSONRequestBody
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]any{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "invalid request",
				"details": err.Error(),
			},
		})
	}

	userID := encoding.DecodeID(input.UserId)
	updated, err := s.teamService.SetUserActiveByID(ctx.Request().Context(), userID, input.IsActive)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]User{
		"user": ToAPIUser(*updated),
	})
}

func (s *Server) GetUsersGetReview(ctx echo.Context, params GetUsersGetReviewParams) error {
	userID := encoding.DecodeID(params.UserId)
	resp, err := s.prService.FindPullRequestsByReviewer(ctx.Request().Context(), userID)
	if err != nil {
		return mapAppErrorToEchoResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"user_id":       params.UserId,
		"pull_requests": ToAPIPullRequestShortList(resp),
	})
}

func (s *Server) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"status": "OK",
	})
}

func mapAppErrorToEchoResponse(ctx echo.Context, err error) error {
	code := http.StatusInternalServerError
	msg := "internal server error"
	apiCode := "INTERNAL"

	switch {
	case errors.Is(err, services.ErrPRAlreadyExists):
		code = http.StatusConflict
		msg = "pull request already exists"
		apiCode = "PR_EXISTS"
	case errors.Is(err, services.ErrPRAlreadyMerged):
		code = http.StatusConflict
		msg = "pull request already merged"
		apiCode = "PR_MERGED"
	case errors.Is(err, services.ErrNoReviewCandidates):
		code = http.StatusConflict
		msg = "no active users to assign as reviewers"
		apiCode = "NO_CANDIDATE"
	case errors.Is(err, services.ErrUserNotReviewer):
		code = http.StatusBadRequest
		msg = "user is not a reviewer"
		apiCode = "NOT_ASSIGNED"
	case errors.Is(err, services.ErrPRNotFound):
		code = http.StatusNotFound
		msg = "pull request not found"
		apiCode = "NOT_FOUND"
	case errors.Is(err, services.ErrTeamNotFound), errors.Is(err, services.ErrAuthorNotFound):
		code = http.StatusNotFound
		msg = err.Error()
		apiCode = "NOT_FOUND"
	case errors.Is(err, services.ErrTeamExists):
		code = http.StatusConflict
		msg = "team already exists"
		apiCode = "TEAM_EXISTS"
	}

	return ctx.JSON(code, map[string]any{
		"error": map[string]string{
			"code":    apiCode,
			"message": msg,
			"details": err.Error(),
		},
	})
}
