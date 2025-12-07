import Fastify from 'fastify';

export async function start() {
  const app = Fastify({ logger: true });

  app.get('/', async () => ({ hello: 'world' }));

  const port = Number(process.env.PORT || 3000);
  await app.listen({ port });
  console.log('listening on ' + port);
}
