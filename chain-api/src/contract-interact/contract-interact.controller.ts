import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
} from '@nestjs/common';
import { ContractInteractService } from './contract-interact.service';
import { CreateContractInteractDto } from './dto/create-contract-interact.dto';
import { UpdateContractInteractDto } from './dto/update-contract-interact.dto';
import { ApiTags } from '@nestjs/swagger';

@ApiTags('contract-interact')
@Controller('contract-interact')
export class ContractInteractController {
  constructor(
    private readonly contractInteractService: ContractInteractService,
  ) {}

  @Post()
  create(@Body() createContractInteractDto: CreateContractInteractDto) {
    return this.contractInteractService.create(createContractInteractDto);
  }

  @Get()
  findAll() {
    return this.contractInteractService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.contractInteractService.findOne(+id);
  }

  @Patch(':id')
  update(
    @Param('id') id: string,
    @Body() updateContractInteractDto: UpdateContractInteractDto,
  ) {
    return this.contractInteractService.update(+id, updateContractInteractDto);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.contractInteractService.remove(+id);
  }
}
