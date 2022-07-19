package driving

import (
	"github.com/gin-gonic/gin"
	"github.com/hex_microservice_template/internal/core/domain"
	"github.com/hex_microservice_template/internal/core/ports/inbound"
	"github.com/hex_microservice_template/internal/core/usecase"
	"github.com/pkg/errors"
	"net/http"
	"sync"
)

type PingResponse struct {
	Name    string
	Version string
	Built   string
	Status  int
}

type Response struct {
	Code    int         `json:"Code"`
	Status  string      `json:"Status"`
	Message string      `json:"Message,omitempty"`
	ErrCode string      `json:"ErrCode,omitempty"`
	Data    interface{} `json:"Data,omitempty"`
}

type Server struct {
	shutdownChan chan bool

	router *gin.Engine
	wg     sync.WaitGroup

	AppName   string
	Version   string
	BuildTime string
}

func InitServer(handler inbound.RedirectHandler) *Server {
	server := &Server{
		router:       gin.New(),
		shutdownChan: make(chan bool),
	}

	server.router.Use(gin.Logger())
	server.router.Use(gin.Recovery())
	server.router.Use(CORS())

	server.router.GET(Ping, func(context *gin.Context) {
		context.JSON(http.StatusOK, PingResponse{
			Status:  http.StatusOK,
			Name:    server.AppName,
			Version: server.Version,
			Built:   server.BuildTime,
		})
	})

	server.router.GET(Get, handler.Get)
	server.router.POST(Post, handler.Post)

	return server
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) ListenAndServe(address string) error {
	return s.router.Run(address)
}

func (s *Server) Close() {
	close(s.shutdownChan)
	s.wg.Wait()
}

type handler struct {
	redirectService inbound.RedirectService
}

func NewHandler(redirectService inbound.RedirectService) inbound.RedirectHandler {
	return &handler{redirectService: redirectService}
}

func (h *handler) Get(context *gin.Context) {
	code := context.Param("code")
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == usecase.ErrRedirectNotFound {
			context.JSON(http.StatusNotFound, Response{
				Code:    http.StatusNotFound,
				Status:  http.StatusText(http.StatusNotFound),
				Message: "Requested URL not found",
				ErrCode: usecase.ErrRedirectNotFound.Error(),
			})
			return
		}
		context.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
		return
	}
	context.JSON(http.StatusMovedPermanently, Response{
		Code:    http.StatusMovedPermanently,
		Status:  http.StatusText(http.StatusMovedPermanently),
		Message: "Success",
		Data:    redirect.URL,
	})
	return
}

func (h *handler) Post(context *gin.Context) {
	var redirectInput domain.Redirect
	err := context.ShouldBindJSON(&redirectInput)
	if err != nil {
		context.JSON(http.StatusBadRequest, Response{
			Code:    http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		})
		return
	}

	err = h.redirectService.Store(&redirectInput)
	if err != nil {
		if errors.Cause(err) == usecase.ErrRedirectInvalid {
			context.JSON(http.StatusBadRequest, Response{
				Code:    http.StatusBadRequest,
				Status:  http.StatusText(http.StatusBadRequest),
				Message: err.Error(),
			})
			return
		}
		context.JSON(http.StatusInternalServerError, Response{
			Code:    http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Status:  http.StatusText(http.StatusCreated),
		Message: "Success",
		Data:    redirectInput,
	})
	return
}
