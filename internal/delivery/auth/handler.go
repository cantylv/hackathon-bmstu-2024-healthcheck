package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity/dto"
	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/usecase/auth"
	f "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/functions"
	mc "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myconstants"
	me "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myerrors"
	"go.uber.org/zap"
)

type AuthHandlerManager struct {
	ucAuth auth.Usecase
	logger *zap.Logger
}

// NewUserHandlerManager возвращает менеджер хендлеров, отвечающих за создание/удаление пользователя из системы
func NewAuthHandlerManager(ucAuth auth.Usecase, logger *zap.Logger) *AuthHandlerManager {
	return &AuthHandlerManager{
		ucAuth: ucAuth,
		logger: logger,
	}
}

func (h *AuthHandlerManager) SignUp(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	jwtToken, err := f.GetJWtToken(r)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		h.logger.Error(fmt.Sprintf("error while getting jwt token: %v", err), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusInternalServerError)
		return
	}
	if jwtToken != "" {
		h.logger.Info("user is already registered", zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrAlreadyRegistered.Error()}, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}

	var signForm dto.CreateData
	err = json.Unmarshal(body, &signForm)
	if err != nil {
		h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}
	err = signForm.Validate()
	if err != nil {
		h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	u, err := h.ucAuth.SignUp(r.Context(), &signForm)
	if err != nil {
		if errors.Is(err, me.ErrUserAlreadyExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}

	w, err = f.SetCookieAndHeaders(w, u.Username)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}

	f.Response(w, getUserWithoutPassword(u), http.StatusOK)
}

func (h *AuthHandlerManager) SignIn(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}

	jwtToken, err := f.GetJWtToken(r)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		h.logger.Error(fmt.Sprintf("error while getting jwt token: %v", err), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusInternalServerError)
		return
	}
	if jwtToken != "" {
		h.logger.Info("user is already registered", zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrAlreadyRegistered.Error()}, http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}
	var signForm dto.AuthData
	err = json.Unmarshal(body, &signForm)
	if err != nil {
		h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInvalidData.Error()}, http.StatusBadRequest)
		return
	}
	err = signForm.Validate()
	if err != nil {
		h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	u, err := h.ucAuth.SignIn(r.Context(), &signForm)
	if err != nil {
		if errors.Is(err, me.ErrIncorrectPwdOrLogin) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}

	w, err = f.SetCookieAndHeaders(w, u.Username)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, getUserWithoutPassword(u), http.StatusOK)
}

func (h *AuthHandlerManager) SignOut(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	jwtToken, err := f.GetJWtToken(r)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		h.logger.Error(fmt.Sprintf("error while getting jwt token: %v", err), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	if jwtToken == "" {
		h.logger.Info("user is not authenticated", zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrNotAuthenticated.Error()}, http.StatusUnauthorized)
		return
	}

	f.FlashCookie(w, r)
	f.Response(w, dto.ResponseDetail{Detail: "ok"}, http.StatusOK)
}
