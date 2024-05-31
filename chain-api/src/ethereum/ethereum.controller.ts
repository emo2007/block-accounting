import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { EthereumService } from './ethereum.service';
import { ApiTags } from '@nestjs/swagger';
import { GetSeedPhraseDto } from './ethereum.dto';

@ApiTags('Ethereum')
@Controller()
export class EthereumController {
  constructor(private readonly ethereumService: EthereumService) {}

  @Get('/address/:privateKey')
  async getAddressFromPrivateKey(@Param('privateKey') privateKey: string) {
    return this.ethereumService.getAddressFromPrivateKey(privateKey);
  }

  @Post('/address-from-seed')
  async getAddressFromSeedPhrase(@Body() body: GetSeedPhraseDto) {
    return this.ethereumService.getAddressFromSeedPhrase(body.seedPhrase);
  }
}
