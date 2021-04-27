# Webmine
#### Webmine - is simplest web-manager for _Minecraft_-forge-server running on docker-container.
The main task of Webmine is to make uploading mods to a forge minecraft-server very easy, without any tools other than a browser!

#### Works only with [itzg/minecraft-server](https://hub.docker.com/r/itzg/minecraft-server) minecraft-server or similar.

## Installation
### Requirements
- docker
- git
#### For linux
```bash
git clone git@github.com:Kiyamov-Bulat/webmine.git
cd webmine
touch docker-compose.yml
... Copy example of docker file to docker-compose.yml...
sudo docker-compose up
```

#### For windows
_alt + R_ to open cmd
```cmd
git clone git@github.com:Kiyamov-Bulat/webmine.git
cd webmine
type NUL > docker-compose.yml
... Copy example of docker file to docker-compose.yml...
docker-compose up
```
_Example of docker-compose.yml file_
```yaml
version: "3.9"
services:
    minecraft-server:
        image: itzg/minecraft-server:java8
        ports:
          - 25565:25565
        restart: always
        environment:
          EULA: "TRUE"
          VERSION: "1.12.2"
          TYPE: "FORGE"
        deploy:
          resources:
            limits:
              memory: 1.5G
        volumes:
          - /home/Kiyamov-Bulat/data:/data
    webmine:
        build: .
        ports: 
            - "80:80"
        restart: always
        env_file:
            - production.env
        volumes:
            - /var/run/docker.sock:/var/run/docker.sock:rw
            - /home/Kiyamov-Bulat/data:/go/src/webmine/data
```

## ENVIROMENT VARIABLES
ENV_KEY | DEFAULT_VALUE
--------|--------------
PORT | 80
USERS | '{"email":"admin", "password":"admin", "name":"admin"}'
AUTH_TOKEN | "aGVsbG8gd29ybGQgMTIzIQo="
SERVER_IMG_NAME | "itzg/minecraft-server:java8"
CONTAINER_ID |

* You can define **USERS** as json array (with fields email, password and name) or as single object.
* If you define **CONTAINER_ID** webmine will manage the container with this id.
* If you define **SERVER_IMG_NAME** webmine will find container id with this img_name. **WORK ONLY WITH SINGLE CONTAIER!!!**
