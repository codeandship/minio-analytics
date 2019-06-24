package http

import (
	"net/http"

	minioanalytics "github.com/codeandship/minio-analytics"
	"github.com/gin-gonic/gin"
)

type API struct {
	address string
	r       *gin.Engine
	store   minioanalytics.Storage
	uam     *minioanalytics.UserAgentMatcher
}

func NewAPI(addr string, s minioanalytics.Storage, uam *minioanalytics.UserAgentMatcher) *API {
	a := &API{
		address: addr,
		store:   s,
		uam:     uam,
	}
	return a
}

func (a *API) Open() error {
	a.setupRouter()
	return a.r.Run(a.address)
}

func (a *API) Close() error {
	return nil
}

func (a *API) setupRouter() {
	gin.SetMode(gin.ReleaseMode)
	a.r = gin.Default()

	a.r.GET("/analytics", a.handleGetAnalytics)
}

func (a *API) handleGetAnalytics(c *gin.Context) {
	res, err := a.store.ListAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for i, r := range res {
		r.UserAgentCount = a.uam.MatchMap(r.UserAgentCount)
		res[i] = r
	}

	c.JSON(http.StatusOK, res)
}
