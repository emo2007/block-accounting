import { Module } from '@nestjs/common';

import { MultiSigWalletService } from './multi-sig.service';
import { MultiSigInteractController } from './multi-sig-interact.controller';
import { BaseModule } from '../../base/base.module';

@Module({
  imports: [BaseModule],
  controllers: [MultiSigInteractController],
  providers: [MultiSigWalletService],
  exports: [MultiSigWalletService],
})
export class MultiSigModule {}
