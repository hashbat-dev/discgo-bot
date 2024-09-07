FROM golang:1.20

RUN apt-get update && apt-get install -y git

# Set the Current Working Directory inside the container
WORKDIR /go/src/femboy-control

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go get github.com/PuerkitoBio/goquery
RUN go get github.com/bwmarrin/discordgo
RUN go get github.com/clinet/discordgo-embed
RUN go get github.com/Yakiyo/nekos_best.go
RUN go get github.com/go-sql-driver/mysql
RUN go get -u golang.org/x/image/webp
RUN go get -u github.com/chai2010/webp
RUN go get github.com/shirou/gopsutil/cpu
RUN go get github.com/shirou/gopsutil/mem
RUN go get github.com/google/uuid
RUN go get github.com/nfnt/resize
COPY . .

RUN set CGO_CFLAGS=-IC:\libwebp\include
RUN set CGO_LDFLAGS=-LC:\libwebp\lib -lwebp

# Build the Go app
RUN go build -o ./out/femboy-control .

# This container exposes ports 3333 to the outside world
EXPOSE 3333

# Run the binary program produced by `go install`
CMD ["./out/femboy-control"]
