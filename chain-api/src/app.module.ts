import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ContractFactoryModule } from './contract-factory/contract-factory.module';
import { ContractInteractModule } from './contract-interact/contract-interact.module';

@Module({
  imports: [ContractFactoryModule, ContractInteractModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
