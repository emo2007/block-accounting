import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { ApiOkResponse, ApiTags } from '@nestjs/swagger';

import { MultiSigWalletDto } from './multi-sig.dto';
import { MultiSigWalletService } from './multi-sig.service';
import {
  ConfirmTransactionDto,
  DeployMultiSigResponseDto,
  DepositContractDto,
  ExecuteTransactionDto,
  GetTransactionDto,
  RevokeConfirmationDto,
  SubmitTransactionDto,
} from '../multi-sig.dto';
import { GetHeader } from '../../decorators/getHeader.decorator';
@ApiTags('multi-sig')
@Controller('multi-sig')
export class MultiSigInteractController {
  constructor(private readonly multiSigWalletService: MultiSigWalletService) {}

  @ApiOkResponse({
    type: DeployMultiSigResponseDto,
  })
  @Post('deploy')
  async deploy(
    @Body() dto: MultiSigWalletDto,
    @GetHeader('X-Seed') seed: string,
  ): Promise<DeployMultiSigResponseDto> {
    const addr = await this.multiSigWalletService.deploy(dto, seed);
    return {
      address: addr,
    };
  }
  @Get('owners/:address')
  async getOwners(
    @Param('address') address: string,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.multiSigWalletService.getOwners(address, seed);
  }

  @ApiOkResponse()
  @Post('submit-transaction')
  async submitTransaction(
    @Body() dto: SubmitTransactionDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.multiSigWalletService.submitTransaction(dto, seed);
  }

  @ApiOkResponse()
  @Post('confirm-transaction')
  async confirmTransaction(
    @Body() dto: ConfirmTransactionDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.multiSigWalletService.confirmTransaction(dto, seed);
  }

  @ApiOkResponse()
  @Post('execute-transaction')
  async executeTransaction(
    @Body() dto: ExecuteTransactionDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.multiSigWalletService.executeTransaction(dto, seed);
  }

  @ApiOkResponse()
  @Post('revoke-confirmation')
  async revokeConfirmation(
    @Body() dto: RevokeConfirmationDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.multiSigWalletService.revokeConfirmation(dto, seed);
  }

  @Get('transaction-count/:contractAddress')
  async getTransactionCount(
    @Param('contractAddress') contractAddress: string,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.multiSigWalletService.getTransactionCount(
      contractAddress,
      seed,
    );
  }

  @Get('transaction')
  async getTransaction(
    @Body() dto: GetTransactionDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.multiSigWalletService.getTransaction(dto, seed);
  }

  @Post('deposit')
  async deposit(
    @Body() dto: DepositContractDto,
    @GetHeader('X-Seed') seed: string,
  ) {
    return this.multiSigWalletService.deposit(dto, seed);
  }
}
