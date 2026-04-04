CREATE TABLE short_link
(
    id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
    short_id VARCHAR(255) NOT NULL ,
    url TEXT NOT NULL
);

CREATE INDEX short_id_idx ON short_link(short_id);