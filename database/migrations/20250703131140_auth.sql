-- migrate:up
CREATE TABLE users
(
    id         VARCHAR(21) PRIMARY KEY DEFAULT nanoid(),
    email      VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255)        NOT NULL,
    last_name  VARCHAR(255)        NOT NULL,
    phone      VARCHAR(255)        NOT NULL,
    created_at TIMESTAMPTZ             DEFAULT now(),
    updated_at TIMESTAMPTZ             DEFAULT now()
);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_phone ON users (phone);

CREATE TABLE passwords
(
    id         VARCHAR(21) PRIMARY KEY DEFAULT nanoid(),
    user_id    VARCHAR(21)  NOT NULL,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ             DEFAULT now(),

    UNIQUE (user_id, password),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE sessions
(
    id         VARCHAR(21) PRIMARY KEY DEFAULT nanoid() NOT NULL,
    user_id    VARCHAR(21) REFERENCES users,
    token      VARCHAR(255)                             NOT NULL UNIQUE,
    created_at TIMESTAMPTZ             DEFAULT now(),
    updated_at TIMESTAMPTZ             DEFAULT now(),
    expires_at TIMESTAMPTZ                              NOT NULL
);

CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_token ON sessions (token);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);

CREATE TABLE auth_codes
(
    id         VARCHAR(21) PRIMARY KEY DEFAULT nanoid() NOT NULL,
    email      VARCHAR(255)                             NOT NULL,
    code       VARCHAR(255)                             NOT NULL UNIQUE,
    created_at TIMESTAMPTZ             DEFAULT now(),
    expires_at TIMESTAMPTZ                              NOT NULL,
    FOREIGN KEY (email) REFERENCES users (email) ON DELETE CASCADE
);

CREATE INDEX idx_login_codes_email ON auth_codes (email);
CREATE INDEX idx_login_codes_code ON auth_codes (code);
CREATE INDEX idx_login_codes_expires_at ON auth_codes (expires_at);

-- migrate:down
DROP TABLE users;
DROP TABLE user_passwords;
DROP TABLE sessions;
DROP TABLE auth_codes;
