FROM node:22.5.1

WORKDIR /

COPY package*.json /
COPY *.mjs /

RUN npm install

CMD ["node", "index.mjs"]
