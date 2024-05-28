FROM node:latest

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

RUN npm run build

RUN npx hardhat compile && npx hardhat typechain

CMD [ "node", "dist/src/main.js" ]
