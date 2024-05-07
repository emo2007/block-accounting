import { Injectable } from '@nestjs/common';
import { CreateContractInteractDto } from './dto/create-contract-interact.dto';
import { UpdateContractInteractDto } from './dto/update-contract-interact.dto';

@Injectable()
export class ContractInteractService {
  create(createContractInteractDto: CreateContractInteractDto) {
    return 'This action adds a new contractInteract';
  }

  findAll() {
    return `This action returns all contractInteract`;
  }

  findOne(id: number) {
    return `This action returns a #${id} contractInteract`;
  }

  update(id: number, updateContractInteractDto: UpdateContractInteractDto) {
    return `This action updates a #${id} contractInteract`;
  }

  remove(id: number) {
    return `This action removes a #${id} contractInteract`;
  }
}
