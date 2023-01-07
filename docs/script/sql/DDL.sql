CREATE TABLE `crawl_proxy`
(
    `id`            int(11)     NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `target_domain` varchar(50) NOT NULL COMMENT '目标网站顶级域名',
    `proxy_url`     varchar(50) NOT NULL COMMENT '代理地址',
    `proxy_status`  tinyint(4)  NOT NULL COMMENT '代理状态. 0-停用,1-使用中',
    `create_user`   varchar(20) DEFAULT NULL COMMENT '创建者',
    `create_time`   datetime    DEFAULT NULL COMMENT '创建时间',
    `update_user`   varchar(20) DEFAULT NULL COMMENT '修改者',
    `update_time`   datetime    DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB COMMENT ='抓取代理';

CREATE TABLE `crawl_queue`
(
    `id`                int(11)      NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `host_type`         int(11)               DEFAULT '0' COMMENT '主机类型。0-普通类型；1-抓付费资源类型',
    `host_ip`           varchar(50)           DEFAULT NULL COMMENT '任务处理的主机的IP。由哪台机器领取的M3U8下载任务就不能变更了',
    `country_code`      varchar(10)  NOT NULL COMMENT '国家二字码.(eg: CN,US,SG等)',
    `video_year`        int(11)      NOT NULL COMMENT '视频发布年份',
    `video_coll_id`     bigint(20)   NOT NULL DEFAULT '-1' COMMENT '视频集ID（视频集ID，不限于电视剧,-1代表单集视频，或者说电影）',
    `video_item_id`     bigint(20)            DEFAULT NULL COMMENT '视频集对应视频项ID（不限于电视剧的剧集）',
    `crawl_type`        tinyint(4)   NOT NULL DEFAULT '1' COMMENT '抓取类型.1-页面URL;2-文件m3u8;3-MP4地址',
    `crawl_status`      tinyint(4)   NOT NULL COMMENT '//抓取状态.0-创建任务;1-M3U8 URL抓取中;2-M3U8 URL抓取失败;3-M3U8 URL抓取完成;4-M3U8下载中;5-M3U8下载异常;6-M3U8下载结束',
    `crawl_seed_url`    varchar(500) NOT NULL COMMENT '种子URL',
    `crawl_seed_params` varchar(200)          DEFAULT NULL COMMENT '种子URL携带的参数。保存Json串',
    `crawl_m3u8_url`    varchar(600)          DEFAULT NULL COMMENT 'm3u8 url',
    `crawl_m3u8_text`   text COMMENT 'M3U8文本',
    `crawl_m3u8_cnt`    int(11)      NOT NULL DEFAULT '0' COMMENT 'm3u8 url抓取次数',
    `crawl_m3u8_notify` tinyint(4)   NOT NULL DEFAULT '0' COMMENT 'crawl_m3u8_cnt次数超过阈值告警,需要人工介入,大概率要优化代码了',
    `error_msg`         varchar(200)          DEFAULT NULL COMMENT '错误信息',
    `create_user`       varchar(20)           DEFAULT NULL COMMENT '创建者',
    `create_time`       datetime              DEFAULT NULL COMMENT '创建时间',
    `update_user`       varchar(20)           DEFAULT NULL COMMENT '修改者',
    `update_time`       datetime              DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `video_coll_id` (`video_coll_id`, `video_item_id`)
) ENGINE = InnoDB COMMENT ='抓取队列';


CREATE TABLE `crawl_vod_config`
(
    `id`              int(11)      NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `host_type`       int(11)               DEFAULT '2' COMMENT '传递给crawlQueue的hostType字段。2-nivod网；3-BananTV',
    `vod_type`        int(11)               DEFAULT '0' COMMENT '点播类型.0-电影；1-剧集（标志给展示逻辑，爬虫统一按剧集逻辑走）',
    `domain_key_part` varchar(50)  NOT NULL COMMENT '域名关键部分.用于配置策略',
    `program_no`      varchar(20)  NOT NULL DEFAULT '0' COMMENT '栏目编号',
    `program_name`    varchar(20)  NOT NULL COMMENT '栏目名称',
    `program_icon`    varchar(255)          DEFAULT NULL COMMENT '栏目图标',
    `category_no`     varchar(20)           DEFAULT '' COMMENT '分类编号',
    `category_name`   varchar(20)  NOT NULL COMMENT '分类名称',
    `seed_url`        varchar(500) NOT NULL COMMENT '种子URL',
    `seed_params`     varchar(200)          DEFAULT NULL COMMENT '种子URL携带的参数。保存Json串',
    `seed_status`     tinyint(1)   NOT NULL DEFAULT '1' COMMENT '状态：1在用 2停用',
    `page_size`       int(11)               DEFAULT NULL COMMENT '翻页次数',
    `seed_desc`       varchar(255)          DEFAULT NULL COMMENT '描述',
    `error_msg`       varchar(200)          DEFAULT NULL COMMENT '错误信息',
    `create_user`     varchar(20)           DEFAULT NULL COMMENT '创建者',
    `create_time`     datetime              DEFAULT NULL COMMENT '创建时间',
    `update_user`     varchar(20)           DEFAULT NULL COMMENT '修改者',
    `update_time`     datetime              DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 4
  DEFAULT CHARSET = utf8mb4 COMMENT ='爬取点播整站爬取配置';

CREATE TABLE `crawl_vod_config_task`
(
    `id`            int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `vod_config_id` int(11) NOT NULL COMMENT 'vod_config表ID',
    `task_status`   int(11) NOT NULL DEFAULT '0' COMMENT '任务状态. 0-初始化; 1-进行中; 2-任务失败; 3-任务结束',
    `create_user`   varchar(20)      DEFAULT NULL COMMENT '创建者',
    `create_time`   datetime         DEFAULT NULL COMMENT '创建时间',
    `update_user`   varchar(20)      DEFAULT NULL COMMENT '修改者',
    `update_time`   datetime         DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 8
  DEFAULT CHARSET = utf8mb4 COMMENT ='整站爬取配置生成的任务实例表';

CREATE TABLE `crawl_vod_tv`
(
    `id`             int(11)      NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `vod_config_id`  int(11)      NOT NULL COMMENT '配置表主键ID',
    `vod_md5`        varchar(32)           DEFAULT NULL COMMENT '电视剧md5。防重复抓取',
    `crawl_status`   int(11)               DEFAULT '0' COMMENT '抓取状态.0-创建任务;1-自动补全视频信息中;2-补充视频信息失败;3-补充视频信息成功;4-补充TV ID信息中;5-补充TV ID信息失败;6-补充TV ID信息成功',
    `video_country`  varchar(20)  NOT NULL DEFAULT '未知' COMMENT '国家',
    `video_year`     varchar(20)  NOT NULL DEFAULT '未知' COMMENT '年份',
    `video_no`       varchar(20)  NOT NULL DEFAULT '0' COMMENT '栏目编号',
    `video_name`     varchar(255) NOT NULL COMMENT '视频名称',
    `video_director` varchar(50)  NOT NULL DEFAULT '未知' COMMENT '视频导演',
    `video_actor`    varchar(255) NOT NULL DEFAULT '未知' COMMENT '视频演员',
    `video_icon`     varchar(255)          DEFAULT '' COMMENT '栏目图标',
    `video_desc`     text COMMENT '栏目描述',
    `seed_url`       varchar(500) NOT NULL COMMENT '种子URL',
    `seed_params`    varchar(200)          DEFAULT NULL COMMENT '种子URL携带的参数。保存Json串',
    `error_cnt`      int(11)      NOT NULL DEFAULT '0' COMMENT '失败次数',
    `error_msg`      varchar(200)          DEFAULT NULL COMMENT '错误信息',
    `video_language` varchar(20)           DEFAULT '未知' COMMENT '语言',
    `video_quality`  varchar(20)           DEFAULT '未知' COMMENT '清晰度',
    `video_tag`      varchar(255)          DEFAULT NULL,
    `video_coll_id`  bigint(20)            DEFAULT NULL COMMENT '剧集ID',
    `create_user`    varchar(20)           DEFAULT NULL COMMENT '创建者',
    `create_time`    datetime              DEFAULT NULL COMMENT '创建时间',
    `update_user`    varchar(20)           DEFAULT NULL COMMENT '修改者',
    `update_time`    datetime              DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB COMMENT ='爬取点播视频清单列表';

CREATE TABLE `crawl_vod_tv_item`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT,
    `tv_id`        int(4)  NOT NULL COMMENT 'tv表ID',
    `tv_item_md5`  varchar(50)  DEFAULT NULL COMMENT '集数MD5',
    `crawl_status` int(11)      DEFAULT NULL COMMENT '抓取状态.0-INIT;1-自动补全视频信息中;2-补充视频信息失败;3-补充视频信息成功;4-补充TV ID信息中;5-补充TV ID信息失败;6-补充TV ID信息成功',
    `seed_url`     varchar(500) DEFAULT NULL COMMENT '种子URL',
    `seed_params`  varchar(500) DEFAULT NULL COMMENT '种子URL参数',
    `error_msg`    varchar(50)  DEFAULT NULL COMMENT '错误信息',
    `episodes`     varchar(50)  DEFAULT NULL COMMENT '集数',
    `create_user`  varchar(20)  DEFAULT NULL COMMENT '创建者',
    `create_time`  datetime     DEFAULT NULL COMMENT '创建时间',
    `update_user`  varchar(20)  DEFAULT NULL COMMENT '修改者',
    `update_time`  datetime     DEFAULT NULL COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB

