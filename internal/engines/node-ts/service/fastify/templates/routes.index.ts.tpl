import { FastifyPluginAsync } from 'fastify'

const routes: FastifyPluginAsync[] = [
{{- range .Items }}
  (await import('./{{ lower .Name }}.route')).default,
{{- end }}
]

export default routes
