import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/banks')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_auth/banks"!</div>
}
