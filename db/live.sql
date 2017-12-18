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
