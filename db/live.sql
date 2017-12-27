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
    -- status 0:创建 1-推流
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
    owner int unsigned NOT NULL,
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
