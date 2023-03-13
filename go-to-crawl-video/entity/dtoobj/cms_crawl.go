package dtoobj

import (
	"easygoadmin/appnewcms/model"
)

type CmsCrawl struct {
	*model.CmsCrawlQueue
	ShowStatus   int    `json:"showStatus"`
	ResourcePath string `json:"resourcePath"`
}
