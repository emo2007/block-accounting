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
@ApiTags('contract-factory')
@Controller('contract-factory')
export class ContractFactoryController {
  constructor(
    private readonly contractFactoryService: ContractFactoryService,
  ) {}

  @Post('')
  create(@Body() createContractFactoryDto: CreateContractFactoryDto) {
    return this.contractFactoryService.create(createContractFactoryDto);
  }
}
