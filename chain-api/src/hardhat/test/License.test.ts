import { StreamingRightsManagement } from '../../../typechain';
import { CHAINLINK } from '../../config/chainlink.config';

const { expect } = require('chai');
const { ethers } = require('hardhat');

describe('StreamingRightsManagement', function () {
  let streamingRightsManagement: StreamingRightsManagement,
    payContract,
    owner,
    addr1,
    addr2;
  const shares = [25, 25, 50];

  beforeEach(async function () {
    [owner, addr1, addr2] = await ethers.getSigners();

    const Payroll = await ethers.getContractFactory('Payroll');
    payContract = await Payroll.deploy(owner.address, owner.address); // assume an oracle price feed address

    const StreamingRightsManagement = await ethers.getContractFactory(
      'StreamingRightsManagement',
    );
    streamingRightsManagement = await StreamingRightsManagement.deploy(
      CHAINLINK.AMOY.CHAINLINK_TOKEN, // Chainlink Token address
      CHAINLINK.AMOY.ORACLE_ADDRESS, // Oracle address
      CHAINLINK.AMOY.JOB_IDS.UINT,
      0,
      owner.address,
      [owner.address, addr1.address, addr2.address],
      shares,
    );
  });

  describe('Initialization', function () {
    it('should set owners and shares correctly', async function () {
      expect(await streamingRightsManagement.getShare(owner.address)).to.equal(
        25,
      );
      expect(await streamingRightsManagement.getShare(addr1.address)).to.equal(
        25,
      );
      expect(await streamingRightsManagement.getShare(addr2.address)).to.equal(
        50,
      );
    });
  });

  describe('Payout Functionality', function () {
    it('should successfully call payout', async function () {
      await streamingRightsManagement.setPayoutContract(payContract.address);
      await expect(streamingRightsManagement.payout()).to.not.be.reverted;
    });
  });

  // More tests as needed for other functions
});
