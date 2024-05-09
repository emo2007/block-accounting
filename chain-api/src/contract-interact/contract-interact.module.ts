import { Module } from '@nestjs/common';
import { ContractInteractService } from './contract-interact.service';

import { HardhatModule } from 'src/hardhat/modules/hardhat.module';
import { MultiSigInteractController } from './multi-sig-interact.controller';

@Module({
  imports: [HardhatModule],
  controllers: [MultiSigInteractController],
  providers: [ContractInteractService],
})
export class ContractInteractModule {}
