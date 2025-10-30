import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 1000,
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 200,
      maxVUs: 1000,
    },
  },
};

const BASE_ORDER_URL = 'http://localhost:4000';
const BASE_PRODUCT_URL = 'http://localhost:3000';

export function setup() {
  const payload = JSON.stringify({
    name: `Product A ${Date.now()}`,
    price: 100000,
    qty: 50000,
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const res = http.post(`${BASE_PRODUCT_URL}/products`, payload, params);

  check(res, {
    'create product success': (r) => r.status === 201 || r.status === 200,
  });

  const product = JSON.parse(res.body);
  const productId = product.id || product.product?.id;

  if (!productId) {
    throw new Error('Gagal mendapatkan product_id dari response create product');
  }

  console.log(`✅ Created product with id=${productId}`);

  return { productId };
}

export default function (data) {
  const payload = JSON.stringify({
    product_id: data.productId,
    quantity: 1,
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const res = http.post(`${BASE_ORDER_URL}/orders`, payload, params);

  check(res, {
    'status is 202 Accepted': (r) => r.status === 202,
  });
}
