import { TransactionReceipt, ethers } from 'ethers';

export const parseLogs = (
  txReceipt: TransactionReceipt,
  contract: ethers.Contract,
  eventName: string,
) => {
  return txReceipt.logs
    .map((log) => contract.interface.parseLog(log))
    .find((log) => !!log && log.fragment.name === eventName);
};
