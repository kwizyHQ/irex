import { FastifyPluginAsync } from 'fastify'

const plugin: FastifyPluginAsync = async (fastify) => {
  fastify.{{ lower .Method }}('/{{ lower .Method }}', async (request, reply) => {
    return { message: 'GET {{ .Method }}' }
  })
}

export default plugin
