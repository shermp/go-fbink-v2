FROM golang:1.22

RUN apt update && apt install -y unzip 
RUN curl -sSL -o /tmp/kobo-tc.zip https://github.com/koreader/koxtoolchain/releases/download/2021.12/kobo.zip && \
    cd /tmp && \
    unzip /tmp/kobo-tc.zip && \
    tar -xzvf /tmp/kobo.tar.gz -C /opt/ && \
    rm -rf /tmp/kobo*

ENV PATH=/opt/x-tools/arm-kobo-linux-gnueabihf/bin:$PATH
ENV CROSS_COMPILE=arm-kobo-linux-gnueabihf