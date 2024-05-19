import {
  ArgumentsHost,
  Catch,
  ExceptionFilter,
  HttpException,
  HttpStatus,
} from '@nestjs/common';
import { error } from 'console';
import { Request, Response } from 'express';

@Catch()
export class AllExceptionsFilter implements ExceptionFilter {
  //   constructor(private readonly httpAdapterHost: HttpAdapterHost) {}

  catch(exception: any, host: ArgumentsHost): void {
    console.log('🚀 ~ AllExceptionsFilter ~ exception:', exception);
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const httpStatus =
      exception instanceof HttpException
        ? exception.getStatus()
        : HttpStatus.INTERNAL_SERVER_ERROR;

    const responseBody = {
      statusCode: httpStatus,
      error: exception?.info?.error.message || exception.toString(),
      timestamp: new Date().toISOString(),
    };

    response.status(500).json(responseBody);
  }
}
