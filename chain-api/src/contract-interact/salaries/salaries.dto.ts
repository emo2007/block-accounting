import { ApiProperty } from '@nestjs/swagger';
import { IsNumber, IsString } from 'class-validator';

export class SalariesDeployDto {
  @ApiProperty()
  @IsString()
  authorizedWallet: string;
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
export class GeneralEmpoyeeSalaryDto {
  @ApiProperty()
  @IsString()
  contractAddress: string;
  @ApiProperty()
  @IsString()
  employeeAddress: string;
}
export class GetEmployeeSalariesDto extends GeneralEmpoyeeSalaryDto {}

export class CreatePayoutDto extends GeneralEmpoyeeSalaryDto {
  @IsString()
  multiSigWallet: string;
}

export class DeployContractResponseDto {
  @IsString()
  address: string;
}
