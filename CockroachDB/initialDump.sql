CREATE TABLE users (
        userid INT8 NOT NULL,
        username VARCHAR(25) NOT NULL,
        password VARCHAR(60) NOT NULL,
        CONSTRAINT "primary" PRIMARY KEY (userid ASC),
        FAMILY "primary" (userid, username, password)
);

CREATE TABLE snippets (
        userid INT8 NOT NULL,
        owner INT8 NOT NULL,
        language VARCHAR(20) NULL,
        title VARCHAR(30) NULL,
        category VARCHAR(50) NULL,
        code VARCHAR NULL,
        CONSTRAINT "primary" PRIMARY KEY (userid ASC),
        INDEX snippets_auto_index_fk_owner_ref_users (owner ASC),
        FAMILY "primary" (userid, owner, language, title, category, code)
);

ALTER TABLE snippets ADD CONSTRAINT fk_owner_ref_users FOREIGN KEY (owner) REFERENCES users(userid);

-- Validate foreign key constraints. These can fail if there was unvalidated data during the dump.
ALTER TABLE snippets VALIDATE CONSTRAINT fk_owner_ref_users;