import { Body, Controller, Post, Get, Param } from '@nestjs/common';
import { AgreementService } from './agreement.service';
import {
  DeployAgreementDto,
  GetAgreementInfoDto,
  RequestAgreementDto,
} from './agreement.dto';
import { ApiTags } from '@nestjs/swagger';
import { GetHeader } from '../../decorators/getHeader.decorator';
@ApiTags('Agreement')
@Controller('agreements')
export class AgreementController {
  constructor(private readonly agreementService: AgreementService) {}

  @Post('deploy')
  async deployAgreement(
    @Body() deployDto: DeployAgreementDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return await this.agreementService.deploy(deployDto, seed);
  }

  @Get(':contractAddress')
  async getAgreementResponse(
    @Param('contractAddress') contractAddress: string,
    @GetHeader('X-Seed') seed: string,
  ) {
    return await this.agreementService.getResponse({ contractAddress }, seed);
  }
  @Post('request')
  async requestAgreement(
    @Body() requestDto: RequestAgreementDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return await this.agreementService.request(requestDto, seed);
  }
}
