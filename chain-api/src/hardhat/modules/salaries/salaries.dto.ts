import { ApiProperty } from '@nestjs/swagger';
import { IsNumber, IsString } from 'class-validator';

export class SalariesDeployDto {
  @ApiProperty()
  @IsString()
  multiSigWallet: string;
}

export class SetSalaryDto {
  @ApiProperty()
  @IsString()
  multiSigWallet: string;
  @ApiProperty()
  @IsString()
  contractAddress: string;
  @ApiProperty()
  @IsString()
  employeeAddress: string;
  @ApiProperty()
  @IsNumber()
  salary: number;
}

export class GetEmployeeSalariesDto {
  @ApiProperty()
  @IsString()
  contractAddress: string;
  @ApiProperty()
  @IsString()
  employeeAddress: string;
}
