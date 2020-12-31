# Docker steps

## Prerequisites
https://docs.docker.com/get-docker/  
https://docs.docker.com/compose/install/

Enable file sharing for docker desktop first.

# Postgres & pgAdmin

| Label      | Inputfield   |
|------------|--------------|
| Host name/address: | `snippets_postgres_db` |
| port:  | `5432`           |
| username:  | `admin`      |
| password:  |   `123`      |

pgAdmin can be reached from [localhost](localhost:5050).  

# Quick docker startup command:  

> docker run -e POSTGRES_PASSWORD=123 -e POSTGRES_USER=admin -e POSTGRES_DB=postgres -d --name snippets_postgres_db -p 5432:5432 andreasroither/snippets_db:latest

