# begin build container definition
FROM registry.access.redhat.com/ubi8/ubi-minimal as build

# Install golang
RUN microdnf install -y golang

ENV GOBIN=/bin \
    GOPATH=/go

# install pod-logger
RUN /usr/bin/go install github.com/rhdedgar/pod-logger@master


# begin run container definition
FROM registry.access.redhat.com/ubi8/ubi-minimal as run

ADD scripts/ /usr/local/bin/

COPY --from=build /bin/pod-logger /usr/local/bin

EXPOSE 8080

CMD /usr/local/bin/start.sh
