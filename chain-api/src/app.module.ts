import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';

import { ContractInteractModule } from './contract-interact/contract-interact.module';
import { ConfigModule } from '@nestjs/config';
import { EthereumModule } from './ethereum/ethereum.module';
import { AgreementModule } from './contract-interact/agreement/agreement.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    ContractInteractModule,
    EthereumModule,
    AgreementModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
