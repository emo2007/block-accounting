import { ApiProperty } from '@nestjs/swagger';
import { IsOptional, IsString } from 'class-validator';

export class SubmitTransactionDto {
  @IsString()
  @ApiProperty()
  contractAddress: string;
  @ApiProperty()
  @IsString()
  destination: string;
  @IsString()
  @ApiProperty()
  value: string;
  @IsOptional()
  @IsString()
  //   @ApiProperty()
  data: string;
}
