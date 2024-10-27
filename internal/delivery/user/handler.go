package user

import (
	"errors"
	"net/http"

	"github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/entity/dto"
	ucUser "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/usecase/user"
	f "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/functions"
	mc "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myconstants"
	me "github.com/cantylv/hackathon-bmstu-2024-healthcheck/internal/utils/myerrors"
	"github.com/gorilla/mux"
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
	username := mux.Vars(r)["username"]
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
	// requestID, err := f.GetCtxRequestID(r)
	//
	//	if err != nil {
	//		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	//	}
	//
	// userEmail := mux.Vars(r)["email"]
	//
	//	if !govalidator.IsEmail(userEmail) {
	//		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
	//		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
	//		return
	//	}
	//
	// userEmailDelete := mux.Vars(r)["email_delete"]
	//
	//	if !govalidator.IsEmail(userEmailDelete) {
	//		h.logger.Info(me.ErrInvalidEmail.Error(), zap.String(mc.RequestID, requestID))
	//		f.Response(w, dto.ResponseError{Error: me.ErrInvalidEmail.Error()}, http.StatusBadRequest)
	//		return
	//	}
	//
	// err = h.ucUser.Delete(r.Context(), userEmail, userEmailDelete)
	//
	//	if err != nil {
	//		if errors.Is(err, me.ErrUserNotExist) || errors.Is(err, me.ErrUserIsResponsible) {
	//			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
	//			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusBadRequest)
	//			return
	//		}
	//		if errors.Is(err, me.ErrOnlyRootCanDeleteUser) || errors.Is(err, me.ErrCantDeleteRoot) {
	//			h.logger.Info(err.Error(), zap.String(mc.RequestID, requestID))
	//			f.Response(w, dto.ResponseError{Error: err.Error()}, http.StatusForbidden)
	//			return
	//		}
	//		h.logger.Error(err.Error(), zap.String(mc.RequestID, requestID))
	//		f.Response(w, dto.ResponseError{Error: me.ErrInternal.Error()}, http.StatusInternalServerError)
	//		return
	//	}
	//
	// f.Response(w, dto.ResponseDetail{Detail: "user was succesful deleted"}, http.StatusOK)
}
