import { Module } from '@nestjs/common';
import { HardhatModule } from 'src/hardhat/modules/hardhat.module';

@Module({
  imports: [HardhatModule],
  controllers: [],
  providers: [],
})
export class ContractInteractModule {}
