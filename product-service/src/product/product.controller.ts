import { Controller, Post, Get, Body, Param } from '@nestjs/common';
import { ProductService } from './product.service';
import { CreateProductDto } from './dto/create-product.dto';

@Controller('products')
export class ProductController {
  constructor(private readonly service: ProductService) {}

  @Post()
  create(@Body() dto: CreateProductDto) {
    return this.service.createProduct(dto);
  }

  @Get(':id')
  get(@Param('id') id: string) {
    const pid = Number(id);
    return this.service.getProductById(pid);
  }
}
