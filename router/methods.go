package routes

import "github.com/gin-gonic/gin"

type renderer interface {
	RenderIndexPage(c *gin.Context)
	RenderSearchPage(c *gin.Context)
	RenderMTESTPage(c *gin.Context)
	RenderUserPage(c *gin.Context)
	RenderAdminPage(c *gin.Context)
}

type mtestHandlers interface {
	GetMTESTHandler(c *gin.Context)
	CreateMTESTHandler(c *gin.Context)
	UpdateMTESTHandler(c *gin.Context)
	DeleteMTESTHandler(c *gin.Context)
}

type userHandlers interface {
	UserCabinetHandler(c *gin.Context)
	ResetPasswordHandler(c *gin.Context)
	EditUserFieldHandler(c *gin.Context)
	PasswordCheckHandler(c *gin.Context)
	SetNewPasswordHandler(c *gin.Context)
	ActivateAccountHandler(c *gin.Context)
	RegistrationHandler(c *gin.Context)
	admActionsHandler
	regionHandler
	executorHandler
	governmentsHandler
	usersHandler
}

type admActionsHandler interface {
	GetAdministrativeActionsHandler(c *gin.Context)
	AddAdministrativeActionHandler(c *gin.Context)
	EditAdministrativeActionsHandler(c *gin.Context)
	DeleteAdministrativeActionsHandler(c *gin.Context)
}

type governmentsHandler interface {
	GetGovernmentsHandlers(c *gin.Context)
	EditGovernmentNameHandler(c *gin.Context)
	AddGovernmentHandler(c *gin.Context)
	RemoveGovernmentHandler(c *gin.Context)
}

type regionHandler interface {
	GetRegionsHandler(c *gin.Context)
	PostEditRegions(c *gin.Context)
}

type executorHandler interface {
	DeleteExecutorHandler(c *gin.Context)
	CreateMTESTExecutorHandler(c *gin.Context)
}

type userValidator interface {
	PasswordCheck(email, password string) bool
	CheckUserActivation(email string) bool
}

type usersHandler interface {
	GetUsersHandler(c *gin.Context)
}
