import { Module } from '@nestjs/common';
import { AgreementController } from './agreement.controller';
import { AgreementService } from './agreement.service';
import { BaseModule } from '../../base/base.module';
import { MultiSigModule } from '../multi-sig/multi-sig.module';

@Module({
  imports: [BaseModule, MultiSigModule],
  controllers: [AgreementController],
  providers: [AgreementService],

  exports: [],
})
export class AgreementModule {}
