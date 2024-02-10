FROM node:20

WORKDIR /

COPY package*.json /
COPY *.mjs /

RUN npm install

CMD ["node", "index.mjs"]
