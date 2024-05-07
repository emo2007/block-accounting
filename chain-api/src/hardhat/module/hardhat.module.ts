import { Module } from '@nestjs/common';
import { HardhatService } from './hardhat.service';

@Module({
  imports: [],
  controllers: [],
  providers: [HardhatService],
  exports: [HardhatService],
})
export class HardhatModule {}
