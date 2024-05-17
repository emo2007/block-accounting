import { PriceFeedMock, Salaries } from '../../../typechain';

const { ethers } = require('hardhat');
const { expect } = require('chai');
import { SignerWithAddress } from '@nomicfoundation/hardhat-ethers/signers';

describe('Salaries', function () {
  let salaries: Salaries;
  let owner: SignerWithAddress;
  let multisigWallet: SignerWithAddress;
  let addr1: SignerWithAddress;
  let priceFeedMock: PriceFeedMock;

  beforeEach(async function () {
    [owner, multisigWallet, addr1] = await ethers.getSigners();

    const PriceFeedMockFactory =
      await ethers.getContractFactory('PriceFeedMock');
    priceFeedMock = await PriceFeedMockFactory.deploy();
    await priceFeedMock.getDeployedCode();
    // Deploy the Salaries contract
    const SalariesFactory = await ethers.getContractFactory('Salaries');
    salaries = (await SalariesFactory.deploy(
      multisigWallet.address,
      await priceFeedMock.getAddress(),
    )) as Salaries;
    await salaries.getDeployedCode();
  });

  it('Should set and get salary correctly', async function () {
    await salaries.connect(multisigWallet).setSalary(addr1.address, 1000);
    expect(await salaries.getUsdtSalary(addr1.address)).to.equal(1000);
  });

  it('Should payout in ETH correctly', async function () {
    // Set the salary in USDT
    await salaries.connect(multisigWallet).setSalary(addr1.address, 100);
    expect(await salaries.getUsdtSalary(addr1.address)).to.equal(100);

    // Fund the contract with ETH
    await owner.sendTransaction({
      to: await salaries.getAddress(),
      value: ethers.parseEther('1'), // 1 ETH
    });

    await expect(() =>
      salaries.connect(multisigWallet).payoutInETH(addr1.address),
    ).to.changeEtherBalances(
      [salaries, addr1],
      ['-32393909944930353', '32393909944930353'],
    );

    // Check events
    expect(salaries.connect(multisigWallet).payoutInETH(addr1.address));
  });
});

describe('PriceFeedMock', function () {
  it('Should return the mocked price', async function () {
    const PriceFeedMockFactory =
      await ethers.getContractFactory('PriceFeedMock');
    const priceFeedMock = await PriceFeedMockFactory.deploy();
    await priceFeedMock.getDeployedCode();

    expect((await priceFeedMock.latestRoundData())[1].toString()).to.equal(
      '3087',
    );
  });
});
