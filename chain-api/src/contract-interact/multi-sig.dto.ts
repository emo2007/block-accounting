import { ApiProperty } from '@nestjs/swagger';
import { IsBoolean, IsNumber, IsOptional, IsString } from 'class-validator';

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
  @ApiProperty()
  data: string;
}

export class ConfirmTransactionDto {
  @IsString()
  @ApiProperty()
  contractAddress: string;
  @ApiProperty()
  @IsNumber()
  index: number;
}

export class ExecuteTransactionDto extends ConfirmTransactionDto {
  @IsOptional()
  @IsBoolean()
  @ApiProperty()
  isDeploy: boolean;
}

export class RevokeConfirmationDto extends ConfirmTransactionDto {}

export class GetTransactionCount {}

export class GetTransactionDto extends ConfirmTransactionDto {}

export class DepositContractDto {
  @IsString()
  @ApiProperty()
  contractAddress: string;
  @IsString()
  @ApiProperty()
  value: string;
}
export class DeployMultiSigResponseDto {
  @IsString()
  address: string;
}
