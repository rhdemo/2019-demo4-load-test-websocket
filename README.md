## RH Summit 2019 - Load Tester on Game (E2E)

This repo contains load test for testing the websocket.

It requires ansible to run it

## Setup

Create a `.env` file with all the required environment variables. An example env file [.env.example](.env.example) is included.

```bash
cp .env.example .env
```

Test that you can login to the server with

```bash
make oc_login
```

## Deploy
After create an .env file, you can create the project with the following command on openshift:

```bash
make deploy-load-test
```

## Removing 
You can remove with the follwing command: 

```bash
make remove-load-test
```