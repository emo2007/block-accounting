import { Test, TestingModule } from '@nestjs/testing';
import { ContractInteractController } from './contract-interact.controller';
import { ContractInteractService } from './contract-interact.service';

describe('ContractInteractController', () => {
  let controller: ContractInteractController;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [ContractInteractController],
      providers: [ContractInteractService],
    }).compile();

    controller = module.get<ContractInteractController>(ContractInteractController);
  });

  it('should be defined', () => {
    expect(controller).toBeDefined();
  });
});
