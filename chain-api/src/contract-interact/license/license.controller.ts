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
@ApiTags('license')
@Controller('license')
export class LicenseController {
  constructor(private readonly licenseService: LicenseService) {}
  @Post('request')
  async getLicenseRequest(@Body() dto: RequestLicenseDto) {
    return this.licenseService.request(dto);
  }

  @Post('deploy')
  async deploy(@Body() dto: DeployLicenseDto) {
    return this.licenseService.deploy(dto);
  }

  @Get('total-payout')
  async getLicenseResponse(@Body() dto: GetLicenseInfoDto) {
    return this.licenseService.getTotalPayoutInUSD(dto);
  }

  @Get('shares')
  async getShares(@Body() dto: GetShareLicense) {
    return this.licenseService.getShares(dto);
  }

  @Get('owners')
  async getOwners(@Body() dto: GetLicenseInfoDto) {
    return this.licenseService.getOwners(dto);
  }

  @Get('payout-contract')
  async getPayoutContract(@Body() dto: GetLicenseInfoDto) {
    return this.licenseService.getPayoutContract(dto);
  }

  @Post('payout')
  async payout(@Body() dto: LicensePayoutDto) {
    return this.licenseService.payout(dto);
  }

  @Post('set-payout-contract')
  async setPayoutContract(@Body() dto: SetPayoutContractDto) {
    return this.licenseService.setPayoutContract(dto);
  }
}
