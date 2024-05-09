import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { ProviderService } from 'src/provider/provider.service';

@Injectable()
export abstract class BaseContractService {
  constructor(
    public readonly configService: ConfigService,
    public readonly providerService: ProviderService,
  ) {}
  abstract deploy(dto: object): Promise<any>;
}
