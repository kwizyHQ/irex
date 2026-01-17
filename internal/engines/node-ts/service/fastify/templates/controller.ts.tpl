import { FastifyPluginAsync } from 'fastify'
import { {{ .Name }} } from '../../models/{{ lower .Name }}'

const controller: FastifyPluginAsync = async (fastify) => {
  fastify.get('/{{ lower .Name }}s', async (request, reply) => {
    return await {{ .Name }}.find()
  })
}

export default controller
