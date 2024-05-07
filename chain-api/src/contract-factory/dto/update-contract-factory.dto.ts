import { PartialType } from '@nestjs/mapped-types';
import { CreateContractFactoryDto } from './create-contract-factory.dto';

export class UpdateContractFactoryDto extends PartialType(CreateContractFactoryDto) {}
