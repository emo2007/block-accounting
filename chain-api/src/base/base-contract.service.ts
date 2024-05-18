import { Injectable } from '@nestjs/common';
import { ProviderService } from './provider/provider.service';

@Injectable()
export abstract class BaseContractService {
  constructor(public readonly providerService: ProviderService) {}
  abstract deploy(dto: object): Promise<any>;
}
