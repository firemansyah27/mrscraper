import { Controller, Get } from '@nestjs/common';

@Controller()
export class MainController {
  @Get('health')
  getHealth() {
    return {
      status: 'ok',
      service: 'product-service',
      timestamp: new Date().toISOString(),
    };
  }
}