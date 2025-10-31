import {createFileRoute, useNavigate} from '@tanstack/react-router'
import {useAuthStore} from "../stores/authStore.ts";
import {Button} from "@mui/material";

export const Route = createFileRoute('/_auth/')({
    component: RouteComponent,
})

function RouteComponent() {
    const { logout } = useAuthStore();
    const navigate = useNavigate();

    const handleClick = () => {
        logout()
        navigate({
            to: '/login',
            replace: true,
        })
    }

  return (
      <>
          <div>Авторизовались))</div>
          <Button variant={'outlined'} onClick={handleClick}>Выйти</Button>
      </>
  )
}
