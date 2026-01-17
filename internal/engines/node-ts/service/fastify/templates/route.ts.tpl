import { FastifyPluginAsync } from 'fastify'

const plugin: FastifyPluginAsync = async (fastify) => {
  // Route: {{ .Name }}
  fastify.get('/{{ lower .Name }}', async (request, reply) => {
    return { message: 'GET {{ .Name }}' }
  })
}

export default plugin
