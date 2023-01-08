package ffmpeg

const (
	LineTypeComment = 1 // 注解行
	LineTypeKey     = 2 // 解密KEY行
	LineTypeSrc     = 3 // TS资源行
	LineTypeXMap    = 4 // #EXT-X-MAP行（所有ts统一的文件头，压缩体积所需）
)

const (
	SrcTypeNormal   = 1 // 普通资源类型 eg: .ts后缀
	SrcTypeImg      = 2 // 图片资源类型 eg: .jpg后缀
	SrcTypeNoSuffix = 3 // 无后缀 eg: http://xxx.com/abc
)

type M3u8DO struct {
	Schema         string
	Host           string
	FromUrl        string
	FromBaseUrl    string
	FromFile       string
	FromDir        string
	StreamLineList []StreamLineDO
	PngHeaderSize  int64
	MP4SaveFile    string // 下载MP4后临时变量

	MapSaveFile   string // 统一TS文件头下载路径临时变量
	MapOriginName string // URI="xxx"中的xxx

	KeySaveFile   string // 解密文件下载路径临时变量
	KeyOriginName string // URI="xxx"中的xxx
}

type StreamLineDO struct {
	LineType           int    // 行类型
	SrcType            int    // TS资源类型
	OriginLine         string // 原始行
	OriginTsSrcName    string // 原始TS资源名称
	TransformedLine    string // 转变后的行
	TransformedSrcName string // 转变后资源名
}
