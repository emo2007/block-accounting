import { ApiProperty } from '@nestjs/swagger';
import { IsString, IsUrl } from 'class-validator';

export class DeployAgreementDto {
  @ApiProperty()
  @IsString()
  multiSigWallet: string;
}
export class GetAgreementInfoDto {
  @ApiProperty()
  @IsString()
  contractAddress: string;
}
export class RequestAgreementDto extends GetAgreementInfoDto {
  @ApiProperty()
  @IsString()
  multiSigWallet: string;
  @ApiProperty()
  @IsUrl()
  url: string;
}
