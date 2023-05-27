FROM golang:1.16-alpine

LABEL description="Docker Template for Go"
ARG NAME=go-docker
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
