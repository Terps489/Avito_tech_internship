CREATE TABLE teams (
    team_name TEXT PRIMARY KEY
);

CREATE TABLE users (
    user_id   TEXT PRIMARY KEY,
    username  TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    team_name TEXT NOT NULL REFERENCES teams(team_name)
);

CREATE TABLE pull_requests (
    pull_request_id   TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id         TEXT NOT NULL REFERENCES users(user_id),
    status            TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at         TIMESTAMPTZ
);

CREATE TABLE pull_request_reviewers (
    pr_id       TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id TEXT NOT NULL REFERENCES users(user_id),
    PRIMARY KEY (pr_id, reviewer_id)
);
