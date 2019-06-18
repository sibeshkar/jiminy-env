FROM ubuntu:16.04

RUN useradd -ms /bin/bash xvnc

RUN apt-get update && apt-get install -y wget gcc default-jdk firefox htop sudo vnc4server \
    && apt-get -y autoremove \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY ./jiminy /usr/bin/

COPY ./plugin-go-grpc/wob-v0.zip /app/

#RUN jiminy install /app/wob-v0.zip
#CMD ["jiminy", "run"]

WORKDIR /app
CMD ["/bin/bash"]

