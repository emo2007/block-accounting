import { Module } from '@nestjs/common';
import { EthereumController } from './ethereum.controller';
import { EthereumService } from './ethereum.service';

@Module({
  imports: [],
  controllers: [EthereumController],
  providers: [EthereumService],
  exports: [],
})
export class EthereumModule {}
