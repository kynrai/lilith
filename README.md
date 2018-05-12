# Lilith
Lilith the evil distributed monolith done right

# NOTE
This is under very heavy development so many sections are still missing in the code

## Examples

Services serve as examples of how something using a particiular backing technology might work

TODO:
- HTTP
    - [x] Hello world server
    - [x] CORS example
    - [ ] AuthN middleware
- Databases
    - [ ] DynamoDB (AWS)
    - [x] Postgres
    - [x] Datastore (google cloud)
    - [ ] Spanner (google cloud)
    - [ ] TableStore (alicloud)
- Big Data
    - [ ] BigQuery (google cloud)
- Streaming
    - [ ] Kinesis (AWS)
    - [ ] Cloud PubSub (google cloud)
    - [ ] Kafka
    - [ ] NATS
- Auth
    - [ ] Ladon
    - [ ] Firebase Auth
    - [ ] AWS Cognito
- CI/CD
    - [x] CircleCI
    - [ ] TravisCI
- Kubernetes
    - [ ] Helm Tiller Setup
    - [ ] Cert Manager on GCE
    - [ ] Cron Jobs
    - [ ] API Services 

## Motivation
Microservice are cool and great for scaling, especailly with kubernetes. However experiance tells us that its rarely a good idea to start a greenfield project buy building out complex microservice infrastrucutre and architecture. Especailly before your team grows.

While it is widely accepted that Distributed monoliths are a bad thing and I agree most of the time, on day one of a project when you just need to get something working, MVP, Demo to stakeholders and gernal real world stuff outside of the developer bubble world we would like to live in, I propose we use a modified single app.

Golang has high enough performance that doing almost everything in one binary is possible for most MVPs especially if all you need is a CRUD API and lets face it, most of what modern backend dev involes a CRUD API of some sort. Further more this can be very easily scaled out simply by replicating the docker container / binary etc.

We can go one step further by orgaising our Monolith in such a way to decouple all features just enough that when that lovely day comes to make true microservices, the conversion is completly trivial.

Finally real world systems often have a need for schedualed and batch services to load data, maniulate huge data sets etc. These fall outside the remit of RESTful APIs, in Golang these will simply be more docker containers but run on schedual via kubernetes.

## Implementation template

The primary purpose of this repo is not as a library or SDK you can use, but rather as a template of how you can structure you code to contian many, many abstractions to make developing "services" much easier. All examples contain demos of how to test them as well as examples of how to interact with popular databases for CRUD operations and more.

The server itself has all the gotchas handled such as easily defining CORS rules, adding authentication, authorization, CI/CD.

The code makes very heavy use of dependancy injection, arguable the best way to develop and unit test functions in go.

## Building
By convention all golang binaries should live in the `cmd/<name>/`, This forms the buildable binaries which are all build with:

    make build

## Testing
Tests in each golang package will be run, You can run all the tests with:

    make local

A conveniance feature to test a CRUD API iterativly is to run a foreman server after the tests. This way you can easily setup environment varibales your app might need when on a server, and setup run scripts.
Foreman uses `.env` and `Procfile` to run the app. Both are included here for demonstration. Local build shsould rebuild all binaries and rerun all tests. To run the server locally run:

    make local

You can edit the Makefile to your own version of foreman

Integration tests can be done against docker versions of some of the datastores.

|Database|Docker hub|
|---|---|
|Google Datastore|[docker-datastore](https://hub.docker.com/r/kynrai/docker-datastore)|
|DynamoDB|[docker-dynamodb](https://hub.docker.com/r/kynrai/docker-dynamodb/)|

## Deployment
The services are best deployed on kubernetes, there are several high quality managed kubernetes out there. Google Cloud, Microsoft Azure both offer production kubernetes clusters. AWS has a preview of a managed kubernetes cluster.

The advantage of kubernetes comes from protability of both images and platform (no vendor lock), Fast A/B deployments, very good control of traffic, and access to schedulers to make deploying background processes easier.

Kubernets is managed from its CLI tool kubectl and deployment files. Samples provided in the repo in due course.

## Documentation

Expect more documentation later
