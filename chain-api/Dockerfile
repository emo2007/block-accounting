FROM node:latest

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .
RUN npx hardhat compile && npx hardhat typechain

RUN npm run build


CMD [ "node", "dist/src/main.js" ]
