FROM registry.access.redhat.com/ubi8-minimal

EXPOSE 8080 5000

RUN microdnf update -y && rm -rf /var/cache/yum && microdnf install git go make -y && microdnf clean all

COPY . /opt/grpc-demo
WORKDIR /opt/grpc-demo

RUN make build

CMD ["/opt/grpc-demo/bin/grpc-demo-product"]
