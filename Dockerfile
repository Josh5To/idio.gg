FROM golang:1.21

ENV CGO_ENABLED=0
# Move to our project folder and run the program
ADD . /web-ttfn
WORKDIR /web-ttfn
RUN go build -o tt-fn-website

EXPOSE 8180

ENTRYPOINT [ "./tt-fn-website" ]


