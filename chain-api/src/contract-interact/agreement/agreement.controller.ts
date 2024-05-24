import { Body, Controller, Post, Get, Param } from '@nestjs/common';
import { AgreementService } from './agreement.service';
import {
  DeployAgreementDto,
  GetAgreementInfoDto,
  RequestAgreementDto,
} from './agreement.dto';
import { ApiTags } from '@nestjs/swagger';
@ApiTags('Agreement')
@Controller('agreements')
export class AgreementController {
  constructor(private readonly agreementService: AgreementService) {}

  @Post('deploy')
  async deployAgreement(@Body() deployDto: DeployAgreementDto) {
    return await this.agreementService.deploy(deployDto);
  }

  @Get(':contractAddress')
  async getAgreementResponse(
    @Param('contractAddress') contractAddress: string,
  ) {
    return await this.agreementService.getResponse({ contractAddress });
  }
  @Post('request')
  async requestAgreement(@Body() requestDto: RequestAgreementDto) {
    return await this.agreementService.request(requestDto);
  }
}
