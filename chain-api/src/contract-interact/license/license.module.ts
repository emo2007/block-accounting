import { Module } from '@nestjs/common';
import { LicenseController } from './license.controller';
import { LicenseService } from './license.service';
import { BaseModule } from '../../base/base.module';
import { MultiSigModule } from '../multi-sig/multi-sig.module';

@Module({
  imports: [BaseModule, MultiSigModule],
  controllers: [LicenseController],
  providers: [LicenseService],
  exports: [LicenseService],
})
export class LicenseModule {}
