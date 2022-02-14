# syntax=docker/dockerfile:1

FROM golang:1.17.6-alpine
WORKDIR /build

COPY . .

# for later moving
COPY ./posts/test.md ./

RUN go get .

RUN go build -o /blogo

WORKDIR /bloGO
RUN blogo init .

ENTRYPOINT [ "/build/docker_entrypoint.sh" ]

EXPOSE 8000
VOLUME [ "/bloGO" ]