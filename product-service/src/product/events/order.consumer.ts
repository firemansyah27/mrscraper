import * as amqp from 'amqplib';
import { Injectable, OnModuleInit } from '@nestjs/common';
import { ProductService } from '../product.service';
import { EventProducer } from './product.producer';


@Injectable()
export class OrderConsumer implements OnModuleInit {
  constructor(
    private productService: ProductService,
    private eventProducer: EventProducer,
  ) {}

  async onModuleInit() {
    try {
      const conn = await amqp.connect(process.env.RABBITMQ_URL);
      const ch = await conn.createChannel();
      await ch.assertExchange('events', 'topic', { durable: true });
      const q = await ch.assertQueue('order-queue', { durable: true });
      await ch.bindQueue(q.queue, 'events', 'update.product.stock');
      const productQueue = await ch.assertQueue('product-queue', { durable: true });
      await ch.bindQueue(productQueue.queue, 'events', 'product.created');

      ch.consume(q.queue, async msg => {
        if (!msg) return;

        try {
          const payload = JSON.parse(msg.content.toString());
          const { order_id, productId, quantity } = payload.data;

          try {
            const availableStock = await this.productService.getProductById(productId);
            if (availableStock.qty < quantity) {
              console.log(`Not enough stock for product ${productId}. Available stock: ${availableStock.qty}, requested: ${quantity}`);
              await this.emitOrderEvent(order_id, 'Out of Stock');
              ch.nack(msg)
              return;
            }

            await this.productService.decreaseStock(productId, quantity);
            console.log(`Reduced stock for product ${productId} by ${quantity}`);

            await this.emitOrderEvent(order_id, 'Sale');

            ch.ack(msg);
          } catch (error) {
            console.error(`Error processing order for product ${productId}: ${error.message}`);
            ch.nack(msg, false, false);
          }

        } catch (e) {
          console.error('Error processing update.product.stock event:', e);
          ch.nack(msg, false, false);
        }
      }, { noAck: false });

      console.log('🎧 Product-service listening for update.product.stock events');
    } catch (err) {
      console.error('Failed to start OrderConsumer', err);
    }
  }

  async emitOrderEvent(orderId: number, status: string) {
    const payload = { orderId: orderId, status: status};
    await this.eventProducer.emitEvent('update.order.status', payload);
  }
}
