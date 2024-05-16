# blockd backend
## Build
### Locally
1. Install Go >= 1.22
``` sh
curl -LO https://get.golang.org/$(uname)/go_installer && \
        chmod +x go_installer && \
        ./go_installer && \
        rm go_installer
```
2. Install docker:
``` sh
output=$(which docker);
if [ -z "${output}" ]; then 
    sudo dnf remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-selinux \
                  docker-engine-selinux \
                  docker-engine
    sudo apt -y install dnf-plugins-core
    sudo dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
    sudo dnf install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    sudo systemctl start docker
fi
```
3. Build it!:
``` sh
make bin.build
```

4. Start the server:
``` sh
make d.net && \
sudo docker compose up blockd-db -d && \
make run.debug
```
Or
``` sh
make d.net && \
sudo docker compose up blockd-db -d && \
make run.local
```

### Docker
Just run
``` sh
make up
```

# API 
Request content type: application/json  
Response content type: application/json  

## POST **/join**  
### Request body:  
name (string, optional)  
credentals (object, optional)  
        credentals.email (string, optional)   
        credentals.phone (string, optional)   
        credentals.telegram (string, optional)   
mnemonic (string, **required**)   

### Example
Request: 
``` bash
curl --location 'http://localhost:8081/join' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Bladee The Grand Drainer",
    "credentals": {
        "email": "bladeee@gmail.com",
        "phone": "+79999999999",
        "telegram": "@thebladee"
    },
    "mnemonic":"airport donate language disagree dumb access insect tribe ozone humor foot jealous much digital confirm"
}'
```

Response: 
``` json 
{
    "token": "token",
    "token_expired_at": 1715975501581,
    "refresh_token": "refresh_token",
    "refresh_token_expired_at": 1716407501581
}
```

## POST **/login**  
### Request body:  
mnemonic (string, **required**)   

### Example
Request: 
``` bash
curl --location 'http://localhost:8081/login' \
--header 'Content-Type: application/json' \
--data '{
    "mnemonic":"airport donate language disagree dumb access insect tribe ozone humor foot jealous much digital confirm"
}'
```

Response: 
``` json 
{
    "token": "token",
    "token_expired_at": 1715975501581,
    "refresh_token": "refresh_token",
    "refresh_token_expired_at": 1716407501581
}
```

## POST **/refresh**  
### Request body:  
token (string, **required**)   
refresh_token (string, **required**)   

### Example
Request: 
``` bash
curl --location --request GET 'http://localhost:8081/refresh' \
--header 'Content-Type: application/json' \
--data '{
        "token": "token",
        "refresh_token": "refresh_token"
}'
```

Response: 
``` json 
{
    "token": "token",
    "token_expired_at": 1715975501581,
    "refresh_token": "refresh_token",
    "refresh_token_expired_at": 1716407501581
}
```

## POST **/organizations**  
### Request body:  
name (string, **required**)  
address (string, optional)
wallet_mnemonic (string, optional. *if not provided, creators mnemonic will me used*)

### Example
Request: 
``` bash
curl --location 'http://localhost:8081/organizations' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU0NTY4Mzg4NTAsInVpZCI6ImI2NmU1Mjk4LTU1ZTctNGIxNy1hYzliLTA0MzU3YjBlN2Q0ZSJ9.K1I0QoZEdDYK_HEsJ0PdWOfZ8ugTcPfLqy7fHhvK9nk' \
--data '{
    "name": "The Drain Gang Inc",
    "address": "Backsippestigen 22, 432 36 Varberg, Sweden"
}'
```

Response: 
``` json 
{
    "id": "dfac7846-0f0a-11ef-9262-0242ac120002"
}
```

## GET **/organizations**  
### Request body:  
cursor (string, optional)  
limit (uint8, optional. Max:50, Default:50)
offset_date (uint63, optional. *time as unix milli*)

### Example
Request: 
``` bash
curl --location --request GET 'http://localhost:8081/organizations' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTU4OTU4MTc3MDYsInVpZCI6IjUyNTNkMzdjLTMxZDQtNDgxMi1iZTcxLWE5ODQwMTVlNGVlMyJ9.YjjHWz7FiMM73e-98pZYHCW9tKDZ_mRWKG3m1PcVTo0' \
--data '{
    "limit":5,
    "cursor":"eyJpZCI6IjAxOGY2ZTc3LWUxNDMtNzcyZi04NjJkLTlkZDM5NzUxYTZkMyJ9"
}'
```

Response: 
``` json 
{
    "_type": "organizations",
    "_links": {
        "self": {
            "href": "/organizations"
        }
    },
    "items": [
        {
            "_links": {
                "self": {
                    "href": "/organizations/018f6e77-ebcc-7547-bc84-2556fbf12300"
                }
            },
            "id": "018f6e77-ebcc-7547-bc84-2556fbf12300",
            "name": "The Drain Gang Inc 6",
            "address": "1",
            "created_at": 1715556104012,
            "updated_at": 1715556104012
        },
        {
            "_links": {
                "self": {
                    "href": "/organizations/018f6e77-f5f5-7bcb-b98f-9966e7a8b706"
                }
            },
            "id": "018f6e77-f5f5-7bcb-b98f-9966e7a8b706",
            "name": "The Drain Gang Inc 7",
            "address": "1",
            "created_at": 1715556106613,
            "updated_at": 1715556106613
        }
    ],
    "pagination": {
        "next_cursor": "eyJpZCI6IjAxOGY2ZTc3LWY1ZjUtN2JjYi1iOThmLTk5NjZlN2E4YjcwNiJ9",
        "total_items": 2
    }
}
```