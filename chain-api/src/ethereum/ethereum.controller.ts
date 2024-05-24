import { Controller, Get, Param } from '@nestjs/common';
import { EthereumService } from './ethereum.service';
import { ApiTags } from '@nestjs/swagger';
@ApiTags('Ethereum')
@Controller()
export class EthereumController {
  constructor(private readonly ethereumService: EthereumService) {}
  @Get('/address/:privateKey')
  async getAddressFromPrivateKey(@Param('privateKey') privateKey: string) {
    return this.ethereumService.getAddressFromPrivateKey(privateKey);
  }
}
