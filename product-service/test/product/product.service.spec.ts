import { Test, TestingModule } from '@nestjs/testing';
import { ProductService } from '../../src/product/product.service';
import { getRepositoryToken } from '@nestjs/typeorm';
import { Product } from '../../src/product/product.entity';
import { EventProducer } from '../../src/product/events/product.producer';
import { mockProductRepository } from './mocks/product.repository.mock';
import { mockEventProducer } from './mocks/event.producer.mock';


describe('ProductService', () => {
  let service: ProductService;

  beforeEach(async () => {
    
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ProductService,
        { provide: getRepositoryToken(Product), useValue: mockProductRepository },
        { provide: EventProducer, useValue: mockEventProducer },
      ],
    }).compile();
    
    service = module.get<ProductService>(ProductService);
    jest.clearAllMocks();
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('createProduct', () => {
    it('should successfully create a product and emit event', async () => {
      const dto = { name: 'Product A', price: 100000, qty: 50000 };
      const savedProduct = { id: 1, ...dto };
  
      mockProductRepository.create.mockImplementation(input => ({ id: 1, ...input }));
      mockProductRepository.save.mockResolvedValue(savedProduct);
  
      const result = await service.createProduct(dto);
  
      expect(mockProductRepository.create).toHaveBeenCalledWith(dto);
      expect(mockProductRepository.save).toHaveBeenCalledWith({ id: 1, ...dto });
      expect(mockEventProducer.emitEvent).toHaveBeenCalledWith('product.created', savedProduct);
      expect(result).toEqual(savedProduct);
    });
  });
});
