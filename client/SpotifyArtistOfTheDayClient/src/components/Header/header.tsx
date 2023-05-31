import { useQuery } from "@tanstack/react-query"
import getUserInfo from "../../api/getUserInfo"

function Header() {

  const isAuthed = window.localStorage.getItem("is_authed")

  const query = useQuery({
    queryKey: ['userInfo'],
    queryFn: getUserInfo,
    enabled: !!isAuthed,
  })

  return (
    <div>
      Hello: {query.data?.user_info.DisplayName}
    </div>
  )
}

export default Header
