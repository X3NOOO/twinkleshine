# builder
FROM golang:latest as builder
LABEL builder=true

WORKDIR /src

COPY . .

RUN mkdir -p out

RUN go mod download
RUN go build -o out/twinkleshine -v .

# runner
FROM photon:latest
LABEL builder=false

WORKDIR /app

COPY --from=builder /src/out/twinkleshine .    

ENTRYPOINT /app/twinkleshine -config /app/config.yaml -env /app/.env -verbose