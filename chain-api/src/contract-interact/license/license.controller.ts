import { Body, Controller, Get, Post } from '@nestjs/common';
import { LicenseService } from './license.service';
import { ApiTags } from '@nestjs/swagger';
import {
  DeployLicenseDto,
  GetLicenseInfoDto,
  GetShareLicense,
  LicensePayoutDto,
  RequestLicenseDto,
  SetPayoutContractDto,
} from './license.dto';
import { GetHeader } from '../../decorators/getHeader.decorator';
@ApiTags('license')
@Controller('license')
export class LicenseController {
  constructor(private readonly licenseService: LicenseService) {}
  @Post('request')
  async getLicenseRequest(
    @Body() dto: RequestLicenseDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.licenseService.request(dto, seed);
  }

  @Post('deploy')
  async deploy(
    @Body() dto: DeployLicenseDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.licenseService.deploy(dto, seed);
  }

  @Get('total-payout')
  async getLicenseResponse(
    @Body() dto: GetLicenseInfoDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.licenseService.getTotalPayoutInUSD(dto, seed);
  }

  @Get('shares')
  async getShares(
    @Body() dto: GetShareLicense,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.licenseService.getShares(dto, seed);
  }

  @Get('owners')
  async getOwners(
    @Body() dto: GetLicenseInfoDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.licenseService.getOwners(dto, seed);
  }

  @Get('payout-contract')
  async getPayoutContract(
    @Body() dto: GetLicenseInfoDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.licenseService.getPayoutContract(dto, seed);
  }

  @Post('payout')
  async payout(
    @Body() dto: LicensePayoutDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.licenseService.payout(dto, seed);
  }

  @Post('set-payout-contract')
  async setPayoutContract(
    @Body() dto: SetPayoutContractDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.licenseService.setPayoutContract(dto, seed);
  }
}
