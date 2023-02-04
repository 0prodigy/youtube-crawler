FROM golang:1.19.1 AS build
WORKDIR /app

RUN --mount=type=secret,id=GITHUB_TOKEN,required=true git config --global url."https://$(cat /run/secrets/GITHUB_TOKEN)@github.com/".insteadOf "https://github.com/"
RUN go env -w GOPRIVATE="github.com/NimbleBoxAI/*"
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build .


FROM gcr.io/distroless/base AS oneclick-server
ARG BINARY
WORKDIR /

COPY --from=build /app/youtube-scraooer /crwaler

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["/crwaler"]
