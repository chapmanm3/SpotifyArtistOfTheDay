import { useQuery } from "@tanstack/react-query"
import { useEffect, useState } from "react"
import getUserInfo from "../../api/getUserInfo"
import LoginButton from "./loginButton"

function Header() {

  const [isAuthed, setIsAuthed] = useState(window.localStorage.getItem("is_authed") === "true")

  useEffect(() => {
    const checkIsAuthed = () => {
      const isAuthed = window.localStorage.getItem("is_authed") === "true"
      setIsAuthed(isAuthed)
      debugger
    }
    window.onstorage = () => checkIsAuthed()
  }, [])

  const query = useQuery({
    queryKey: ['userInfo'],
    queryFn: getUserInfo,
    enabled: isAuthed,
  })

  return (
    isAuthed ?
      <div>
        Hello: {query.data?.user_info.DisplayName}
      </div> :
      <LoginButton />
  )
}

export default Header
