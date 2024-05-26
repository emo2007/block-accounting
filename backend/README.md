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
* name (string, optional)  
* credentals (object, optional)  
        credentals.email (string, optional)   
        credentals.phone (string, optional)   
        credentals.telegram (string, optional)   
* mnemonic (string, **required**)   

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
* mnemonic (string, **required**)   

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
* token (string, **required**)   
* refresh_token (string, **required**)   

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
* name (string, **required**)  
* address (string, optional)
* wallet_mnemonic (string, optional. *if not provided, creators mnemonic will me used*)

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
* cursor (string, optional)  
* limit (uint8, optional. Max:50, Default:50)
* offset_date (uint63, optional. *time as unix milli*)

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

## GET **/{organization_id}/participants**
### Request body:
```json
{
  "limit":1
}
```

### Example
Request:
```bash
curl --request GET \
  --url http://localhost:8081/organizations/018faff4-481f-73ec-a4b8-27ef07b4029b/participants \
  --header 'Authorization: Bearer TOKEN' \
  --header 'content-type: application/json' \
  --data '{
  "limit":1
}'
```
Response: 
```json
{
  "_type": "participants",
  "_links": {
    "self": {
      "href": "/organizations/018faff4-481f-73ec-a4b8-27ef07b4029b/participants"
    }
  },
  "participants": [
    {
      "_type": "participant",
      "_links": {
        "self": {
          "href": "/organizations/018faff4-481f-73ec-a4b8-27ef07b4029b/participants018faff4-25fb-7973-860a-59eb69b766a4"
        }
      },
      "id": "018faff4-25fb-7973-860a-59eb69b766a4",
      "name": "Bladee The Grand Drainer",
      "credentials": {
        "email": "bladeee@gmail.com",
        "phone": "+79999999999",
        "telegram": "@thebladee"
      },
      "created_at": 1716654773151,
      "updated_at": 1716654773151,
      "is_user": true,
      "is_admin": true,
      "is_owner": true,
      "is_active": true
    }
  ]
}
```

## GET **/{organization_id}/transactions**  
### Request body:  


### Example
Request: 
``` bash
curl --location --request GET 'http://localhost:8081/organizations/018f9078-af60-7589-af64-9312b97aa7be/transactions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer TOKEN' \
--data '{
    "description":"New test tx!",
    "amount": 100,  
    "to":"0x323b5d4c32345ced77393b3530b1eed0f346429d",
    "limit":1
}'
```

Response: 
``` json 
{
    "_type": "transaction",
    "_links": {
        "self": {
            "href": "/organizations/{organization_id}/transactions"
        }
    },
    "id": "018f8ce2-dada-75fb-9745-8560e5736bec",
    "description": "New test tx!",
    "organization_id": "018f8ccd-2431-7d21-a0c2-a2735c852764",
    "created_by": "018f8ccc-e4fc-7a46-9628-15f9c3301f5b",
    "amount": 100,
    "to": "MjtdTDI0XO13OTs1MLHu0PNGQp0=",
    "max_fee_allowed": 5,
    "deadline": 123456767,
    "created_at": 1716055628507,
    "updated_at": 1716055628507
}
```

## POST **/{organization_id}/transactions**  
### Request body:  
* ids ([]uuid)  
* created_by (uuid)  
* to (string)  
* cancelled (bool)  
* confirmed (bool)  
* commited (bool)  
* expired (bool)  
* pending (bool)   
* cursor (string)  
* limit (int)  
* offset_date (unix milli time)  

### Example
Request: 
``` bash
curl --location --request GET 'http://localhost:8081/organizations/018f9112-1805-7b5e-ae30-7fc2151810f3/transactions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer TOKEN' \
--data '{
    "to": "0xD53990543641Ee27E2FC670ad2cf3cA65ccDc8BD",
    "pending":true,
    "limit":1,
    "created_by":"018f9111-f0fb-708a-aec1-55295f5496d6"
}'
```

Response: 
``` json 
{
    "_type": "organizations",
    "_links": {
        "self": {
            "href": "/organizations/018f9112-1805-7b5e-ae30-7fc2151810f3/transactions"
        }
    },
    "next_cursor": "eyJpZCI6IjAxOGY5MTE1LWU5NmItN2IxMi04Y2JiLWQxNTY5NDNkYjk5NCJ9",
    "transactions": [
        {
            "_type": "transaction",
            "_links": {
                "self": {
                    "href": "/organizations/018f9112-1805-7b5e-ae30-7fc2151810f3/transactions/018f9115-e96b-7b12-8cbb-d156943db994"
                }
            },
            "id": "018f9115-e96b-7b12-8cbb-d156943db994",
            "description": "Test filter by TO!!!!!",
            "organization_id": "018f9112-1805-7b5e-ae30-7fc2151810f3",
            "created_by": "018f9111-f0fb-708a-aec1-55295f5496d6",
            "amount": 1234,
            "to": "0xD53990543641Ee27E2FC670ad2cf3cA65ccDc8BD",
            "max_fee_allowed": 2.5,
            "created_at": 1716136883437,
            "updated_at": 1716136883437
        }
    ]
}
```
## POST **/organizations/{organization_id}/participants**  
### Request body:  
* name (string)
* position (string)
* wallet_address (string)

