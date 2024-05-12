import { ethers } from 'ethers';

// Define a TypeScript type for the EventLog based on the provided structure
export type TransactionLogs = {
  provider: ethers.JsonRpcApiProvider;
  transactionHash: string;
  blockHash: string;
  blockNumber: number;
  removed: boolean | undefined;
  address: string;
  data: string;
  topics: string[];
  index: number;
  transactionIndex: number;
  interface: Interface;
  fragment: EventFragment;
};

type Interface = {
  fragments: Fragment[];
  deploy: ConstructorFragment[];
  fallback: any | null;
  receive: boolean;
};

type Fragment = {};

type ConstructorFragment = {};

type EventFragment = {
  type: string;
  inputs: any[];
  name: string;
  anonymous: boolean;
};

type SubmitArgs = {
  args: [
    owner: string,
    txIndex: bigint,
    to: string,
    value: bigint,
    data: string,
  ];
};

type ConfirmArgs = {
  args: [owner: string, txIndex: bigint];
};
type ExecuteArgs = {
  args: [owner: string, txIndex: bigint];
};
type DepositArgs = {
  args: [owner: string, value: bigint, address: string];
};

export type SubmitTransactionLogs = TransactionLogs & SubmitArgs;
export type ConfirmTransactionLogs = TransactionLogs & ConfirmArgs;
export type ExecuteTransactionLogs = TransactionLogs & ExecuteArgs;
export type DepositLogs = TransactionLogs & DepositArgs;
