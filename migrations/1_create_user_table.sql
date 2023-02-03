-- Create Users Table
CREATE TABLE users(
    id SERIAL NOT NULL,
    PRIMARY KEY (id),
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    score INTEGER NOT NULL DEFAULT 0);