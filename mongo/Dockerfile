FROM mongo:latest

MAINTAINER <vlad>

ADD mongod.conf /etc/mongod.conf
ADD create_users.js /tmp/
RUN echo "Creating mongo users..."
RUN mongod --fork --logpath /var/log/mongodb.log \
    && echo "Waiting..." \
    && sleep 5 \
    && echo "Start..." \
    && mongo admin /tmp/create_users.js