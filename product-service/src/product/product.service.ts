import { Injectable, NotFoundException } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Product } from './product.entity';
import { CreateProductDto } from './dto/create-product.dto';
import { EventProducer } from './events/product.producer';
import Redis from 'ioredis';

@Injectable()
export class ProductService {
  private redis = new Redis(process.env.REDIS_URL);

  constructor(
    @InjectRepository(Product) private repo: Repository<Product>,
    private producer: EventProducer
  ) {}

  async createProduct(dto: CreateProductDto) {
    const product = this.repo.create(dto);
    await this.repo.save(product);
    await this.producer.emitEvent('product.created',product);
    return product;
  }

  async getProductById(id: number) {
    const cacheKey = `product:${id}`;
    const cached = await this.redis.get(cacheKey);
    if (cached) return JSON.parse(cached);

    const product = await this.repo.findOne({ where: { id } });
    if (product) {
      await this.redis.set(cacheKey, JSON.stringify(product), 'EX', 60);
      return product;
    }
    throw new NotFoundException(`Product with ID ${id} not found.`);
  }

  async decreaseStock(productId: number, qty: number) {
    await this.repo.decrement({ id: productId }, 'qty', qty);
    await this.redis.del(`product:${productId}`);
  }

}
