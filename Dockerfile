FROM ubuntu:16.04

RUN useradd -ms /bin/bash xvnc

RUN apt-get update
RUN apt-get install -y wget git gcc default-jdk firefox htop sudo nano vnc4server unzip

COPY ./jiminy /app/

COPY ./wob-v0.zip /app/

WORKDIR /app
CMD ["/bin/bash"]