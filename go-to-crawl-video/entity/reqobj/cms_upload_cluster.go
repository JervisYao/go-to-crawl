package reqobj

import "go-to-crawl-common/app/utils/reqobj"

type CmsUploadClusterCreate struct {
	CreateUser       int    `p:"createUser" v:"required#登录信息过期"`
	ClusterNo        string `p:"clusterNo" v:"required#请指定一个非重复的集群编码"`
	ClusterIps       string `p:"clusterIps" v:"required#请指定集群IP列表,用英文逗号隔开"`
	ClusterAvailable int    `p:"clusterAvailable" v:"required#请设置集群是否可用"`
	ClusterNote      int    `p:"clusterNote"`
}

type CmsUploadClusterQry struct {
	*reqobj.CmsBasePageQry
	ClusterNo string `p:"clusterNo"`
}

type CmsUploadClusterUpdate struct {
	Id               int `p:"id" v:"required#未指定集群ID"`
	ClusterAvailable int `p:"clusterAvailable" v:"required#请设置集群是否可用"`
}
