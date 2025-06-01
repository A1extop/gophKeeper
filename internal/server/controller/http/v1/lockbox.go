package v1

import (
	"github.com/gin-gonic/gin"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/domain"
	"gophKeeper/internal/server/middleware"
	"gophKeeper/internal/server/services/lockbox/models"
	"gophKeeper/internal/server/services/lockbox/usecase"
	"gophKeeper/util"
	"log"
	"net/http"
)

type LockBoxHandler struct {
	config         *config.Config
	lockBoxService usecase.ILockBoxUsecase
	mware          middleware.IMiddlewareService
}

func NewLockBoxHandlerHandler(config *config.Config, router *gin.RouterGroup, lockBoxService usecase.ILockBoxUsecase, mware middleware.IMiddlewareService) {
	lockBoxHandler := LockBoxHandler{
		config:         config,
		lockBoxService: lockBoxService,
		mware:          mware,
	}

	lockBoxRouter := router.Group("/lock_boxes")
	{
		lockBoxRouter.POST("/create", mware.MiddlewareJWT(), mware.AuthorizeRoles(util.Admin, util.Attendee), lockBoxHandler.createLockBox)
		lockBoxRouter.DELETE("/:name", mware.MiddlewareJWT(), mware.AuthorizeRoles(util.Admin, util.Attendee), lockBoxHandler.deleteLockBox)
		lockBoxRouter.GET("/:name", mware.MiddlewareJWT(), mware.AuthorizeRoles(util.Admin, util.Attendee), lockBoxHandler.getLockBox)
		lockBoxRouter.GET("/", mware.MiddlewareJWT(), mware.AuthorizeRoles(util.Admin, util.Attendee), lockBoxHandler.getLockBoxes)
		lockBoxRouter.POST("/create/update", mware.MiddlewareJWT(), mware.AuthorizeRoles(util.Admin, util.Attendee), lockBoxHandler.createOrUpdateLockBox)
		lockBoxRouter.PUT("/", mware.MiddlewareJWT(), mware.AuthorizeRoles(util.Admin, util.Attendee), lockBoxHandler.updateLockBox)

	}
}

func (l *LockBoxHandler) createLockBox(ctx *gin.Context) {
	var data models.Data

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := ctx.GetInt("userId")
	data.UserID = userId
	id, err := l.lockBoxService.CreateLock(ctx, &data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (l *LockBoxHandler) deleteLockBox(ctx *gin.Context) {
	name := ctx.Param("name")

	userId := ctx.GetInt("userId")
	err := l.lockBoxService.DeleteLock(ctx, name, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (l *LockBoxHandler) updateLockBox(ctx *gin.Context) {
	userId := ctx.GetInt("userId")
	var data models.Data
	if err := ctx.ShouldBindJSON(&data); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data.UserID = userId
	log.Println(data)
	err := l.lockBoxService.UpdateLock(ctx, &data)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)

}
func (l *LockBoxHandler) getLockBox(ctx *gin.Context) {
	name := ctx.Param("name")

	userId := ctx.GetInt("userId")
	lockBox, err := l.lockBoxService.GetLockByName(ctx, name, userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	log.Println(lockBox)
	ctx.JSON(http.StatusOK, lockBox)
}

func (l *LockBoxHandler) getLockBoxes(ctx *gin.Context) {

	userId := ctx.GetInt("userId")
	lockBoxes, err := l.lockBoxService.GetAllLocks(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(*lockBoxes) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": domain.ErrLockBoxNotFound.Error()})
		return
	}
	ctx.JSON(http.StatusOK, lockBoxes)

}

func (l *LockBoxHandler) createOrUpdateLockBox(ctx *gin.Context) {
	var data models.Data
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data.UserID = ctx.GetInt("userId")

	id, err := l.lockBoxService.CreateOrUpdateLock(ctx, &data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}
