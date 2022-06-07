FROM golang:1.18.3 as builder

WORKDIR /go/src/github.com/fbdlampayan/k8s-device-plugin/

COPY . /go/src/github.com/fbdlampayan/k8s-device-plugin/

RUN go get -v github.com/intel/intel-device-plugins-for-kubernetes/pkg/deviceplugin && \
    go get -v github.com/pkg/errors && \
    go get -v k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o fbdl -v main.go

FROM alpine:3.16.0

WORKDIR /work/

COPY --from=builder /go/src/github.com/fbdlampayan/k8s-device-plugin/fbdl .

ENTRYPOINT ["./fbdl"]
