import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";
//npx hardhat ignition deploy ignition/modules/Salaries.ts --network amoy
//SalariesModule#Salaries - 0xac45e95Dd5C7F9B1a6C3e4883d04952B9C974b05
const SalariesModule = buildModule("SalariesModule", (m) => {
  const salaryContract = m.contract("Salaries");

  const answer = m.call(salaryContract, "getChainlinkDataFeedLatestAnswer", []);
  console.log("🚀 ~ SalariesModule ~ answer:", answer);
  return { salaryContract };
});
export default SalariesModule;
