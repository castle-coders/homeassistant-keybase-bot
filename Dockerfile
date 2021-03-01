FROM keybaseio/client:stable-slim
ENV KEYBASE_SERVICE=1
COPY bin/castlebot /castlebot
CMD /castlebot