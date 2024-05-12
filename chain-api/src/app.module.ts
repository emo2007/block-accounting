import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';

import { ContractInteractModule } from './contract-interact/contract-interact.module';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    ContractInteractModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
