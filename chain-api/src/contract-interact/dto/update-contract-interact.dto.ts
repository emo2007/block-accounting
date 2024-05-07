import { PartialType } from '@nestjs/swagger';
import { CreateContractInteractDto } from './create-contract-interact.dto';

export class UpdateContractInteractDto extends PartialType(CreateContractInteractDto) {}
