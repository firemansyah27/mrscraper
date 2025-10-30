import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Product } from './product.entity';
import { ProductService } from './product.service';
import { ProductController } from './product.controller';
import { EventProducer } from './events/product.producer';
import { OrderConsumer } from './events/order.consumer';

@Module({
  imports: [TypeOrmModule.forFeature([Product])],
  controllers: [ProductController],
  providers: [ProductService, EventProducer, OrderConsumer],
})
export class ProductModule {}
