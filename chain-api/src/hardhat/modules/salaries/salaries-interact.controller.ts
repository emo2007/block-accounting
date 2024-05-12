import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { SalariesService } from './salaries.service';
import {
  GetEmployeeSalariesDto,
  SalariesDeployDto,
  SetSalaryDto,
} from './salaries.dto';
import { ApiTags } from '@nestjs/swagger';
@ApiTags('salaries')
@Controller('salaries')
export class SalariesController {
  constructor(private readonly salariesService: SalariesService) {}

  @Post('deploy')
  async deploy(@Body() dto: SalariesDeployDto) {
    return this.salariesService.deploy(dto);
  }

  @Get('usdt-price/:contractAddress')
  async getUsdtPrice(@Param('contractAddress') contractAddress: string) {
    return this.salariesService.getLatestUSDTPrice(contractAddress);
  }

  @Post('set-salary')
  async setSalary(@Body() dto: SetSalaryDto) {
    return this.salariesService.setSalary(dto);
  }

  @Get('salary')
  async getSalary(@Body() dto: GetEmployeeSalariesDto) {
    return this.salariesService.getSalary(dto);
  }
}
