import { createParamDecorator, ExecutionContext } from '@nestjs/common';

export const GetHeader = createParamDecorator(
  (data: string, ctx: ExecutionContext) => {
    const request = ctx.switchToHttp().getRequest();
    return request.headers[data.toLowerCase()];
  },
);
