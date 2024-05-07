import { Test, TestingModule } from '@nestjs/testing';
import { ContractFactoryController } from './contract-factory.controller';
import { ContractFactoryService } from './contract-factory.service';

describe('ContractFactoryController', () => {
  let controller: ContractFactoryController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [ContractFactoryController],
      providers: [ContractFactoryService],
    }).compile();

    controller = module.get<ContractFactoryController>(ContractFactoryController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});
