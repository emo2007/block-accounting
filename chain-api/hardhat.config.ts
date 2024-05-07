require('@nomicfoundation/hardhat-toolbox');
require('@nomicfoundation/hardhat-ethers');
const dotenv = require('dotenv');
dotenv.config();

const config = {
  solidity: '0.8.24',
  networks: {
    amoy: {
      url: `https://polygon-amoy.g.alchemy.com/v2/pEtFFy_Qr_NrM1vMnlzSXmYXkozVNzLy`,
      accounts: [process.env.POLYGON_PK || ''],
    },
  },
  paths: {
    sources: './src/hardhat/contracts',
    // tests: './src/hardhat/test',
    cache: './src/hardhat/cache',
    artifacts: './src/hardhat/artifacts',
  },
};

module.exports = config;
