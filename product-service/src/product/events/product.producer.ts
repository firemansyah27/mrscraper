import * as amqp from 'amqplib';
import { Injectable } from '@nestjs/common';

@Injectable()
export class EventProducer {
  private exchangeName = 'events';

  async emitEvent(routingKey: string, data: any) {
    try {
      const conn = await amqp.connect(process.env.RABBITMQ_URL);
      const channel = await conn.createChannel();

      await channel.assertExchange(this.exchangeName, 'topic', { durable: true });

      const payload = {
        event: routingKey,
        timestamp: new Date().toISOString(),
        data: data,
      };

      channel.publish(this.exchangeName, routingKey, Buffer.from(JSON.stringify(payload)));

      await channel.close();
      await conn.close();

      console.log(`Emitted event with routing key: ${routingKey}`);
    } catch (err) {
      console.error(`Failed to emit event with routing key: ${routingKey}`, err);
    }
  }
}
