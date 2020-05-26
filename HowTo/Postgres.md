# Postgres steps

#### Preparation
- Setup databases since we are not using the standard db (postgres) | since we have separate databases for test and production
  - Login with `admin` and `123` at `localhost:8000` for postgres
  - Create db Snippets
  - Create db SnippetsTest
- Create extension for the gen_random_uuid() function with an sql command
  - > CREATE EXTENSION pgcrypto;