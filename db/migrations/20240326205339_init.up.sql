CREATE SEQUENCE IF NOT EXISTS user_id_seq;
CREATE SEQUENCE IF NOT EXISTS product_id_seq;

CREATE TABLE IF NOT EXISTS public."user"
(
    id          BIGINT                   DEFAULT NEXTVAL('user_id_seq'::regclass) NOT NULL PRIMARY KEY,
    login       TEXT UNIQUE                                                       NOT NULL CHECK (login <> '')
    CONSTRAINT  max_len_name CHECK (LENGTH(login) <= 256),
    password    TEXT                                                              NOT NULL CHECK (password <> '')
    CONSTRAINT  max_len_password CHECK (LENGTH(password) <= 256),
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()                            NOT NULL
);

CREATE TABLE IF NOT EXISTS public."product"
(
    id              BIGINT                   DEFAULT NEXTVAL('product_id_seq'::regclass) NOT NULL PRIMARY KEY,
    saler_id        BIGINT                                                               NOT NULL REFERENCES public."user" (id),
    title           TEXT                                                                 NOT NULL CHECK (title <> '')
    CONSTRAINT max_len_title CHECK (LENGTH(title) <= 256),
    description     TEXT                                                                 NOT NULL CHECK (description <> '')
    CONSTRAINT max_len_description CHECK (LENGTH(description) <= 4000),
    price           BIGINT                   DEFAULT 0                                   NOT NULL
    CONSTRAINT not_negative_price CHECK (price >= 0),
    image_url       TEXT                                                                 DEFAULT ''
    CONSTRAINT max_len_image_url CHECK (LENGTH(image_url) <= 256),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT NOW()                               NOT NULL
);