import { TransactionReceipt, ethers } from 'ethers';

export const parseLogs = (
  txReceipt: TransactionReceipt,
  contract: ethers.Contract,
) => {
  return txReceipt.logs
    .map((log) => contract.interface.parseLog(log))
    .find((log) => !!log);
};
