start transaction;

-- Users Table --
CREATE TABLE IF NOT EXISTS Users (
    id BINARY(16) unique not null COMMENT 'UUIDv4',
    name varchar(255) not null,
    username varchar(255) not null,
    email varchar(255) not null,
    password varchar(128) not null,
    apiToken BINARY(16) not null unique,

    PRIMARY KEY (id),
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_apiToken (apiToken)
) engine=InnoDB default charset=utf8;

-- Galleries Table --
CREATE TABLE IF NOT EXISTS Galleries (
    id BINARY(16) unique not null COMMENT 'UUIDv4',
    owner BINARY(16) UNIQUE NOT NULL COMMENT 'UUID of Gallery Creator',
    name varchar(255) not null default 'Default Gallery',
    canonicalUrl tinytext not null unique COMMENT 'Alternate URL for Gallery',

    PRIMARY KEY (id),
    FOREIGN KEY (owner)
        REFERENCES Users(id)
        ON DELETE CASCADE ON UPDATE CASCADE
) engine=InnoDB default charset=utf8;

-- Gallery Photos Table --
CREATE TABLE IF NOT EXISTS GalleryPhotos (
    gallery BINARY(16) UNIQUE NOT NULL COMMENT 'UUID of parent Gallery',
    photoPath VARCHAR(255) NOT NULL COMMENT 'File path to the photo',

    FOREIGN KEY (gallery)
        REFERENCES Galleries(id)
        ON UPDATE CASCADE ON DELETE CASCADE
) engine=InnoDB default charset=utf8;


commit;