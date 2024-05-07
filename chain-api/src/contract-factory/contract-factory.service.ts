import { HardhatService } from '../hardhat/module/hardhat.service';
import { Injectable } from '@nestjs/common';
import { CreateContractFactoryDto } from './dto/create-contract-factory.dto';

@Injectable()
export class ContractFactoryService {
  constructor(private readonly hhService: HardhatService) {}
  async create(createContractFactoryDto: CreateContractFactoryDto) {
    return await this.hhService.deploySalaryContract();
  }
}
