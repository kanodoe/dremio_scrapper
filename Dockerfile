FROM golang:alpine AS build

RUN apk add --update make

WORKDIR /app

COPY . .

RUN make build

FROM alpine
WORKDIR /app
COPY --from=build /app/build/bin/dremio_scrapper /app


ENTRYPOINT [ "/app/dremio_scrapper" ]