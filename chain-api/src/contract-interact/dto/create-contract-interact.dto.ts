import { ApiProperty } from '@nestjs/swagger';

export class CreateContractInteractDto {
  @ApiProperty()
  contractAddress: string;
  @ApiProperty()
  sender: string;
}
