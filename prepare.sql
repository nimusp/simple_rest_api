CREATE TABLE books (
    id SERIES,
    title VARCHAR(50) NOT NULL,
    author_name VARCHAR(50),
    author_last_name VARCHAR(50),
    PRIMARY KEY (id)
);


INSERT INTO books (title, author_name, author_last_name)
VALUES ('Notre-Dame de Paris', 'Victor', 'Hugo');

INSERT INTO books (title, author_name, author_last_name) 
VALUES ('The Three Musketeers', 'Alexandre', 'Dumas');

INSERT INTO books (title, author_name, author_last_name) 
VALUES ('War and Peace', 'Leo', 'Tolstoy');

INSERT INTO books (title, author_name, author_last_name) 
VALUES ('Crime and Punishment', 'Fyodor', 'Dostoevsky');