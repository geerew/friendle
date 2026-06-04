-- +goose Up

CREATE TABLE users (
    id            TEXT PRIMARY KEY NOT NULL,
    username      TEXT UNIQUE NOT NULL COLLATE NOCASE,
    display_name  TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    site_role     TEXT NOT NULL CHECK(site_role IN ('site_admin', 'site_user')),
    created_at    TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    updated_at    TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW'))
);

CREATE TABLE sessions (
    id      TEXT PRIMARY KEY NOT NULL,
    data    BLOB NOT NULL,
    expires BIGINT NOT NULL,
    user_id TEXT NOT NULL DEFAULT ''
);

CREATE TABLE groups (
    id              TEXT PRIMARY KEY NOT NULL,
    name            TEXT UNIQUE NOT NULL COLLATE NOCASE,
    created_by      TEXT NOT NULL,
    created_at      TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    updated_at      TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE group_members (
    id            TEXT PRIMARY KEY NOT NULL,
    group_id      TEXT NOT NULL,
    user_id       TEXT NOT NULL,
    group_role    TEXT NOT NULL CHECK(group_role IN ('group_admin', 'group_user')),
    times_picked  INTEGER NOT NULL DEFAULT 0,
    picker_skips  INTEGER NOT NULL DEFAULT 0,
    created_at    TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    updated_at    TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE(group_id, user_id)
);

CREATE TABLE group_join_requests (
    id         TEXT PRIMARY KEY NOT NULL,
    group_id   TEXT NOT NULL,
    user_id    TEXT NOT NULL,
    status     TEXT NOT NULL CHECK(status IN ('pending', 'approved', 'rejected')) DEFAULT 'pending',
    created_at TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    updated_at TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE(group_id, user_id)
);

CREATE TABLE rounds (
    id              TEXT PRIMARY KEY NOT NULL,
    group_id        TEXT NOT NULL,
    round_date      TEXT NOT NULL,
    picker_user_id  TEXT NOT NULL,
    word_plain      TEXT,
    status          TEXT NOT NULL CHECK(status IN ('awaiting_word', 'active', 'completed', 'skipped')),
    created_at      TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    updated_at      TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE,
    FOREIGN KEY (picker_user_id) REFERENCES users (id),
    UNIQUE(group_id, round_date)
);

CREATE TABLE round_participations (
    id              TEXT PRIMARY KEY NOT NULL,
    round_id        TEXT NOT NULL,
    user_id         TEXT NOT NULL,
    solved          BOOLEAN NOT NULL DEFAULT FALSE,
    finished        BOOLEAN NOT NULL DEFAULT FALSE,
    score           INTEGER NOT NULL DEFAULT 0,
    first_guess_at  TEXT,
    completed_at    TEXT,
    created_at      TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    updated_at      TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    FOREIGN KEY (round_id) REFERENCES rounds (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE(round_id, user_id)
);

CREATE TABLE guesses (
    id              TEXT PRIMARY KEY NOT NULL,
    round_id        TEXT NOT NULL,
    user_id         TEXT NOT NULL,
    attempt         INTEGER NOT NULL CHECK(attempt BETWEEN 1 AND 6),
    word            TEXT NOT NULL,
    result          TEXT NOT NULL,
    outcome         TEXT NOT NULL CHECK(outcome IN ('correct', 'partial', 'incorrect')),
    created_at      TEXT NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
    FOREIGN KEY (round_id) REFERENCES rounds (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE(round_id, user_id, attempt)
);

CREATE INDEX idx_sessions_expires ON sessions(expires);
CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_group_members_group ON group_members(group_id);
CREATE INDEX idx_group_members_user ON group_members(user_id);
CREATE INDEX idx_join_requests_group_status ON group_join_requests(group_id, status);
CREATE INDEX idx_rounds_group_date ON rounds(group_id, round_date);
CREATE INDEX idx_rounds_group_status ON rounds(group_id, status);
CREATE INDEX idx_participations_round ON round_participations(round_id);
CREATE INDEX idx_guesses_round_user ON guesses(round_id, user_id);

-- +goose Down

DROP TABLE IF EXISTS guesses;
DROP TABLE IF EXISTS round_participations;
DROP TABLE IF EXISTS rounds;
DROP TABLE IF EXISTS group_join_requests;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
