import { useQuery } from "@tanstack/react-query"
import { useEffect, useState } from "react"
import getUserInfo from "../../api/getUserInfo"
import LoginButton from "./loginButton"
import UserInfo from "./userInfo"

import './header.css'

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
    <div className="header-children">
      {(isAuthed && query.data) ?
        <UserInfo user={query.data.user_info} /> :
        <LoginButton />
      }
    </div>
  )
}

export default Header
