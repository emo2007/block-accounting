import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { ContractFactoryModule } from './contract-factory/contract-factory.module';
import { ContractInteractModule } from './contract-interact/contract-interact.module';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    ContractFactoryModule,
    ContractInteractModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
