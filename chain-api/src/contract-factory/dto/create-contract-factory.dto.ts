import { ApiProperty } from '@nestjs/swagger';
export enum ContractType {
  SALARY,
  AGREEMENT,
}

export class CreateContractFactoryDto {
  @ApiProperty({
    enum: ContractType,
  })
  contractType: ContractType;
}
