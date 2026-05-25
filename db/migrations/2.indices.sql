-- Posts
CREATE INDEX posts_by_user_and_date ON posts (userid, created);
CREATE INDEX posts_by_vote ON posts (votes);
CREATE INDEX posts_by_date ON posts (created);

-- Comments
CREATE INDEX comments_by_date ON comments (created);
CREATE INDEX comments_by_post_and_date ON comments (postid, created);

-- Replies
CREATE INDEX replies_by_comment_and_date ON replies (parentid, created);

-- Votes
-- We already have an index on votes by user through the composite primary key,
-- so we only need an additional index for votes by post.
CREATE INDEX votes_by_post ON votes (postid);

-- Users
CREATE INDEX users_by_username ON users (username);
