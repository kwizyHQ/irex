import { FastifyPluginAsync } from 'fastify'

const controllers: FastifyPluginAsync[] = [
{{- range .Items }}
  (await import('./{{ lower .Name }}.controller')).default,
{{- end }}
]

export default controllers
