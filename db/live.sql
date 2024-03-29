use live;

CREATE TABLE IF NOT EXISTS users
(
    uid bigint unsigned NOT NULL AUTO_INCREMENT,
    name    varchar(36) NOT NULL,
    passwd  varchar(32) NOT NULL,
    salt    varchar(32) NOT NULL,
    phone   varchar(16) NOT NULL DEFAULT '',
    headurl     varchar(256) NOT NULL DEFAULT '',
    nickname    varchar(64) NOT NULL DEFAULT '', 
    -- role 0:普通用户 1:VIP用户
    role        tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
    stime   datetime NOT NULL DEFAULT '2017-12-01',
    etime   datetime NOT NULL DEFAULT '2017-12-01',
    PRIMARY KEY(uid),
    UNIQUE KEY(name)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS user_info
(
    id      bigint unsigned NOT NULL AUTO_INCREMENT,
    uid     int unsigned NOT NULL,
    income  int unsigned NOT NULL DEFAULT 0,
    apply   int unsigned NOT NULL DEFAULT 0,
    withdraw    int unsigned NOT NULL DEFAULT 0,
    recharge    int unsigned NOT NULL DEFAULT 0,
    expense     int unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY(id),
    UNIQUE KEY(uid)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS withdraw_history
(
    id      bigint unsigned NOT NULL AUTO_INCREMENT,
    uid     int unsigned NOT NULL,
    amount  int unsigned NOT NULL DEFAULT 0,
    -- status 0:申请  1:成功 2:拒绝
    status  tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
    PRIMARY KEY(id),
    KEY(uid)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS stream
(
    id  bigint unsigned NOT NULL AUTO_INCREMENT,
    name    varchar(36) NOT NULL DEFAULT '',
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
    uid     int unsigned NOT NULL DEFAULT 0,
    -- status 0:创建 1-推流 2-停止 3-暂停
    status  tinyint unsigned NOT NULL DEFAULT 0,
    cur_id  int unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY(id),
    UNIQUE KEY(name),
    KEY(uid)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS live_history
(
    id      bigint unsigned NOT NULL AUTO_INCREMENT,
    uid     int unsigned NOT NULL,
    title   varchar(128) NOT NULL DEFAULT '',
    cover   varchar(128) NOT NULL DEFAULT '',
    depict  varchar(256) NOT NULL DEFAULT '',
    -- authority 访问权限 0:公开 1:密码 2:付费
    authority   tinyint unsigned NOT NULL DEFAULT 0,
    passwd  varchar(16) NOT NULL DEFAULT '',
    price   int unsigned NOT NULL DEFAULT 0,
    resolution tinyint unsigned NOT NULL DEFAULT 0,
    push    varchar(256) NOT NULL DEFAULT '',
    replay  varchar(256) NOT NULL DEFAULT '',
    -- status 0:创建 1-推流 2-停止 3-暂停
    status  tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
    ftime   datetime NOT NULL DEFAULT '2017-12-01',
    PRIMARY KEY(id),
    KEY(uid)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS orders
(
    id  bigint unsigned NOT NULL AUTO_INCREMENT,
    oid varchar(64) NOT NULL,
    uid int unsigned NOT NULL,
    -- type 0:支付视频  1-充值
    type    tinyint unsigned NOT NULL DEFAULT 0,
    owner int unsigned NOT NULL DEFAULT 0,
    hid   int unsigned NOT NULL,
    price int unsigned NOT NULL DEFAULT 0,
    fee   int unsigned NOT NULL DEFAULT 0,
    status  tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
    ptime   datetime NOT NULL DEFAULT '2017-12-01',
    PRIMARY KEY(id),
    UNIQUE KEY(oid),
    KEY(uid),
    KEY(owner)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS pay_items
(
    id      int unsigned NOT NULL AUTO_INCREMENT,
    price   int unsigned NOT NULL DEFAULT 0,
    product varchar(64) NOT NULL DEFAULT '',
    img     varchar(128) NOT NULL DEFAULT '',
    online  tinyint unsigned NOT NULL DEFAULT 0,
    deleted tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
    PRIMARY KEY(id)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS image
(
    id     bigint unsigned NOT NULL AUTO_INCREMENT,
    name    varchar(48) NOT NULL,
    uid     int unsigned NOT NULL,
    filesize    int unsigned NOT NULL,
    height      int unsigned NOT NULL,
    width       int unsigned NOT NULL,
    status      tinyint unsigned NOT NULL default 0,
    deleted     tinyint unsigned NOT NULL default 0,
    ctime       datetime NOT NULL default '0000-00-00 00:00:00',
    ftime       datetime NOT NULL default '0000-00-00 00:00:00',
    PRIMARY KEY(id),
    UNIQUE KEY(name)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS channel
(
    id      int unsigned NOT NULL AUTO_INCREMENT,
    uid     int unsigned NOT NULL,
    -- 频道名称
    title   varchar(128) NOT NULL DEFAULT '',
    -- 频道轮播图
    cover1  varchar(128) NOT NULL DEFAULT '',
    cover2  varchar(128) NOT NULL DEFAULT '',
    cover3  varchar(128) NOT NULL DEFAULT '',
    -- 公众号二维码
    qrcode  varchar(128) NOT NULL DEFAULT '',
    -- 直播简介
    depict  varchar(256) NOT NULL DEFAULT '',
    -- 频道介绍
    chan_intro  varchar(128) NOT NULL DEFAULT '',
    -- 直播介绍
    live_intro  varchar(128) NOT NULL DEFAULT '',
    -- 企业公众号
    wxmp        varchar(128) NOT NULL DEFAULT '',
    -- 是否显示频道简介
    display     tinyint unsigned NOT NULL DEFAULT 0, 
    -- 频道地址
    dst         varchar(128) NOT NULL DEFAULT '',
    -- 附加功能
    extra       tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
    PRIMARY KEY(id),
    UNIQUE KEY(uid)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS tags
(
    id      bigint unsigned NOT NULL AUTO_INCREMENT,
    uid     int unsigned NOT NULL,
    content varchar(128) NOT NULL DEFAULT '',
    priority    int unsigned NOT NULL DEFAULT 0,
    online  tinyint unsigned NOT NULL DEFAULT 0,
    deleted tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
    PRIMARY KEY(id),
    KEY(uid)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS feedback
(
    id      bigint unsigned NOT NULL AUTO_INCREMENT,
    uid     int unsigned NOT NULL,
    title   varchar(128) NOT NULL DEFAULT '',
    content varchar(512) NOT NULL DEFAULT '',
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
    img     varchar(128) NOT NULL DEFAULT '',
    status  tinyint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY(id),
    KEY(uid)
) ENGINE = InnoDB;
