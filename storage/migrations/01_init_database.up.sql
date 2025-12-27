CREATE TABLE IF NOT EXISTS users(
    user_id serial PRIMARY KEY,
    login VARCHAR (50) UNIQUE NOT NULL,
    password VARCHAR (64) NOT NULL
);

CREATE TABLE IF NOT EXISTS metadata(
    md_id serial PRIMARY KEY,
    user_id INTEGER NOT NULL,
    data_type INTEGER,
    descript VARCHAR (50),
    md_hash VARCHAR(64) UNIQUE NOT NULL,
    upload_data VARCHAR (50) NOT NULL
);

CREATE TABLE IF NOT EXISTS users_data(
    md_id INTEGER PRIMARY KEY,
    user_data BYTEA
);

-- CREATE TABLE IF NOT EXISTS balances(
--     user_id INTEGER PRIMARY KEY,
--     accrual REAL NOT NULL,
--     withdraw REAL NOT NULL
-- );
