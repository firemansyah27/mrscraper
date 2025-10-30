import { Entity, Column, PrimaryGeneratedColumn, CreateDateColumn } from 'typeorm';

@Entity('products')
export class Product {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  name: string;

  @Column('float')
  price: number;

  @Column('int')
  qty: number;

  @CreateDateColumn()
  createdAt: Date;
}
