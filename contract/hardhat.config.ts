import { HardhatUserConfig, vars } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import dotenv from "dotenv";
dotenv.config();
const config: HardhatUserConfig = {
  solidity: "0.8.24",
  networks: {
    polygon: {
      url: `https://polygon-amoy.g.alchemy.com/v2/pEtFFy_Qr_NrM1vMnlzSXmYXkozVNzLy`,
      accounts: [process.env.POLYGON_PK || ""],
    },
  },
};

export default config;
