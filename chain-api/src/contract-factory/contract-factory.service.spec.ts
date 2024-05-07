import { Test, TestingModule } from '@nestjs/testing';
import { ContractFactoryService } from './contract-factory.service';

describe('ContractFactoryService', () => {
  let service: ContractFactoryService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [ContractFactoryService],
    }).compile();

    service = module.get<ContractFactoryService>(ContractFactoryService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});
