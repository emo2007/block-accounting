import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
} from '@nestjs/common';
import { ContractFactoryService } from './contract-factory.service';
import { CreateContractFactoryDto } from './dto/create-contract-factory.dto';
import { UpdateContractFactoryDto } from './dto/update-contract-factory.dto';
import { ApiTags } from '@nestjs/swagger';
import { MultiSigWalletDto } from 'src/hardhat/modules/dto/multi-sig.dto';
@ApiTags('contract-factory')
@Controller('contract-factory')
export class ContractFactoryController {
  constructor(
    private readonly contractFactoryService: ContractFactoryService,
  ) {}

  @Post('multi-sig')
  create(@Body() createContractFactoryDto: MultiSigWalletDto) {
    return this.contractFactoryService.createMultiSig(createContractFactoryDto);
  }
}
