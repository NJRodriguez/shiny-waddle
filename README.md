# shiny-waddle

Shiny-waddle is a demo project written in Golang using Pulumi for infrastructure deployment and Docker for containerizing our web service.

## Tools Required
- [Golang](https://golang.org/dl/)
- [Pulumi](https://www.pulumi.com/docs/get-started/install/)
- [AWS CLI](https://aws.amazon.com/cli/)
- [Docker](https://www.docker.com/products/docker-desktop)

## Installation

This project uses [go modules](https://blog.golang.org/using-go-modules) for dependency management so be sure to have the latest version of Golang installed.

For infrastructure deployment you'll need to install [Pulumi](https://www.pulumi.com/docs/get-started/install/) and set up an account. You can login in your environment using pulumi login command.

```
pulumi login
```

Pulumi utilizes the [AWS Command Line interface](https://aws.amazon.com/cli/) for secure credentials retrieval. Once installed, you must set up your account.

```
aws configure
```

This will prompt you for the following:

- AWS Access Key ID
- AWS Secret Access Key
- Region name
- Output format (default)

Once both Pulumi and AWS CLI are properly configured, you can deploy your stack and build your docker image.

### Infrastructure deployment

```
pulumi up
```

Pulumi will preview the stack that will be deployed to your AWS account and after accepting, will begin creating the necessary infrastructure for this project.

### Container deploy
You can build the docker image by moving to `project/api` directory and running the following commands:

1. Build the GO executable
```BASH
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
```

2. Building the docker image
```BASH
docker build -f DockerFile -t test .
```

3. Running the docker image
```BASH
docker run --publish 80:80 --env TABLE_NAME=sucursal_table --env AWS_REGION=us-east-1 --env AWS_ACCESS_KEY_ID=<YOUR ACCESS KEY HERE> --env AWS_SECRET_ACCESS_KEY=<YOUR SECRET ACCESS KEY HERE> test
```

Your docker container is up and running! You can see logs from the command line or using Docker Desktop.

If you want to avoid running the docker image, you can also launch this locally using VSCode and the following launch configuration:

```JSON
{
    "name": "Debug Server",
    "type": "go",
    "request": "launch",
    "mode": "debug",
    "remotePath": "",
    "port": 80,
    "host": "0.0.0.0",
    "program": "${workspaceRoot}/api/main.go",
    "env": {
        "TABLE_NAME": "sucursal_table",
        "AWS_REGION": "us-east-1"
    },
    "args": [],
    "showLog": true
}
```

If you set up your AWS CLI correctly, credentials will be loaded automagically when your debug session is created.


## Usage

Web service is deployed to `0.0.0.0:80`. The following are the endpoints available:

### /sucursal POST
Will create a new Sucursal in the database.

```
+-----------+---------+-------------------------------------+--------------------------------------+
| Property  |  Type   |             Description             |               Example                |
+-----------+---------+-------------------------------------+--------------------------------------+
| ID        | UUID    | Unique identifier for this Sucursal | b309060a-ce7b-4649-abc1-4cf3f6e51d1b |
| Address   | String  | Physical address of the Sucursal    | 123 Fake St.                         |
| Latitude  | Float64 | Precise latitude of Sucursal        | -34.604258                           |
| Longitude | Float64 | Precise longitude of Sucursal       | -58.375094                           |
+-----------+---------+-------------------------------------+--------------------------------------+
```

#### Example request

```JSON
{
    "id":"b309060a-ce7b-4649-abc1-4cf3f6e51d1b",
    "address": "Florida 296, C1005 CABA",
    "latitude": -34.604258,
    "longitude": -58.375094
}
```

#### Example response

```JSON
{
    "id": "b309060a-ce7b-4649-abc1-4cf3f6e51d1b",
    "message":"Sucursal has been created successfully!"
}
```

### /sucursal/{id} GET
Will retrieve the sucursal ID from the database.

```
+----------+-------+--------------------------------------+
| Property | Type  |               Example                |
+----------+-------+--------------------------------------+
| ID       | UUID  | b309060a-ce7b-4649-abc1-4cf3f6e51d1b |
+----------+-------+--------------------------------------+
```

#### Example request
```HTTP
http://0.0.0.0:80/sucursal/b309060a-ce7b-4649-abc1-4cf3f6e51d1b
```

#### Example response
```JSON
{
    "id": "b309060a-ce7b-4649-abc1-4cf3f6e51d1b",
    "address": "Florida 296, C1005 CABA",
    "latitude": -34.604258,
    "longitude": -58.375094
}
```
### /sucursal/{lat}/{lon} GET
Will retrieve the closest sucursal based on the latitude and longitude path arguments.

```
+-----------+---------+-------------+----------+
| Property  |  Type   | Description | Example  |
+-----------+---------+-------------+----------+
| Latitude  | float64 | -90 ~ 90    |  58.2314 |
| Longitude | float64 | -120 ~ 120  | 102.3644 |
+-----------+---------+-------------+----------+
```

#### Example request
```HTTP
http://0.0.0.0:80/sucursal/-34.613217/-58.374625
```

#### Example response
```JSON
{
    "Sucursal": {
        "id": "b309060a-ce7b-4649-abc1-4cf3f6e51d1b",
        "address": "Florida 296, C1005 CABA",
        "latitude": -34.604258,
        "longitude": -58.375094
    },
    "DistanceInKm": 0.9970716278723797
}
```