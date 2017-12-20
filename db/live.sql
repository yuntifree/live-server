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
    -- role 0:普通用户 
    role        tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
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
    PRIMARY KEY(id),
    UNIQUE KEY(name),
    KEY(uid)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS orders
(
    id  bigint unsigned NOT NULL AUTO_INCREMENT,
    oid varchar(64) NOT NULL,
    uid int unsigned NOT NULL,
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
