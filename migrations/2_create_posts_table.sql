CREATE TABLE posts(
    id SERIAL NOT NULL,
    PRIMARY KEY (id),
    user_id SERIAL NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    body TEXT NOT NULL,
    upvote_count INTEGER DEFAULT 0,
    downvote_count INTEGER DEFAULT 0);

CREATE TABLE  user_post_votes(
    id SERIAL NOT NULL,
    PRIMARY KEY (id),
    user_id SERIAL NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    post_id SERIAL NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id),
    vote_type INTEGER NOT NULL);