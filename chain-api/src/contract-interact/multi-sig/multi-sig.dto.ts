import { ApiProperty } from '@nestjs/swagger';
import { IsArray, IsNumber } from 'class-validator';

export class MultiSigWalletDto {
  @IsArray()
  @ApiProperty()
  owners: string[];
  @IsNumber()
  @ApiProperty()
  confirmations: number;
}
