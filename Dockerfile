FROM golang:1.16-alpine

LABEL description="Mail sending microservice"
ARG NAME=mail-serve
ENV NAME_DIR=${NAME}

# setup
RUN mkdir /${NAME}
COPY . /${NAME}

# build the project
WORKDIR /${NAME}
RUN go mod download && go mod verify
RUN go build -o application .

# run the app
CMD /${NAME_DIR}/application
