start transaction;
    -- Insert Users Test Data
    INSERT INTO Users (id,name,username,email,password,apiToken)
    VALUES (
        'ef45499a-8cc6-44f8-a33f-3df7275b20ee',
        'Taylor Lindquist',
        'tlindquist',
        'tlindquistt@gmail.com',
        '$2a$10$22fvBj4nx1f90eVVw2lpZuIlGJQnCCPRgqtrI7d6.3c5ZkK6htv3C', -- asdf, bcrypt 10 rounds --
        'bb0f0933-f619-4533-8da1-910988e92e31'
    );

    -- Insert Galleries Test Data --
    INSERT INTO Galleries (id, owner, name, canonicalUrl)
    VALUES (
        '1f579937-d99c-4e0f-98a7-6a7e0976bcab',
        'ef45499a-8cc6-44f8-a33f-3df7275b20ee',
        'Fun Time Gallery :)',
        'funTimeGallery'
    );

    -- Insert GalleryPhotos Test Data --
    INSERT INTO GalleryPhotos (gallery, collection, photoPath)
    VALUES ('1f579937-d99c-4e0f-98a7-6a7e0976bcab', 'The most fun collection', '/photos/kitten.jpg');
    INSERT INTO GalleryPhotos (gallery, collection, photoPath)
    VALUES ('1f579937-d99c-4e0f-98a7-6a7e0976bcab', 'The most fun collection', '/photos/tayrick.png');
    INSERT INTO GalleryPhotos (gallery, collection, photoPath)
    VALUES ('1f579937-d99c-4e0f-98a7-6a7e0976bcab', 'The 2nd most fun collection', '/photos/tayrick.png');
    INSERT INTO GalleryPhotos (gallery, collection, photoPath)
    VALUES ('1f579937-d99c-4e0f-98a7-6a7e0976bcab', 'The 2nd most fun collection', '/photos/tayrick.png');

commit;