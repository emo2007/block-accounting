import { IsArray, IsNumber, IsString } from 'class-validator';
import { ApiProperty } from '@nestjs/swagger';
export class GetLicenseInfoDto {
  @ApiProperty()
  @IsString()
  contractAddress: string;
}
export class DeployLicenseDto {
  @ApiProperty()
  @IsString()
  multiSigWallet: string;
  @ApiProperty()
  @IsArray()
  owners: string[];
  @ApiProperty({
    isArray: true,
    type: Number,
  })
  @IsNumber({}, { each: true })
  shares: number[];
}

export class RequestLicenseDto extends GetLicenseInfoDto {
  @ApiProperty()
  @IsString()
  multiSigWallet: string;
}

export class GetLicenseResponseDto extends GetLicenseInfoDto {}

export class GetShareLicense extends GetLicenseInfoDto {
  @IsString()
  @ApiProperty()
  ownerAddress: string;
}

export class LicensePayoutDto extends RequestLicenseDto {}
export class SetPayoutContractDto extends RequestLicenseDto {
  @IsString()
  @ApiProperty()
  payoutContract: string;
}
