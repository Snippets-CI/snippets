CREATE TABLE users (
        user_id UUID NOT NULL DEFAULT gen_random_uuid(),
        username VARCHAR(25) NOT NULL,
        password VARCHAR(60) NOT NULL,
        CONSTRAINT "primary" PRIMARY KEY (user_id ASC),
        FAMILY "primary" (user_id, username, password)
);

CREATE TABLE snippets (
        snippet_id UUID NOT NULL DEFAULT gen_random_uuid(),
        owner UUID NOT NULL,
        language VARCHAR(20) NULL,
        title VARCHAR(30) NULL,
        category VARCHAR(50) NULL,
        code VARCHAR NULL,
        CONSTRAINT "primary" PRIMARY KEY (snippet_id ASC),
        INDEX snippets_auto_index_fk_owner_ref_users (owner ASC),
        FAMILY "primary" (snippet_id, owner, language, title, category, code)
);

INSERT INTO users (user_id, username, password) VALUES
        ('99ce9c3c-6210-4ed0-a759-5d289f4d2cd0', 'markus', 'clearPassword');

ALTER TABLE snippets ADD CONSTRAINT fk_owner_ref_users FOREIGN KEY (owner) REFERENCES users(user_id);

-- Validate foreign key constraints. These can fail if there was unvalidated data during the dump.
ALTER TABLE snippets VALIDATE CONSTRAINT fk_owner_ref_users;