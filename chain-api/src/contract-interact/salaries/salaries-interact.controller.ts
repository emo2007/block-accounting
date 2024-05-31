import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { SalariesService } from './salaries.service';
import {
  CreatePayoutDto,
  DeployContractResponseDto,
  GetEmployeeSalariesDto,
  SalariesDeployDto,
  SetSalaryDto,
} from './salaries.dto';
import { ApiOkResponse, ApiTags } from '@nestjs/swagger';
import { DepositContractDto } from '../multi-sig.dto';
import { GetHeader } from '../../decorators/getHeader.decorator';

@ApiTags('salaries')
@Controller('salaries')
export class SalariesController {
  constructor(private readonly salariesService: SalariesService) {}

  @ApiOkResponse({
    type: DeployContractResponseDto,
  })
  @Post('deploy')
  async deploy(
    @Body() dto: SalariesDeployDto,
    @GetHeader('X-Seed') seed: string,
  ): Promise<DeployContractResponseDto> {
    const address = await this.salariesService.deploy(dto, seed);
    return {
      address,
    };
  }

  @Get('usdt-price/:contractAddress')
  async getUsdtPrice(
    @Param('contractAddress') contractAddress: string,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.salariesService.getLatestUSDTPrice(contractAddress, seed);
  }

  @Post('set-salary')
  async setSalary(
    @Body() dto: SetSalaryDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.salariesService.setSalary(dto, seed);
  }

  @Get('salary')
  async getSalary(
    @Body() dto: GetEmployeeSalariesDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.salariesService.getSalary(dto, seed);
  }

  @Post('payout')
  async createPayout(
    @Body() dto: CreatePayoutDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.salariesService.createPayout(dto, seed);
  }

  @Post('deposit')
  async deposit(
    @Body() dto: DepositContractDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.salariesService.deposit(dto, seed);
  }
}
