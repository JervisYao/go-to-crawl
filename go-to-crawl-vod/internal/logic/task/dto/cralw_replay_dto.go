package dto

type CrawlReplayInterface interface {
	CreateProgram(replayConfig *model.CmsCrawlReplayConfig, manifestTask *model.CmsCrawlReplayManifestTask)
}

type AbstractCrawlReplayUrl struct {
	CrawlReplayInterface
}

func (receiver *AbstractCrawlReplayUrl) CreateProgram(replayConfig *model.CmsCrawlReplayConfig, manifestTask *model.CmsCrawlReplayManifestTask) {

}
