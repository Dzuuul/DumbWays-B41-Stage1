CREATE TABLE tb_projects (
    id              SERIAL PRIMARY KEY,
    user_id         integer NOT NULL REFERENCES tb_users(id),
    name            VARCHAR(128) NOT NULL,
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    description     VARCHAR(1028) NOT NULL DEFAULT '',
    technologies    VARCHAR(1028) NOT NULL DEFAULT '',
    image           VARCHAR(1028) NOT NULL DEFAULT ''
)