package crawl

//抓取状态.0-创建任务;1-M3U8 URL抓取中;2-M3U8 URL抓取失败;3-M3U8 URL抓取完成;4-M3U8下载中;5-M3U8下载异常;6-M3U8下载结束';
const (
	Init        = 0
	Crawling    = 1
	CrawlErr    = 2
	CrawlFinish = 3
	M3U8Parsing = 4
	M3U8Err     = 5
	M3U8Parsed  = 6
)

const (
	// 抓取类型.1-页面URL;2-文件m3u8;3-MP4地址
	QueueTypePageUrl  = 1
	QueueTypeM3u8File = 2
	QueueTypeMp4Url   = 3
)

const (
	TypePageUrl = 1
	TypeM3U8Url = 2
	TypeMP4Url  = 3
)

const (
	CrawProxyClose = 0
	CrawProxyOpen  = 1
)

const (
	BusinessTypeNormal     = 0
	BusinessTypeCrawlLogin = 1 // 抓需要登录的资源
	BusinessTypeNiVod      = 2 // NiVod抓取模式
	BusinessTypeBananTV    = 3 // BananTV抓取模式
)

const (
	CrawlM3U8NotifyNo  = 0
	CrawlM3U8NotifyYes = 1
)

const (
	CrawlTVInit       = 0 // INIT
	CrawlTVPadInfo    = 1 // 自动补全视频信息中
	CrawlTVPadInfoErr = 2 // 补充视频信息失败
	CrawlTVPadInfoOK  = 3 // 补充视频信息成功
	CrawlTVPadId      = 4 // 补充TV ID信息中
	CrawlTVPadIdErr   = 5 // 补充TV ID信息失败
	CrawlTVPadIdOk    = 6 // 补充TV ID信息成功
)

const (
	CrawlTVItemInit     = 0 // INIT
	CrawlTVItemPadId    = 1 // 补充TV Item ID信息中
	CrawlTVItemPadIdErr = 2 // 补充TV Item ID信息失败
	CrawlTVItemPadIdOk  = 3 // 补充TV Item ID信息成功
)
