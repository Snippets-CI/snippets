# Docker steps

## Prerequisites
https://docs.docker.com/get-docker/  
https://docs.docker.com/compose/install/

Enable file sharing for docker desktop first.

# Postgres & Adminer
Build and start in background:  
docker-compose up -d

Then adminer can be reached from [localhost](localhost:8000).  
Credentials:  
| Label      | Inputfield   |
|------------|--------------|
| System:    | `PostgreSQL` |
| Server:    | `db`         |
| Username:  | `admin`      |
| Password:  | `123`        |
| Database:  |              |


# Postgres & pgAdmin

| Label      | Inputfield   |
|------------|--------------|
| Host name/address: | `db` |
| port:  | `5432`           |
| username:  | `admin`      |
| password:  |   `123`      |

Adminer can be reached from [localhost](localhost:5050).  

# Quick docker startup command:  

> docker run --name postgres-snippets -e POSTGRES_PASSWORD=123 -d -p 5432:5432 postgres

