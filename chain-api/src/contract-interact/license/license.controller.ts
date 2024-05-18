import { Body, Controller, Get, Post } from '@nestjs/common';
import { LicenseService } from './license.service';
import { ApiTags } from '@nestjs/swagger';
import {
  DeployLicenseDto,
  GetShareLicense,
  RequestLicenseDto,
} from './license.dto';
@ApiTags('license')
@Controller('license')
export class LicenseController {
  constructor(private readonly licenseService: LicenseService) {}
  @Get('request')
  async getLicenseRequest(@Body() dto: RequestLicenseDto) {
    return this.licenseService.request(dto);
  }

  @Post('deploy')
  async deploy(@Body() dto: DeployLicenseDto) {
    return this.licenseService.deploy(dto);
  }

  @Get('total-payout')
  async getLicenseResponse(@Body() dto: RequestLicenseDto) {
    return this.licenseService.getTotalPayoutInUSD(dto);
  }

  @Get('shares')
  async getShares(@Body() dto: GetShareLicense) {
    return this.licenseService.getShares(dto);
  }

  @Get('owners')
  async getOwners(@Body() dto: GetShareLicense) {
    return this.licenseService.getOwners(dto);
  }

  @Get('payout-contract')
  async getPayoutContract(@Body() dto: GetShareLicense) {
    return this.licenseService.getPayoutContract(dto);
  }
}
