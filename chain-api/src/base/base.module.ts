import { Module } from '@nestjs/common';
import { ProviderModule } from './provider/provider.module';

@Module({
  imports: [ProviderModule],
  controllers: [],
  providers: [],
  exports: [ProviderModule],
})
export class BaseModule {}
