# AWS STS Proxy

A simple proxy that can be used to proxy AWS STS based off an oidc token. An oidc token contains all the information about a user that is required to create a temporary sts session and return that session back to the user. The STS session maps back to the user, because the session name is created based on the oidc token's username. Only the server application is allowed to assume the role it is creating tokens for, this allows for us to trust the name of the session name.


### Usage

#### Configuration

```
EKS_ASSUME_ROLE: The Role to assume from the server
STRING_REQUIREMENT: A string to require in the users email address, or a 403 is thrown.
PORT: A port to run the application on. Default is 8080
HEALTHCHECK: A path to serve the healthcheck on. Default is /hc
```

#### Run Locally

```
dep ensure
go run main.go
```

#### POST `/sts/token`

> Returns temproary credentials for a role the server assumes. User must pass Authentication TOKEN Header with request from oidc application. The server creats a session with the email retrieved from the oidc token.

##### Params

```
ROLE_ARN = A role the user wants to assume. The server must be able to assume this role or it will return a 403
Duration = The Duration of the temporary credentials. If the role does not accept this duration, the server will return 403.
ExternalId = An optional ExternalID if the role that is being assumed requires it. If this is not passed in and the role expects it, the server will return 403.
```


> request a temporary sts token that maps back to your user

#### Example

```
curl -XPOST -H"Authorization: $TOKEN" localhost:8080/sts/token
```

### Build

```
docker build -t sts-proxy .
docker-compose up
```
## License
aws-sts-proxy is © copyright by WW International.

aws-sts-proxy is licensed under the [Apache-2.0 Open Source license](http://choosealicense.com/licenses/apache-2.0/).
