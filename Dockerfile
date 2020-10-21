FROM latonaio/l4t:latest

# Definition of a Device & Service
ENV POSITION=Runtime \
    SERVICE=container-image-sweeper \
    AION_HOME=/var/lib/aion

# 1day
ENV INTERVAL_TIME_SECOND 2592000

ADD . .
RUN pip3 install -U pip && pip3 install -r requirement.txt

CMD ["python3", "-u", "main.py"]