### Example
Request: 
``` bash
curl --request POST \
  --url http://localhost:8081/organizations/018fb419-c3ad-7cda-81b8-cad30211b5fb/participants \
  --header 'Authorization: Bearer token' \
  --header 'content-type: application/json' \
  --data '{
  "name":"dodik",
  "position":"employee", 
  "wallet_address":"0x8b1bc2590A3C9A1FEb349f1BacAfbc92CBC50156"
}'
```

Response: 
``` json 
{
  "_type": "participant",
  "_links": {
    "self": {
      "href": "/organizations/018fb419-c3ad-7cda-81b8-cad30211b5fb/participants/018fb42c-81dc-77f8-9eac-9d0540b34441"
    }
  },
  "id": "018fb42c-81dc-77f8-9eac-9d0540b34441",
  "name": "dodik2",
  "created_at": 1716714766812,
  "updated_at": 1716714766812,
  "is_user": false,
  "is_admin": false,
  "is_owner": false,
  "is_active": false
}
```
## GET **/organizations/{organization_id}/participants**  
### Request body:  
* ids (string array)

### Example
Request: 
``` bash
curl --request GET \
  --url http://localhost:8081/organizations/018fb419-c3ad-7cda-81b8-cad30211b5fb/participants \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTY3OTk5MjgzMzIsInVpZCI6IjAxOGZiNDE5LTliZjctN2QwOS05MzViLTNiOTAyNDc3ZDJkYiJ9.V_d3b8MvuOp01xDGX0g5Ab2nOvdyGL84WO01xPodTro' \
  --header 'content-type: application/json' \
  --data '{
	"ids":[
      "018fb419-9bf7-7d09-935b-3b902477d2db",
      "018fb42c-1a60-7dd3-841a-fccce8575091",
      "018fb42c-20b7-7c06-acb6-03b3ccb8b7e5"
    ]
}'
```

Response: 
``` json 
{
  "_type": "participants",
  "_links": {
    "self": {
      "href": "/organizations/018fb419-c3ad-7cda-81b8-cad30211b5fb/participants"
    }
  },
  "participants": [
    {
      "_type": "participant",
      "_links": {
        "self": {
          "href": "/organizations/018fb419-c3ad-7cda-81b8-cad30211b5fb/participants/018fb419-9bf7-7d09-935b-3b902477d2db"
        }
      },
      "id": "018fb419-9bf7-7d09-935b-3b902477d2db",
      "name": "Bladee The Grand Drainer",
      "credentials": {
        "email": "bladeee@gmail.com",
        "phone": "+79999999999",
        "telegram": "@thebladee"
      },
      "created_at": 1716724338478,
      "updated_at": 1716724338478,
      "is_user": true,
      "is_admin": true,
      "is_owner": true,
      "is_active": true
    },
    {
      "_type": "participant",
      "_links": {
        "self": {
          "href": "/organizations/018fb419-c3ad-7cda-81b8-cad30211b5fb/participants/018fb42c-1a60-7dd3-841a-fccce8575091"
        }
      },
      "id": "018fb42c-1a60-7dd3-841a-fccce8575091",
      "name": "New Employee",
      "created_at": 1716725540320,
      "updated_at": 1716725540320,
      "is_user": false,
      "is_admin": false,
      "is_owner": false,
      "is_active": false
    },
    {
      "_type": "participant",
      "_links": {
        "self": {
          "href": "/organizations/018fb419-c3ad-7cda-81b8-cad30211b5fb/participants/018fb42c-20b7-7c06-acb6-03b3ccb8b7e5"
        }
      },
      "id": "018fb42c-20b7-7c06-acb6-03b3ccb8b7e5",
      "name": "New Employee",
      "created_at": 1716725541943,
      "updated_at": 1716725541943,
      "is_user": false,
      "is_admin": false,
      "is_owner": false,
      "is_active": false
    }
  ]
}
```
## POST **/organizations/{organization_id}/participants**  
### Request body:  
* title (string)
* owners (array of object { "public_key":"string" })
* confirmations (uint) 

### Example
Request: 
``` bash
curl --request POST \
  --url http://localhost:8081/organizations/018fb246-1616-7f1b-9fe2-1a3202224695/multisig \
  --header 'Authorization: Bearer token' \
  --header 'content-type: application/json' \
  --data '{
  "title":"new sig",
  "owners":[
    "0x5810f45ac87c0be03b4d8174132e2bc81ba1a928"
  ],
  "confirmations":1
}'
```

Response: 
``` json 
{
  "Ok": true
}
```

