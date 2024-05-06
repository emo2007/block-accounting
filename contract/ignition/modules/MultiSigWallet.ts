import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const JAN_1ST_2030 = 1893456000;
const ONE_GWEI: bigint = 1_000_000_000n;

const owners = ["0xfE87F7EF2a58a1f363a444332df6c131C683e35f"];

const MultiSigModule = buildModule("MultiSigWallet", (m) => {
  const ownerP = m.getParameter("owners", owners);
  const confirmationsP = m.getParameter("_numConfirmationsRequired", 1);

  const deploy = m.contract("MultiSigWallet", [ownerP, confirmationsP]);

  return { deploy };
});

export default MultiSigModule;
