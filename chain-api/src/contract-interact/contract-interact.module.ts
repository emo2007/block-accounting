import { Module } from '@nestjs/common';
import { ContractInteractService } from './contract-interact.service';
import { ContractInteractController } from './contract-interact.controller';

@Module({
  controllers: [ContractInteractController],
  providers: [ContractInteractService],
})
export class ContractInteractModule {}
