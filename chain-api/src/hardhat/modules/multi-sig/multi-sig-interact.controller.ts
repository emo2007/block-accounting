import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { ApiOkResponse, ApiTags } from '@nestjs/swagger';
import { MultiSigWalletService } from 'src/hardhat/modules/multi-sig/multi-sig.service';
import {
  ConfirmTransactionDto,
  DepositContractDto,
  ExecuteTransactionDto,
  GetTransactionDto,
  RevokeConfirmationDto,
  SubmitTransactionDto,
} from '../../../contract-interact/dto/multi-sig.dto';
import { MultiSigWalletDto } from './multi-sig.dto';
@ApiTags('multi-sig')
@Controller('multi-sig')
export class MultiSigInteractController {
  constructor(private readonly multiSigWalletService: MultiSigWalletService) {}

  @Post('deploy')
  async deploy(@Body() dto: MultiSigWalletDto) {
    return this.multiSigWalletService.deploy(dto);
  }
  @Get('owners/:address')
  async getOwners(@Param('address') address: string) {
    return this.multiSigWalletService.getOwners(address);
  }

  @ApiOkResponse()
  @Post('submit-transaction')
  async submitTransaction(@Body() dto: SubmitTransactionDto) {
    return this.multiSigWalletService.submitTransaction(dto);
  }

  @ApiOkResponse()
  @Post('confirm-transaction')
  async confirmTransaction(@Body() dto: ConfirmTransactionDto) {
    return this.multiSigWalletService.confirmTransaction(dto);
  }

  @ApiOkResponse()
  @Post('execute-transaction')
  async executeTransaction(@Body() dto: ExecuteTransactionDto) {
    return this.multiSigWalletService.executeTransaction(dto);
  }

  @ApiOkResponse()
  @Post('revoke-confirmation')
  async revokeConfirmation(@Body() dto: RevokeConfirmationDto) {
    return this.multiSigWalletService.revokeConfirmation(dto);
  }

  @Get('transaction-count/:contractAddress')
  async getTransactionCount(@Param('contractAddress') contractAddress: string) {
    return this.multiSigWalletService.getTransactionCount(contractAddress);
  }

  @Get('transaction')
  async getTransaction(@Body() dto: GetTransactionDto) {
    return this.multiSigWalletService.getTransaction(dto);
  }

  @Post('deposit')
  async deposit(@Body() dto: DepositContractDto) {
    return this.multiSigWalletService.deposit(dto);
  }
}
