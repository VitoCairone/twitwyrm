FROM golang

ADD . /go/src/github.com/vitocairone/twitwyrm
RUN cd /go/src/github.com/vitocairone/twitwyrm && go install

CMD ["/go/bin/twitwyrm"]