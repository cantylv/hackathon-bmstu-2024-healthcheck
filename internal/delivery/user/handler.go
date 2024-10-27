package user

import (
	"errors"
	"net/http"

	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity/dto"
	ucUser "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/usecase/user"
	f "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/functions"
	mc "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myconstants"
	me "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myerrors"
	"go.uber.org/zap"
)

type UserHandlerManager struct {
	ucUser ucUser.Usecase
	logger *zap.Logger
}

// NewUserHandlerManager возвращает менеджер хендлеров, отвечающих за создание/удаление пользователя из системы
func NewUserHandlerManager(ucUser ucUser.Usecase, logger *zap.Logger) *UserHandlerManager {
	return &UserHandlerManager{
		ucUser: ucUser,
		logger: logger,
	}
}

// Read метод чтения данных пользователя, в случае успеха возвращает пользователю его данные.
func (h *UserHandlerManager) Read(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	username := f.GetUsernameCtx(r)
	if username == "" {
		h.logger.Info(me.ErrNotAuthenticated.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrNotAuthenticated.Error()}, http.StatusUnauthorized)
		return
	}
	if err := dto.ValidateUsername(username); err != nil {
		h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	u, err := h.ucUser.Read(r.Context(), username)
	if err != nil {
		if errors.Is(err, me.ErrUserNotExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}
	f.Response(w, getUserWithoutPassword(u), http.StatusOK)
}

// Delete метод удаление пользователя, в случае успеха возвращает сообщение о том, что пользователь был удален.
func (h *UserHandlerManager) Delete(w http.ResponseWriter, r *http.Request) {
	requestID, err := f.GetCtxRequestID(r)
	if err != nil {
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	}
	username := f.GetUsernameCtx(r)
	if username == "" {
		h.logger.Info(me.ErrNotAuthenticated.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrNotAuthenticated.Error()}, http.StatusUnauthorized)
		return
	}
	err = h.ucUser.Delete(r.Context(), username)

	if err != nil {
		if errors.Is(err, me.ErrUserNotExist) {
			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
			return
		}
		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
		return
	}

	f.FlashCookie(w, r)
	f.Response(w, dto.ResponseDetail{Detail: "Вы успешно удалили себя из приложения 'Healthcheck'"}, http.StatusOK)
}
