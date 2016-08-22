# spital-server

To run the server you need a *docker-compose.yml* file :

```yaml
version: '2'
services:
  db:
    image: mysql:5.7
    volumes:
      - "~/.hospital_data/db:/var/lib/mysql"
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: xxx
      MYSQL_DATABASE: xxx
      MYSQL_USER: xxx
      MYSQL_PASSWORD: xx
  spital-server:
    depends_on:
      - db
    image: fabienfoerster/spital-server:latest
    links:
      - db
    ports:
      - "5000:5000"
    restart: always
    environment:
      MYSQL_HOST: db:3306
      MYSQL_USER: xxx
      MYSQL_PASSWORD: xxx
```
