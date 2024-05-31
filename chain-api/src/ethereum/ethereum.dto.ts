import {IsString} from "class-validator";
import {ApiProperty} from "@nestjs/swagger";

export class GetSeedPhraseDto {
    @ApiProperty()
    @IsString()
    seedPhrase: string;
}