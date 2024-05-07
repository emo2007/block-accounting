import { Test, TestingModule } from '@nestjs/testing';
import { ContractInteractService } from './contract-interact.service';

describe('ContractInteractService', () => {
  let service: ContractInteractService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [ContractInteractService],
    }).compile();

    service = module.get<ContractInteractService>(ContractInteractService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});
