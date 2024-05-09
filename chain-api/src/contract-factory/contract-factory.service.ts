import { HardhatService } from '../hardhat/modules/hardhat.service';
import { Injectable } from '@nestjs/common';
import { CreateContractFactoryDto } from './dto/create-contract-factory.dto';
import { SalariesService } from 'src/hardhat/modules/salary.service';
import { MultiSigWalletService } from 'src/hardhat/modules/multi-sig/multi-sig.service';
import { MultiSigWalletDto } from 'src/hardhat/modules/dto/multi-sig.dto';

@Injectable()
export class ContractFactoryService {
  constructor(
    private readonly salaryService: SalariesService,
    private readonly multiSigService: MultiSigWalletService,
  ) {}
  async createSalary(createContractFactoryDto: CreateContractFactoryDto) {
    return await this.salaryService.deploy();
  }

  async createMultiSig(dto: MultiSigWalletDto) {
    return await this.multiSigService.deploy(dto);
  }
}
