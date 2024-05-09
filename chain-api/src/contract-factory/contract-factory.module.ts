import { HardhatModule } from '../hardhat/modules/hardhat.module';
import { Module } from '@nestjs/common';
import { ContractFactoryService } from './contract-factory.service';
import { ContractFactoryController } from './contract-factory.controller';

@Module({
  imports: [HardhatModule],
  controllers: [ContractFactoryController],
  providers: [ContractFactoryService],
})
export class ContractFactoryModule {}
