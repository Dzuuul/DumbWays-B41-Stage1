CREATE TABLE tb_projects (
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER NOT NULL REFERENCES tb_users(id),
    name            VARCHAR(128) NOT NULL,
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    description     TEXT NOT NULL,
    technologies    VARCHAR[] NOT NULL,
    image           VARCHAR NOT NULL 
)