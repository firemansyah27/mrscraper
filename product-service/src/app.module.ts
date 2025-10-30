import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { ProductModule } from './product/product.module';
import { Product } from './product/product.entity';
import { MainController } from './main.controller';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: 'postgres',
      url: process.env.DATABASE_URL,
      entities: [Product],
      synchronize: true,
    }),
    ProductModule,
  ],
  controllers: [MainController]
})
export class AppModule {}
