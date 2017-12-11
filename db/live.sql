use live;

CREATE TABLE IF NOT EXISTS stream
(
    id  bigint unsigned NOT NULL AUTO_INCREMENT,
    name    varchar(36) NOT NULL DEFAULT '',
    ctime   datetime NOT NULL DEFAULT '2017-12-01',
    uid     int unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY(id),
    UNIQUE KEY(name),
    KEY(uid)
) ENGINE = InnoDB;
