package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Edbeer/Project/config"
	"github.com/Edbeer/Project/internal/entity"
	mockservice "github.com/Edbeer/Project/internal/service/mock"
	"github.com/Edbeer/Project/pkg/converter"
	"github.com/Edbeer/Project/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockservice.NewMockUser(ctrl)
	mockSessionService := mockservice.NewMockSession(ctrl)

	config := &config.Config{
		Cookie: config.Cookie {
			MaxAge: 10,
		},
	}

	userHandler := NewUserHandler(config, mockUserService, mockSessionService)

	user := &entity.User{
		Name: "PavelV",
		Email: "edbeermtn@gmail.com",
		Password: "12345678",
	}

	buffer, err := converter.AnyToBytesBuffer(user)
	require.NoError(t, err)
	require.NotNil(t, buffer)
	require.Nil(t, err)

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/user/sign-up", strings.NewReader(buffer.String()))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	c := e.NewContext(request, recorder)
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userHandler.SignUp")
	defer span.Finish()

	handlerFunc := userHandler.SignUp()

	userID := uuid.New()
	userWithToken := &entity.UserWithToken{
		User: &entity.User{
			ID: userID,
		},
	}
	sess := &entity.Session{
		UserID: userID,
	}
	token := "token"

	mockUserService.EXPECT().SignUp(ctxWithTrace, gomock.Eq(user)).Return(userWithToken, nil)
	mockSessionService.EXPECT().CreateSession(ctxWithTrace, gomock.Eq(sess), 10).Return(token, nil)

	err = handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestHandler_SignIn(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockservice.NewMockUser(ctrl)
	mockSessionService := mockservice.NewMockSession(ctrl)

	config := &config.Config{
		Cookie: config.Cookie {
			MaxAge: 10,
		},
	}

	userHandler := NewUserHandler(config, mockUserService, mockSessionService)

	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}

	login := &Login{
		Email:    "edbeermtn@gmail.com",
		Password: "12345678",
	}

	user := &entity.User{
		Email:    login.Email,
		Password: login.Password,
	}

	buffer, err := converter.AnyToBytesBuffer(user)
	require.NoError(t, err)
	require.NotNil(t, buffer)
	require.Nil(t, err)

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/user/sign-in", strings.NewReader(buffer.String()))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()

	c := e.NewContext(request, recorder)
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userHandler.SignIn")
	defer span.Finish()

	handlerFunc := userHandler.SignIn()

	userID := uuid.New()
	userWithToken := &entity.UserWithToken{
		User: &entity.User{
			ID: userID,
		},
	}
	sess := &entity.Session{
		UserID: userID,
	}
	token := "refresh token"

	mockUserService.EXPECT().SignIn(ctxWithTrace, gomock.Eq(user)).Return(userWithToken, nil)
	mockSessionService.EXPECT().CreateSession(ctxWithTrace, gomock.Eq(sess), 10).Return(token, nil)

	err = handlerFunc(c)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestHandler_SignOut(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mockservice.NewMockUser(ctrl)
	mockSessionService := mockservice.NewMockSession(ctrl)

	config := &config.Config{
		Cookie: config.Cookie {
			MaxAge: 10,
		},
	}

	userHandler := NewUserHandler(config, mockUserService, mockSessionService)
	token := "jwt-token"
	cookieValue := "cookieValue"

	e := echo.New()
	request := httptest.NewRequest(http.MethodPost, "/api/user/sign-out", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.AddCookie(&http.Cookie{Name: token, Value: cookieValue})
	recorder := httptest.NewRecorder()

	c := e.NewContext(request, recorder)
	ctx := utils.GetRequestCtx(c)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userHandler.SignOut")
	defer span.Finish()

	logout := userHandler.SignOut()

	cookie, err := request.Cookie(token)
	require.NoError(t, err)
	require.NotNil(t, cookie)
	require.NotEqual(t, cookie.Value, "")
	require.Equal(t, cookie.Value, cookieValue)

	mockSessionService.EXPECT().DeleteSession(ctxWithTrace, gomock.Eq(cookie.Value)).Return(nil)

	err = logout(c)
	require.NoError(t, err)
	require.Nil(t, err)
}