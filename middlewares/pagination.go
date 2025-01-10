package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mcuadros/go-defaults"
	"github.com/nft-rainbow/rainbow-goutils/utils/ginutils"
)

type Pagination struct {
	Page  int `json:"page" form:"page" default:"1"`
	Limit int `json:"limit" form:"limit" default:"10"`
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

func PaginationMiddleware(c *gin.Context) {
	var pagination Pagination
	defaults.SetDefaults(&pagination)

	if err := c.ShouldBindQuery(&pagination); err != nil {
		ginutils.RenderError(c, ginutils.NewBadRequestNormalGinError(err.Error()))
		return
	}

	c.Set("page", pagination.Page)
	c.Set("limit", pagination.Limit)
	c.Set("offset", pagination.Offset())
	c.Set("pagination", pagination)
	c.Next()
}
