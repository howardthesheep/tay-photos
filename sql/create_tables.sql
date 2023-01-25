start transaction;

-- Users Table --
CREATE TABLE IF NOT EXISTS Users (
    id VARCHAR(36) unique not null default UUID() COMMENT 'UUIDv4',
    name varchar(255) not null,
    username varchar(255) not null,
    email varchar(255) not null,
    password varchar(128) not null,
    apiToken VARCHAR(36) not null unique,

    PRIMARY KEY (id),
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_apiToken (apiToken)
) engine=InnoDB default charset=utf8;

-- Galleries Table --
CREATE TABLE IF NOT EXISTS Galleries (
    id VARCHAR(36) unique not null default UUID() COMMENT 'UUIDv4',
    owner VARCHAR(36) UNIQUE NOT NULL COMMENT 'UUID of Gallery Creator',
    name varchar(255) not null default 'Default Gallery',
    canonicalUrl varchar(255) not null unique COMMENT 'Alternate URL for Gallery',
    createTime Date not null default Date(NOW()),

    PRIMARY KEY (id),
    FOREIGN KEY (owner)
        REFERENCES Users(id)
        ON DELETE CASCADE ON UPDATE CASCADE
) engine=InnoDB default charset=utf8;

-- Gallery Photos Table --
CREATE TABLE IF NOT EXISTS GalleryPhotos (
    id VARCHAR(36) unique not null default UUID() COMMENT 'UUIDv4',
    gallery VARCHAR(36) NOT NULL COMMENT 'UUID of parent Gallery',
    collection VARCHAR(255) COMMENT 'The Gallery Subsection the Photo exists in, if null then the photo is the Header Photo',
    photoPath VARCHAR(255) NOT NULL COMMENT 'File path to the photo',

    FOREIGN KEY (gallery)
        REFERENCES Galleries(id)
        ON UPDATE CASCADE ON DELETE CASCADE
) engine=InnoDB default charset=utf8;


commit;