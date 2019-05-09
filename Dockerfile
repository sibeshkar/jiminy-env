FROM ubuntu:16.04

RUN useradd -ms /bin/bash xvnc

RUN apt-get update
RUN apt-get install -y wget gcc default-jdk firefox htop sudo nano vnc4server \
    && apt-get -y autoremove \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY ./jiminy /app/

COPY ./plugin-go-grpc/wob-v0.zip /app/

RUN mv /app/jiminy /usr/bin/

WORKDIR /app
CMD ["/bin/bash"]