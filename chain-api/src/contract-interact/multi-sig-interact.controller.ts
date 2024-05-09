import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { ApiOkResponse, ApiTags } from '@nestjs/swagger';
import { MultiSigWalletService } from 'src/hardhat/modules/multi-sig/multi-sig.service';
import { SubmitTransactionDto } from './dto/multi-sig.dto';
@ApiTags('multi-sig-interact')
@Controller()
export class MultiSigInteractController {
  constructor(private readonly multiSigWalletService: MultiSigWalletService) {}

  @Get('owners/:address')
  async getOwners(@Param('address') address: string) {
    return this.multiSigWalletService.getOwners(address);
  }

  @ApiOkResponse()
  @Post('submit-transaction')
  async submitTransaction(@Body() dto: SubmitTransactionDto) {
    return this.multiSigWalletService.submitTransaction(dto);
  }
}
