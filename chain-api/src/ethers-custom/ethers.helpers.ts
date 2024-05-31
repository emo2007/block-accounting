import { TransactionReceipt, ethers } from 'ethers';

export const parseLogs = (
  txReceipt: TransactionReceipt,
  contract: ethers.Contract,
  eventName: string,
) => {
  const parsedLogs = txReceipt.logs.map((log) =>
    contract.interface.parseLog(log),
  );
  console.log('=>(ethers.helpers.ts:10) parsedLogs', parsedLogs);
  return parsedLogs.find((log) => !!log && log.fragment.name === eventName);
};
