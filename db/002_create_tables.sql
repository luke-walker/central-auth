DROP TABLE IF EXISTS servers CASCADE;
CREATE TABLE servers (
    id              uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    name            text UNIQUE NOT NULL,
    addresses       text[],
    proxy_url       text,
    redirect_url    text NOT NULL,
    token           uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    PRIMARY KEY (id)
);

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
    id              uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    token           uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    username        text UNIQUE NOT NULL,
    password        text NOT NULL,
    last_ip         text,
    admin           boolean NOT NULL DEFAULT FALSE,
    PRIMARY KEY (id, token)
);

DROP TABLE IF EXISTS sessions CASCADE;
CREATE TABLE sessions (
    user_ip         text NOT NULL,
    user_token      uuid NOT NULL,
    access_token    uuid UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    expires         date,
    PRIMARY KEY (user_ip, user_token),
    FOREIGN KEY (user_token) REFERENCES users(token) ON DELETE CASCADE
);
