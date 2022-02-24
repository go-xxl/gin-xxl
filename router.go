package xxl

import (
	"github.com/gin-gonic/gin"
	sdk "github.com/go-xxl/xxl"
	"github.com/go-xxl/xxl/admin"
	"github.com/go-xxl/xxl/job"
	"github.com/go-xxl/xxl/log"
	"net/http"
)

type PathHandler struct {
	Path    string
	Handler gin.HandlerFunc
}

type GinRouter struct {
	handlers map[string][]PathHandler
	option   sdk.Options
	mw       []gin.HandlerFunc
}

// NewGinRouter create a gin router struct
func NewGinRouter(opt sdk.Options) *GinRouter {
	return &GinRouter{
		handlers: make(map[string][]PathHandler),
		option:   opt,
	}
}

// SetLog set xxl log
func (r *GinRouter) SetLog(l log.Logger) {
	log.SetLog(l)
}

// Router gin group func
func (r *GinRouter) Router(group *gin.RouterGroup) {

	group.Use(r.mw...)

	for method, handlers := range r.handlers {
		for _, pathHandler := range handlers {
			group.Handle(method, pathHandler.Path, pathHandler.Handler)
		}
	}

	biz := sdk.NewExecutorService()
	group.Any("/run", convertGinFunc(biz.Run))
	group.Any("/kill", convertGinFunc(biz.Kill))
	group.Any("/log", convertGinFunc(biz.Log))
	group.Any("/beat", convertGinFunc(biz.Beat))
	group.Any("/idleBeat", convertGinFunc(biz.IdleBeat))
}

// POST add post handler
func (r *GinRouter) POST(path string, handler gin.HandlerFunc) {
	r.handlers[http.MethodPost] = append(r.handlers[http.MethodPost], PathHandler{
		Path:    path,
		Handler: handler,
	})
}

// GET add post handler
func (r *GinRouter) GET(path string, handler gin.HandlerFunc) {
	r.handlers[http.MethodGet] = append(r.handlers[http.MethodGet], PathHandler{
		Path:    path,
		Handler: handler,
	})
}

// PUT add post handler
func (r *GinRouter) PUT(path string, handler gin.HandlerFunc) {
	r.handlers[http.MethodPut] = append(r.handlers[http.MethodPut], PathHandler{
		Path:    path,
		Handler: handler,
	})
}

// PATCH add post handler
func (r *GinRouter) PATCH(path string, handler gin.HandlerFunc) {
	r.handlers[http.MethodPatch] = append(r.handlers[http.MethodPatch], PathHandler{
		Path:    path,
		Handler: handler,
	})
}

// DELETE add post handler
func (r *GinRouter) DELETE(path string, handler gin.HandlerFunc) {
	r.handlers[http.MethodDelete] = append(r.handlers[http.MethodDelete], PathHandler{
		Path:    path,
		Handler: handler,
	})
}

// CONNECT add post handler
func (r *GinRouter) CONNECT(path string, handler gin.HandlerFunc) {
	r.handlers[http.MethodConnect] = append(r.handlers[http.MethodConnect], PathHandler{
		Path:    path,
		Handler: handler,
	})
}

// OPTIONS add post handler
func (r *GinRouter) OPTIONS(path string, handler gin.HandlerFunc) {
	r.handlers[http.MethodOptions] = append(r.handlers[http.MethodOptions], PathHandler{
		Path:    path,
		Handler: handler,
	})
}

// TRACE add post handler
func (r *GinRouter) TRACE(path string, handler gin.HandlerFunc) {
	r.handlers[http.MethodTrace] = append(r.handlers[http.MethodTrace], PathHandler{
		Path:    path,
		Handler: handler,
	})
}

// Register register node
func (r *GinRouter) Register() {
	r.getAdm().Register()
}

// Deregister deregister node
func (r *GinRouter) Deregister() {
	r.getAdm().RegistryRemove()
}

// getAdm
func (r *GinRouter) getAdm() *admin.AdmApi {
	adm := admin.NewAdmApi()

	adm.SetOpt(
		admin.SetAccessToken(r.option.AccessToken),
		admin.SetAdmAddresses(r.option.AdmAddresses),
		admin.SetTimeout(r.option.Timeout),
		admin.SetRegistryKey(r.option.RegistryKey),
		admin.SetRegistryValue(r.option.RegistryValue),
		admin.SetRegistryGroup(admin.GetGroupName(admin.EXECUTOR)),
	)

	return adm
}

func (r *GinRouter) Job(handlerName string, handler job.Func) {
	job.GetAllHandlerList().Set(handlerName, job.Job{
		Id:        0,
		Name:      "",
		Ext:       nil,
		Cancel:    nil,
		Param:     nil,
		Fn:        handler,
		StartTime: 0,
		EndTime:   0,
	})
}

func (r *GinRouter) MiddleWare(mw ...gin.HandlerFunc) {
	r.mw = mw
}
