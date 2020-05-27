CREATE TABLE IF NOT EXISTS users (
	user_id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
	username VARCHAR(25) NOT NULL,
        mail VARCHAR(50) NOT NULL,
	password VARCHAR(60) NOT NULL,
        UNIQUE(mail)
);

CREATE TABLE IF NOT EXISTS snippets (
        snippet_id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
        owner UUID REFERENCES users(user_id) NOT NULL,
        language VARCHAR(20) NULL,
        title VARCHAR(30) NULL,
        category VARCHAR(50) NULL,
        code VARCHAR NULL
);

INSERT INTO users (user_id, username, password) VALUES
        ('99ce9c3c-6210-4ed0-a759-5d289f4d2cd0', 'markus', 'clearPassword');

